package render

import (
	"encoding/json"
	"github.com/samber/lo"
	"slices"
)

type M map[string]interface{}

type Model struct {
	Label  string
	Legend []string
	XAxis  []string
	YAxis  []string
	Values []float64
	Series []Series
}

type Series struct {
	Type string    `json:"type"`
	Name string    `json:"name"`
	Data []float64 `json:"data"`
}

func (m Model) LegendJSON() string {
	rsp, _ := json.Marshal(m.Legend)
	return string(rsp)
}

func (m Model) XAxisJSON() string {
	rsp, _ := json.Marshal(m.XAxis)
	return string(rsp)
}

func (m Model) YAxisJSON() string {
	rsp, _ := json.Marshal(m.YAxis)
	return string(rsp)
}

func (m Model) XAxisUniqueJSON() string {
	rsp, _ := json.Marshal(lo.Uniq(m.XAxis))
	return string(rsp)
}

func (m Model) YAxisUniqueJSON() string {
	rsp, _ := json.Marshal(lo.Uniq(m.YAxis))
	return string(rsp)
}

func (m Model) ValuesJSON() string {
	rsp, _ := json.Marshal(m.Values)
	return string(rsp)
}

func (m Model) ValuesJSONMap() string {
	result := make([]map[string]any, len(m.XAxis))
	for idx := range len(m.XAxis) {
		result[idx] = map[string]any{"name": m.XAxis[idx], "value": m.Values[idx]}
	}

	rsp, _ := json.Marshal(result)
	return string(rsp)
}

func (m Model) SeriesJSON() string {
	rsp, _ := json.Marshal(m.Series)
	return string(rsp)
}

func (m Model) ValuesTupleJSON() string {
	xAxis := lo.Uniq(m.XAxis)
	yAxis := lo.Uniq(m.YAxis)

	result := make([][]float64, 0)
	for i := range len(m.XAxis) {
		x, y, z := m.XAxis[i], m.YAxis[i], m.Values[i]
		x1, y1 := slices.Index(xAxis, x), slices.Index(yAxis, y)
		tuple := []float64{float64(x1), float64(y1), z}
		result = append(result, tuple)
	}

	rsp, _ := json.Marshal(result)
	return string(rsp)
}

func (m Model) ValuesMin() float64 {
	return slices.Min(m.Values)
}

func (m Model) ValuesMax() float64 {
	return slices.Max(m.Values)
}

func (m Model) LabelString() string {
	return m.Label
}

func (m Model) LinkCategoriesJSON() string {
	type Record struct {
		Name string `json:"name"`
	}

	unique := make(map[string]bool)
	for _, x := range m.XAxis {
		unique[x] = true
	}
	for _, y := range m.YAxis {
		unique[y] = true
	}

	result := make([]Record, 0)
	for key := range unique {
		result = append(result, Record{Name: key})
	}

	if len(result) > 0 {
		result = append(result, Record{Name: ""})
	}

	rsp, _ := json.Marshal(result)
	return string(rsp)
}

func (m Model) LinkValuesJSON() string {
	type Record struct {
		Source string  `json:"source"`
		Target string  `json:"target"`
		Value  float64 `json:"value"`
	}

	result := make([]Record, 0)
	if !slices.Equal(m.XAxis, m.YAxis) {
		for idx, y := range m.YAxis {
			result = append(result, Record{Source: y, Target: m.XAxis[idx], Value: m.Values[idx]})
		}

		for _, x := range m.XAxis {
			var subtotal float64
			for _, r := range result {
				if x == r.Target {
					subtotal += r.Value
				}
			}

			result = append(result, Record{Source: x, Target: "", Value: subtotal})
		}
	} else {
		for idx, x := range m.XAxis {
			result = append(result, Record{Source: x, Target: "", Value: m.Values[idx]})
		}
	}

	rsp, _ := json.Marshal(result)
	return string(rsp)
}
