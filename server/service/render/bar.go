package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"text/template"

	"github.com/amukoski/aaa/model"
)

var barExample = M{
	"title": M{
		"text": "DEMO DATA - CONFIGURE TO PREVIEW",
		"left": "center",
		"textStyle": M{
			"fontSize":   24,
			"fontWeight": "bold",
		},
	},
	"legend": M{
		"bottom": 0,
	},
	"xAxis": M{
		"data": []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
	},
	"yAxis": M{},
	"series": []M{
		{
			"type": "bar",
			"data": []float64{150, 230, 224, 218, 135, 147, 260},
		},
	},
}

var barSchema = model.ChartSchemaRules{
	Dimensions: model.FieldRule{Min: 1, Max: 2, Values: []string{}},
	Metrics:    model.FieldRule{Min: 1, Max: 1, Values: []string{}},
	Filters:    model.FieldRule{Min: 0, Max: 5, Values: model.SupportedFilters},
}

var barTemplate = `
{
  "title": {
    "text": "{{.LabelString}}",
    "left": "center",
	"textStyle": {
	  "fontSize":   24,
	  "fontWeight": "bold"
	}
  },
  "legend": {
	"bottom": 0,
	"data": {{.LegendJSON}}
  },
  "xAxis": {
    "data": {{.XAxisJSON}}
  },
  "yAxis": {},
  "series": {{.SeriesJSON}}
}
`

type BarChart struct {
	tmpl *template.Template
}

func NewBarChart() (*BarChart, error) {
	tmpl, err := template.New(string(model.BAR)).Parse(barTemplate)
	return &BarChart{tmpl: tmpl}, err
}

func (l *BarChart) Schema() model.ChartSchema {
	return model.ChartSchema{
		Type:    model.BAR,
		Schema:  barSchema,
		Example: barExample,
	}
}

func (l *BarChart) Render(name string, groups [][]string, values [][]float64) (any, error) {
	req := Model{Label: name}
	is2D := len(groups) == 1
	is3D := len(groups) == 2

	if is2D {
		req.XAxis = groups[0]
		req.Series = []Series{
			{
				Type: string(model.BAR),
				Data: values[0],
			},
		}
	}

	if is3D {
		req.XAxis = lo.Uniq(groups[0])
		req.Legend = lo.Uniq(groups[1])
		req.Series = make([]Series, 0)
		aggregate := make(map[string][]float64)

		for idx, group := range groups[1] {
			aggregate[group] = append(aggregate[group], values[0][idx])
		}

		for group, data := range aggregate {
			req.Series = append(req.Series, Series{
				Type: string(model.BAR),
				Name: group,
				Data: data,
			})
		}
	}

	var buf bytes.Buffer
	if err := l.tmpl.Execute(&buf, req); err != nil {
		return nil, fmt.Errorf("failed to render chart renderer: %w", err)
	}

	var data map[string]any
	return data, json.Unmarshal(buf.Bytes(), &data)
}
