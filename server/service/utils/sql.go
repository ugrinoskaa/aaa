package utils

import (
	"fmt"
	"github.com/samber/lo"
	"slices"
	"strings"
)

func BuildSQLQuery(table string, dimensions []string, metrics []string, filters []string, columns []string) string {
	dimensionsSQL := buildDimensions(dimensions)
	metricsSQL := buildMetrics(metrics)
	selectSQL := buildSelect(dimensionsSQL, metricsSQL)
	whereSQL := buildWhere(filters, columns)

	return fmt.Sprintf(`SELECT %s FROM %s WHERE %s GROUP BY %s ORDER BY %s`, selectSQL, table, whereSQL, dimensionsSQL, dimensionsSQL)
}

func buildWhere(filters []string, columns []string) interface{} {
	normalized := make([]string, 0)
	normalized = append(normalized, "1=1")

	for _, filter := range filters {
		parts := strings.Split(filter, "/")
		if len(parts) != 3 {
			continue
		}

		dimension, operator, value := parts[0], parts[1], parts[2]

		column, precision := ParseColumn(dimension)
		if precision != "" {
			query := fmt.Sprintf("EXTRACT(%s FROM %s)::text %s '%s'", precision, column, operator, value)
			normalized = append(normalized, query)
			continue
		}

		idx := slices.IndexFunc(columns, func(v string) bool {
			dbcol, _ := ParseColumn(v)
			if strings.EqualFold(dbcol, dimension) {
				return true
			}

			return false
		})

		if idx != -1 {
			_, dataType := ParseColumn(columns[idx])
			if !IsColumnNumeric(dataType) {
				if operator == "LIKE" {
					value = fmt.Sprintf("'%%%s%%'", value)
				} else if operator == "IN" {
					splits := strings.Split(value, ",")
					inop := lo.Map(splits, func(item string, index int) string {
						return fmt.Sprintf("'%s'", item)
					})
					value = fmt.Sprintf("(%s)", strings.Join(inop, ","))
				} else {
					value = fmt.Sprintf("'%s'", value)
				}
			}
		}

		normalized = append(normalized, fmt.Sprintf("%s %s %s", column, operator, value))
	}

	return strings.Join(normalized, " AND ")
}

func buildDimensions(dimensions []string) string {
	normalized := make([]string, len(dimensions))

	for idx, col := range dimensions {
		column, precision := ParseColumn(col)
		if precision != "" {
			normalized[idx] = fmt.Sprintf("DATE_TRUNC('%s', %s)::date::text", precision, column)
			continue
		}

		normalized[idx] = column
	}

	return strings.Join(normalized, ",")
}

func buildMetrics(metrics []string) string {
	return strings.Join(metrics, ",")
}

func buildSelect(dimensions, metrics string) string {
	if dimensions == "" {
		return metrics
	}

	if metrics == "" {
		return dimensions
	}

	return fmt.Sprintf("%s,%s", dimensions, metrics)
}
