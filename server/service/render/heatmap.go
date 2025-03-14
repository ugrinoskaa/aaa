package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/amukoski/aaa/model"
	"text/template"
)

var heatmapExample = M{
	"title": M{
		"text": "DEMO DATA - CONFIGURE TO PREVIEW",
		"left": "center",
		"textStyle": M{
			"fontSize":   24,
			"fontWeight": "bold",
		},
	},
	"grid": M{
		"height": "70%",
	},
	"xAxis": M{
		"data": []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
	},
	"yAxis": M{
		"data": []string{"Morning", "Afternoon", "Evening"},
	},
	"visualMap": M{
		"min":        20,
		"max":        40,
		"calculable": true,
		"orient":     "horizontal",
		"left":       "center",
	},
	"series": []M{
		{
			"type": "heatmap",
			"label": M{
				"show": true,
			},
			"itemStyle": M{
				"shadowColor": "rgba(0, 0, 0, 0.5)",
				"shadowBlur":  2,
			},
			"data": [][]float64{
				// Morning
				{0, 0, 20}, {1, 0, 22}, {2, 0, 21}, {3, 0, 23}, {4, 0, 24}, {5, 0, 20}, {6, 0, 19},
				// Afternoon
				{0, 1, 30}, {1, 1, 32}, {2, 1, 33}, {3, 1, 31}, {4, 1, 34}, {5, 1, 35}, {6, 1, 29},
				// Evening
				{0, 2, 24}, {1, 2, 26}, {2, 2, 25}, {3, 2, 22}, {4, 2, 23}, {5, 2, 27}, {6, 2, 21},
			},
		},
	},
}

var heatmapSchema = model.ChartSchemaRules{
	Dimensions: model.FieldRule{Min: 1, Max: 2, Values: []string{}},
	Metrics:    model.FieldRule{Min: 1, Max: 1, Values: []string{}},
	Filters:    model.FieldRule{Min: 0, Max: 5, Values: model.SupportedFilters},
}

var heatmapTemplate = `
{
  "title": {
    "text": "{{.LabelString}}",
    "left": "center",
	"textStyle": {
	  "fontSize":   24,
	  "fontWeight": "bold"
	}
  },
  "grid": {
    "height": "70%"
  },
  "xAxis": {
    "data": {{.XAxisUniqueJSON}}
  },
  "yAxis": {
    "data": {{.YAxisUniqueJSON}}
  },
  "visualMap": {
    "min": {{.ValuesMin}},
    "max": {{.ValuesMax}},
    "calculable": true,
    "orient": "horizontal",
    "left": "center"
  },
  "series": [
    {
      "type": "heatmap",
      "label": {
        "show": true
      },
	  "itemStyle": {
	    "shadowBlur": 2,
	    "shadowColor": "rgba(0, 0, 0, 0.5)"
	  },
      "data": {{.ValuesTupleJSON}}
    }
  ]
}
`

type HeatmapChart struct {
	tmpl *template.Template
}

func NewHeatmapChart() (*HeatmapChart, error) {
	tmpl, err := template.New(string(model.HEATMAP)).Parse(heatmapTemplate)
	return &HeatmapChart{tmpl: tmpl}, err
}

func (l *HeatmapChart) Schema() model.ChartSchema {
	return model.ChartSchema{
		Type:    model.HEATMAP,
		Schema:  heatmapSchema,
		Example: heatmapExample,
	}
}

func (l *HeatmapChart) Render(name string, groups [][]string, values [][]float64) (any, error) {
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
