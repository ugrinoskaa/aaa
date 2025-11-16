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

type SourceConfig struct {
	DatabaseURI string
	Datasets    []DatasetConfig
}
