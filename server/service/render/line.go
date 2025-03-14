package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/samber/lo"

	"github.com/amukoski/aaa/model"
	"text/template"
)

var lineExample = M{
	"title": M{
		"text": "DEMO DATA - CONFIGURE TO PREVIEW",
		"left": "center",
		"textStyle": M{
			"fontSize":   24,
			"fontWeight": "bold",
		},
	},
	"xAxis": M{
		"data": []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
	},
	"yAxis": M{},
	"series": []M{
		{
			"type": "line",
			"data": []float64{150, 230, 224, 218, 135, 147, 260},
		},
	},
}

var lineSchema = model.ChartSchemaRules{
	Dimensions: model.FieldRule{Min: 1, Max: 2, Values: []string{}},
	Metrics:    model.FieldRule{Min: 1, Max: 1, Values: []string{}},
	Filters:    model.FieldRule{Min: 0, Max: 5, Values: model.SupportedFilters},
}

var lineTemplate = `
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

type LineChart struct {
	tmpl *template.Template
}

func NewLineChart() (*LineChart, error) {
	tmpl, err := template.New(string(model.LINE)).Parse(lineTemplate)
	return &LineChart{tmpl: tmpl}, err
}

func (l *LineChart) Schema() model.ChartSchema {
	return model.ChartSchema{
		Type:    model.LINE,
		Schema:  lineSchema,
		Example: lineExample,
	}
}

func (l *LineChart) Render(name string, groups [][]string, values [][]float64) (any, error) {
	req := Model{Label: name}
	is2D := len(groups) == 1
	is3D := len(groups) == 2

	if is2D {
		req.XAxis = groups[0]
		req.Series = []Series{
			{
				Type: string(model.LINE),
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
				Type: string(model.LINE),
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
