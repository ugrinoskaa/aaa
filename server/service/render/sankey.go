package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/amukoski/aaa/model"
)

var sankeyExample = M{
	"title": M{
		"text": "DEMO DATA - CONFIGURE TO PREVIEW",
		"left": "center",
		"textStyle": M{
			"fontSize":   24,
			"fontWeight": "bold",
		},
	},
	"series": []M{
		{
			"type":      "sankey",
			"draggable": false,
			"left":      "10%",
			"top":       "10%",
			"right":     "10%",
			"bottom":    "10%",
			"data": []M{
				{"name": "Online"},
				{"name": "In-Store"},
				{"name": "Monday"},
				{"name": "Tuesday"},
				{"name": "Wednesday"},
				{"name": "Thursday"},
				{"name": "Friday"},
				{"name": ""},
			},
			"links": []M{
				{
					"source": "",
					"target": "Online",
					"value":  18,
				},
				{
					"source": "",
					"target": "In-Store",
					"value":  13,
				},
				{
					"source": "Online",
					"target": "Monday",
					"value":  8,
				},
				{
					"source": "In-Store",
					"target": "Monday",
					"value":  2,
				},
				{
					"source": "Online",
					"target": "Tuesday",
					"value":  2,
				},
				{
					"source": "Online",
					"target": "Wednesday",
					"value":  3,
				},
				{
					"source": "Online",
					"target": "Thursday",
					"value":  5,
				},
				{
					"source": "In-Store",
					"target": "Thursday",
					"value":  10,
				},
				{
					"source": "In-Store",
					"target": "Friday",
					"value":  1,
				},
			},
		},
	},
}

var sankeySchema = model.ChartSchemaRules{
	Dimensions: model.FieldRule{Min: 1, Max: 2, Values: []string{}},
	Metrics:    model.FieldRule{Min: 1, Max: 1, Values: []string{}},
	Filters:    model.FieldRule{Min: 0, Max: 5, Values: model.SupportedFilters},
}

var sankeyTemplate = `
{
  "title": {
    "text": "{{.LabelString}}",
    "left": "center",
	"textStyle": {
	  "fontSize":   24,
	  "fontWeight": "bold"
	}
  },
  "series": {
    "type": "sankey",
    "draggable": false,
	"left": "5%",
	"top": "5%",
	"right": "5%",
	"bottom": "5%",
    "data": {{.LinkCategoriesJSON}},
    "links": {{.LinkValuesJSON}}
  }
}
`

type SankeyChart struct {
	tmpl *template.Template
}

func NewSankeyChart() (*SankeyChart, error) {
	tmpl, err := template.New(string(model.SANKEY)).Parse(sankeyTemplate)
	return &SankeyChart{tmpl: tmpl}, err
}

func (l *SankeyChart) Schema() model.ChartSchema {
	return model.ChartSchema{
		Type:    model.SANKEY,
		Schema:  sankeySchema,
		Example: sankeyExample,
	}
}

func (l *SankeyChart) Render(name string, groups [][]string, values [][]float64) (any, error) {
	req := Model{
		Label:  name,
		XAxis:  groups[0],
		YAxis:  groups[0],
		Values: values[0],
	}

	is3D := len(groups) == 2
	if is3D {
		req.YAxis = groups[1]
	}

	var buf bytes.Buffer
	if err := l.tmpl.Execute(&buf, req); err != nil {
		return nil, fmt.Errorf("failed to render chart renderer: %w", err)
	}

	var data map[string]any
	return data, json.Unmarshal(buf.Bytes(), &data)
}
