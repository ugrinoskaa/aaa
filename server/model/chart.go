package model

type ChartType string

const (
	BAR     ChartType = "bar"
	PIE     ChartType = "pie"
	LINE    ChartType = "line"
	SCATTER ChartType = "scatter"
	HEATMAP ChartType = "heatmap"
	SANKEY  ChartType = "sankey"
)

var (
	SupportedChartTypes = []ChartType{BAR, PIE, LINE, SCATTER, HEATMAP, SANKEY}
	SupportedPrecisions = []string{"year", "quarter", "month", "week", "day"}
	SupportedFilters    = []string{"=", "!=", ">", "<", "LIKE", "IN"}
)

type Chart struct {
	ID        int         `json:"id"`
	DatasetID int         `json:"datasetId"`
	Name      string      `json:"name"`
	Type      ChartType   `json:"type"`
	Config    ChartConfig `json:"config"`
}

type ChartConfig struct {
	Dimensions []string `json:"dimensions"`
	Metrics    []string `json:"metrics"`
	Filters    []string `json:"filters"`
}

type ChartSchema struct {
	Type    ChartType        `json:"type"`
	Schema  ChartSchemaRules `json:"schema"`
	Example interface{}      `json:"example"`
}

type ChartSchemaRules struct {
	Dimensions FieldRule `json:"dimensions"`
	Metrics    FieldRule `json:"metrics"`
	Filters    FieldRule `json:"filters"`
}

type FieldRule struct {
	Min    int      `json:"min"`
	Max    int      `json:"max"`
	Values []string `json:"values"`
}
