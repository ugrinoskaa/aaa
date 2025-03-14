package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/amukoski/aaa/model"
	"text/template"
)

var pieExample = M{
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
			"type":              "pie",
			"radius":            []string{"40%", "70%"},
			"avoidLabelOverlap": false,
			"itemStyle": M{
				"borderRadius": 10,
				"borderColor":  "#fff",
				"borderWidth":  2,
			},
			"data": []M{
				{"value": 150, "name": "Mon"},
				{"value": 230, "name": "Tue"},
				{"value": 224, "name": "Wed"},
				{"value": 218, "name": "Thu"},
				{"value": 135, "name": "Fri"},
				{"value": 147, "name": "Sat"},
				{"value": 260, "name": "Sun"},
			},
		},
	},
}

var pieTemplate = `
{
  "title": {
    "text": "{{.LabelString}}",
    "left": "center",
	"textStyle": {
	  "fontSize":   24,
	  "fontWeight": "bold"
	}
  },
  "series": [
    {
      "type": "pie",
	  "radius": ["40%","70%"],
	  "avoidLabelOverlap": false,
	  "itemStyle": {
		"borderRadius": 10,
		"borderColor":  "#fff",
		"borderWidth":  2
	  },
      "data": {{.ValuesJSONMap}}
    }
  ]
}
`

var pieSchema = model.ChartSchemaRules{
	Dimensions: model.FieldRule{Min: 1, Max: 1, Values: []string{}},
	Metrics:    model.FieldRule{Min: 1, Max: 1, Values: []string{}},
	Filters:    model.FieldRule{Min: 0, Max: 5, Values: model.SupportedFilters},
}

type PieChart struct {
	tmpl *template.Template
}

func NewPieChart() (*PieChart, error) {
	tmpl, err := template.New(string(model.PIE)).Parse(pieTemplate)
	return &PieChart{tmpl: tmpl}, err
}

func (l *PieChart) Schema() model.ChartSchema {
	return model.ChartSchema{
		Type:    model.PIE,
		Schema:  pieSchema,
		Example: pieExample,
	}
}

func (l *PieChart) Render(name string, groups [][]string, values [][]float64) (any, error) {
	req := Model{
		Label:  name,
		XAxis:  groups[0],
		Values: values[0],
	}

	var buf bytes.Buffer
	if err := l.tmpl.Execute(&buf, req); err != nil {
		return nil, fmt.Errorf("failed to render chart renderer: %w", err)
	}

	var data map[string]any
	return data, json.Unmarshal(buf.Bytes(), &data)
}
