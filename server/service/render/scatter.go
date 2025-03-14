package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"text/template"

	"github.com/amukoski/aaa/model"
)

var scatterExample = M{
	"title": M{
		"text": "DEMO DATA - CONFIGURE TO PREVIEW",
		"left": "center",
		"textStyle": M{
			"fontSize":   24,
			"fontWeight": "bold",
		},
	},
	"xAxis": M{
		"data": []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"},
	},
	"yAxis": M{},
	"series": []M{
		{
			"type": "scatter",
			"data": []float64{220, 182, 191, 234, 290, 330, 310},
		},
		{
			"type": "scatter",
			"data": []float64{22, 18, 19, 23, 29, 33, 31},
		},
	},
}

var scatterSchema = model.ChartSchemaRules{
	Dimensions: model.FieldRule{Min: 1, Max: 2, Values: []string{}},
	Metrics:    model.FieldRule{Min: 1, Max: 1, Values: []string{}},
	Filters:    model.FieldRule{Min: 0, Max: 5, Values: model.SupportedFilters},
}

var scatterTemplate = `
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

type ScatterChart struct {
	tmpl *template.Template
}

func NewScatterChart() (*ScatterChart, error) {
	tmpl, err := template.New(string(model.SCATTER)).Parse(scatterTemplate)
	return &ScatterChart{tmpl: tmpl}, err
}

func (l *ScatterChart) Schema() model.ChartSchema {
	return model.ChartSchema{
		Type:    model.SCATTER,
		Schema:  scatterSchema,
		Example: scatterExample,
	}
}

func (l *ScatterChart) Render(name string, groups [][]string, values [][]float64) (any, error) {
	req := Model{Label: name}
	is2D := len(groups) == 1
	is3D := len(groups) == 2

	if is2D {
		req.XAxis = groups[0]
		req.Series = []Series{
			{
				Type: string(model.SCATTER),
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
				Type: string(model.SCATTER),
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
