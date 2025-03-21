package api

import (
	"net/http"

	"github.com/amukoski/aaa/service"
	"github.com/gofiber/fiber/v2"
)

type ChartRsp struct {
	ID         int      `json:"id,omitempty"`
	DatasetID  int      `json:"datasetId,omitempty"`
	Name       string   `json:"name,omitempty"`
	Type       string   `json:"type,omitempty"`
	Dimensions []string `json:"dimensions,omitempty"`
	Metrics    []string `json:"metrics,omitempty"`
	Filters    []string `json:"filters,omitempty"`
}

type ChartAllRsp []ChartRsp

func (h *Handler) ChartAll(c *fiber.Ctx) error {
	charts, err := h.Charts.All(c.Context())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Error{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	result := make([]ChartRsp, len(charts))
	for idx, chart := range charts {
		result[idx] = ChartRsp{
			ID:         chart.ID,
			DatasetID:  chart.DatasetID,
			Name:       chart.Name,
			Type:       string(chart.Type),
			Dimensions: chart.Config.Dimensions,
			Metrics:    chart.Config.Metrics,
			Filters:    chart.Config.Filters,
		}
	}

	return c.JSON(result)
}

func (h *Handler) ChartGet(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return c.Status(http.StatusBadRequest).JSON(Error{
			Status:  http.StatusBadRequest,
			Message: "invalid chart id",
		})
	}

	chart, err := h.Charts.Get(c.Context(), id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Error{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.JSON(ChartRsp{
		ID:         chart.ID,
		DatasetID:  chart.DatasetID,
		Name:       chart.Name,
		Type:       string(chart.Type),
		Dimensions: chart.Config.Dimensions,
		Metrics:    chart.Config.Metrics,
		Filters:    chart.Config.Filters,
	})
}

type CreateChartReq struct {
	DatasetID  int      `json:"datasetId"`
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Dimensions []string `json:"dimensions"`
	Metrics    []string `json:"metrics"`
	Filters    []string `json:"filters"`
}

func (h *Handler) ChartCreate(c *fiber.Ctx) error {
	var req CreateChartReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(Error{
			Status:  http.StatusBadRequest,
			Message: "invalid request",
		})
	}

	id, err := h.Charts.Create(c.Context(), service.CreateChartReq{
		DatasetID:  req.DatasetID,
		Name:       req.Name,
		Type:       req.Type,
		Dimensions: req.Dimensions,
		Metrics:    req.Metrics,
		Filters:    req.Filters,
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Error{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.JSON(ChartRsp{ID: id})
}

func (h *Handler) ChartDelete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return c.Status(http.StatusBadRequest).JSON(Error{
			Status:  http.StatusBadRequest,
			Message: "invalid request",
		})
	}

	err = h.Charts.Delete(c.Context(), id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Error{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.SendStatus(http.StatusNoContent)
}

type ValidateChartReq struct {
	DatasetID  int      `json:"datasetId"`
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Dimensions []string `json:"dimensions"`
	Metrics    []string `json:"metrics"`
	Filters    []string `json:"filters"`
}

type ValidateChartRsp struct {
	Valid   bool `json:"valid,omitempty"`
	Options any  `json:"options,omitempty"`
}

func (h *Handler) ChartValidate(c *fiber.Ctx) error {
	var req ValidateChartReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(Error{
			Status:  http.StatusBadRequest,
			Message: "invalid request",
		})
	}

	result, err := h.Charts.Validate(c.Context(), service.ValidateChartReq{
		DatasetID:  req.DatasetID,
		Name:       req.Name,
		Type:       req.Type,
		Dimensions: req.Dimensions,
		Metrics:    req.Metrics,
		Filters:    req.Filters,
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Error{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.JSON(ValidateChartRsp{Valid: true, Options: result})
}

func (h *Handler) ChartRun(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return c.Status(http.StatusBadRequest).JSON(Error{
			Status:  http.StatusBadRequest,
			Message: "invalid chart id",
		})
	}

	result, err := h.Charts.Run(c.Context(), id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Error{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.JSON(ValidateChartRsp{Valid: true, Options: result})
}

type ChartTypesAll []string

func (h *Handler) ChartTypesAll(c *fiber.Ctx) error {
	charts := h.Charts.AllTypes()

	result := make(ChartTypesAll, len(charts))
	for idx, chart := range charts {
		result[idx] = string(chart)
	}

	return c.JSON(result)
}

type ChartSchemaRsp struct {
	Type    string      `json:"type,omitempty"`
	Schema  interface{} `json:"schema,omitempty"`
	Example interface{} `json:"example,omitempty"`
}

func (h *Handler) ChartGetType(c *fiber.Ctx) error {
	ctype := c.Params("type")
	if ctype == "" {
		return c.Status(http.StatusBadRequest).JSON(Error{
			Status:  http.StatusBadRequest,
			Message: "invalid chart type",
		})
	}

	chart, err := h.Charts.GetType(ctype)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(Error{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return c.JSON(ChartSchemaRsp{
		Type:    string(chart.Type),
		Schema:  chart.Schema,
		Example: chart.Example,
	})
}
