package model

type Dashboard struct {
	ID   int              `json:"id"`
	Name string           `json:"name"`
	Grid []map[string]any `json:"grid"`
}
