package utils

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/jackc/pgtype"
)

var (
	ColumnSeparator = "::"
)

func ToFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case int:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case pgtype.Numeric:
		var f float64
		err := v.AssignTo(&f)
		if err != nil {
			return 0, fmt.Errorf("failed to convert pgtype.Numeric to float64: %w", err)
		}
		return f, nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("unsupported value type: %v", reflect.TypeOf(value))
	}
}

func IsColumnNumeric(dataType string) bool {
	numeric := []string{"integer", "numeric", "decimal", "real", "double precision"}
	return slices.Contains(numeric, dataType)
}

func IsColumnDateTime(dataType string) bool {
	dateTime := []string{"date", "timestamp without time zone", "timestamp with time zone"}
	return slices.Contains(dateTime, dataType)
}

func ParseColumn(column string) (string, string) {
	parts := strings.SplitN(column, ColumnSeparator, 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}

	return parts[0], ""
}

func FormatColumn(name string, dtype string) string {
	return name + ColumnSeparator + dtype
}

func ColumnNames(columns []string) []string {
	names := make([]string, len(columns))
	for i, column := range columns {
		name, _ := ParseColumn(column)
		names[i] = name
	}

	return names
}
