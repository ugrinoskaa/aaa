package api

import (
	"net/http"

	"github.com/amukoski/aaa/service"
	"github.com/gofiber/fiber/v2"
)

type DatasetRsp struct {
	ID         int      `json:"id,omitempty"`
	SourceID   int      `json:"sourceId,omitempty"`
	Name       string   `json:"name,omitempty"`
	Schema     string   `json:"schema,omitempty"`
	Table      string   `json:"table,omitempty"`
	Columns    []string `json:"columns,omitempty"`
	Dimensions []string `json:"dimensions,omitempty"`
	Metrics    []string `json:"metrics,omitempty"`
}

type DatasetAllRsp []DatasetRsp

func (h *Handler) DatasetAll(c *fiber.Ctx) error {
	datasets, err := h.Datasets.All(c.Context())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Error{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	result := make([]DatasetRsp, len(datasets))
	for idx, dataset := range datasets {
		result[idx] = DatasetRsp{
			ID:       dataset.ID,
			SourceID: dataset.SourceID,
			Name:     dataset.Name,
			Schema:   dataset.Config.Schema,
			Table:    dataset.Config.Table,
		}
	}

	return c.JSON(result)
}

func (h *Handler) DatasetGet(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return c.Status(http.StatusBadRequest).JSON(Error{
			Status:  http.StatusBadRequest,
			Message: "invalid dataset id",
		})
	}

	dataset, err := h.Datasets.Get(c.Context(), id)
	if err != nil {
		return c.SendStatus(http.StatusNotFound)
	}

	return c.JSON(DatasetRsp{
		ID:         dataset.ID,
		SourceID:   dataset.SourceID,
		Name:       dataset.Name,
		Schema:     dataset.Config.Schema,
		Table:      dataset.Config.Table,
		Columns:    dataset.Config.Columns,
		Dimensions: dataset.Config.Dimensions(),
		Metrics:    dataset.Config.Metrics(),
	})
}

type DatasetCreateReq struct {
	Name         string `json:"name"`
	SourceID     int    `json:"sourceId"`
	SourceTable  string `json:"sourceTable"`
	SourceSchema string `json:"sourceSchema"`
}

func (h *Handler) DatasetCreate(c *fiber.Ctx) error {
	var req DatasetCreateReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(Error{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
		})
	}

	id, err := h.Datasets.Create(c.Context(), service.CreateDatasetReq{
		Name:           req.Name,
		SourceID:       req.SourceID,
		DatabaseSchema: req.SourceSchema,
		DatabaseTable:  req.SourceTable,
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Error{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.JSON(DatasetRsp{ID: id})
}

func (h *Handler) DatasetDelete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return c.Status(http.StatusBadRequest).JSON(Error{
			Status:  http.StatusBadRequest,
			Message: "invalid id param",
		})
	}

	err = h.Datasets.Delete(c.Context(), id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Error{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.SendStatus(http.StatusNoContent)
}
