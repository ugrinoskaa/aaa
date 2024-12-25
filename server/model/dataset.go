package model

import (
	"fmt"
	"github.com/amukoski/aaa/service/utils"
)

var (
	SQLAggregationsGlobal = []string{"COUNT(*)"}
	SQLAggregationsColumn = []string{"COUNT(%s)", "AVG(%s)", "SUM(%s)", "MIN(%s)", "MAX(%s)"}
)

var (
	SupportedPrecisions = []string{"year", "quarter", "month", "week", "day"}
	SupportedFilters    = []string{"=", "!=", ">", "<", "LIKE", "IN"}
)

type Dataset struct {
	ID       int
	SourceID int
	Name     string
	Config   DatasetConfig
}

func (ds DatasetConfig) Dimensions() []string {
	dimensions := make([]string, 0, len(ds.Columns))

	for _, col := range ds.Columns {
		column, dataType := utils.ParseColumn(col)
		dimensions = append(dimensions, column)

		if utils.IsColumnDateTime(dataType) {
			for _, precision := range SupportedPrecisions {
				dimensions = append(dimensions, utils.FormatColumn(column, precision))
			}
		}
	}

	return dimensions
}

func (ds DatasetConfig) Metrics() []string {
	metrics := make([]string, 0)

	for _, op := range SQLAggregationsGlobal {
		metrics = append(metrics, op)
	}

	for _, col := range ds.Columns {
		column, dataType := utils.ParseColumn(col)
		if !utils.IsColumnNumeric(dataType) {
			continue
		}

		for _, op := range SQLAggregationsColumn {
			metrics = append(metrics, fmt.Sprintf(op, column))
		}
	}

	return metrics
}

func (ds DatasetConfig) Precisions() []string {
	precisions := make([]string, 0)

	for _, col := range ds.Columns {
		column, dataType := utils.ParseColumn(col)
		if !utils.IsColumnDateTime(dataType) {
			continue
		}

		precisions = append(precisions, column)
	}

	return precisions
}
