package api

import (
	"errors"
	"net/http"

	"github.com/amukoski/aaa/model"
	"github.com/amukoski/aaa/service"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"
)

type SourceRsp struct {
	ID       int          `json:"id,omitempty"`
	Name     string       `json:"name,omitempty"`
	Type     string       `json:"type,omitempty"`
	Datasets []DatasetRsp `json:"datasets,omitempty"`
}

type SourceAllRsp []SourceRsp

func (h *Handler) SourceAll(c *fiber.Ctx) error {
	sources, err := h.Sources.All(c.Context())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Error{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	result := make([]SourceRsp, len(sources))
	for idx, source := range sources {
		result[idx] = SourceRsp{
			ID:   source.ID,
			Name: source.Name,
			Type: string(source.Type),
		}
	}

	return c.JSON(result)
}

func (h *Handler) SourceGet(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return c.Status(http.StatusBadRequest).JSON(Error{
			Status:  http.StatusBadRequest,
			Message: "invalid source id",
		})
	}

	source, datasets, err := h.Sources.Get(c.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.SendStatus(http.StatusNotFound)
		}

		return c.Status(http.StatusInternalServerError).JSON(Error{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	result := SourceRsp{
		ID:       source.ID,
		Name:     source.Name,
		Type:     string(source.Type),
		Datasets: make([]DatasetRsp, len(datasets)),
	}

	for idx, ds := range datasets {
		result.Datasets[idx] = DatasetRsp{
			Schema:  ds.Schema,
			Table:   ds.Table,
			Columns: ds.Columns,
		}
	}

	return c.JSON(result)
}

type SourceDiscoveryReq struct {
	Type string `json:"type"`
	URI  string `json:"uri"`
}

func (h *Handler) SourceDiscovery(c *fiber.Ctx) error {
	var req SourceDiscoveryReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(Error{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if req.Type == string(model.POSTGRES) {
		return h.SourcePostgresConnect(c, req.URI)
	}

	if req.Type == string(model.CSV) {
		return h.SourceCSVUpload(c)
	}

	return c.Status(http.StatusBadRequest).JSON(Error{
		Status:  http.StatusBadRequest,
		Message: "invalid type",
	})
}

func (h *Handler) SourcePostgresConnect(c *fiber.Ctx, uri string) error {
	datasets, err := h.Sources.DiscoverDB(c.Context(), uri)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Error{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	items := make([]DatasetRsp, len(datasets))
	result := SourceRsp{Type: string(model.POSTGRES), Datasets: items}

	for idx, ds := range datasets {
		result.Datasets[idx] = DatasetRsp{
			Table:   ds.Config.Table,
			Schema:  ds.Config.Schema,
			Columns: ds.Config.Columns,
		}
	}

	return c.JSON(result)
}

func (h *Handler) SourceCSVUpload(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(Error{
			Status:  http.StatusBadRequest,
			Message: "invalid multipart form",
		})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(http.StatusBadRequest).JSON(Error{
			Status:  http.StatusBadRequest,
			Message: "no files uploaded",
		})
	}

	datasets, err := h.Sources.Upload(c.Context(), files)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Error{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	items := make([]DatasetRsp, len(datasets))
	result := SourceRsp{Type: string(model.CSV), Datasets: items}

	for idx, ds := range datasets {
		result.Datasets[idx] = DatasetRsp{
			Table:   ds.Config.Table,
			Schema:  ds.Config.Schema,
			Columns: ds.Config.Columns,
		}
	}

	return c.JSON(result)
}

type SourceCreateReq struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Resource string `json:"resource"`
}

func (h *Handler) SourceCreate(c *fiber.Ctx) error {
	var req SourceCreateReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(Error{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	id, err := h.Sources.Create(c.Context(), service.CreateSourceReq{
		Name:     req.Name,
		Type:     req.Type,
		Resource: req.Resource,
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Error{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.JSON(SourceRsp{ID: id})
}

func (h *Handler) SourceDelete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return c.Status(http.StatusBadRequest).JSON(Error{
			Status:  http.StatusBadRequest,
			Message: "invalid id param",
		})
	}

	err = h.Sources.Delete(c.Context(), id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Error{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.SendStatus(http.StatusNoContent)
}
