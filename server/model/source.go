package model

type SourceType string

const (
	POSTGRES SourceType = "postgres"
	CSV      SourceType = "csv"
)

type Source struct {
	ID     int
	Name   string
	Type   SourceType
	Config SourceConfig
}

type DatasetConfig struct {
	Schema  string   // samples
	Table   string   // country_vaccinations_by_manufacturer
	Columns []string // total_vaccinations::integer
}

type SourceConfig struct {
	DatabaseURI string
	Datasets    []DatasetConfig
}
