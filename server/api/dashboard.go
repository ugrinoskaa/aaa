package api

import (
	"net/http"

	"github.com/amukoski/aaa/service"
	"github.com/gofiber/fiber/v2"
)

type DashboardRsp struct {
	ID   int              `json:"id,omitempty"`
	Name string           `json:"name,omitempty"`
	Grid []map[string]any `json:"grid,omitempty"`
}

type DashboardAllRsp []DashboardRsp

func (h *Handler) DashboardAll(c *fiber.Ctx) error {
	dashboards, err := h.Dashboard.All(c.Context())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Error{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	result := make([]DashboardRsp, len(dashboards))
	for idx, dashboard := range dashboards {
		result[idx] = DashboardRsp{
			ID:   dashboard.ID,
			Name: dashboard.Name,
			Grid: dashboard.Grid,
		}
	}

	return c.JSON(result)
}

func (h *Handler) DashboardGet(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return c.Status(http.StatusBadRequest).JSON(Error{
			Status:  http.StatusBadRequest,
			Message: "invalid dashboard id",
		})
	}

	dashboard, err := h.Dashboard.Get(c.Context(), id)
	if err != nil {
		return c.SendStatus(http.StatusNotFound)
	}

	return c.JSON(DashboardRsp{
		ID:   dashboard.ID,
		Name: dashboard.Name,
		Grid: dashboard.Grid,
	})
}

type DashboardCreateReq struct {
	Name string           `json:"name"`
	Grid []map[string]any `json:"grid"`
}

func (h *Handler) DashboardCreate(c *fiber.Ctx) error {
	var req DashboardCreateReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(Error{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
		})
	}

	id, err := h.Dashboard.Create(c.Context(), service.CreateDashboardReq{
		Name: req.Name,
		Grid: req.Grid,
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Error{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.JSON(DashboardRsp{ID: id})
}

func (h *Handler) DashboardDelete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return c.Status(http.StatusBadRequest).JSON(Error{
			Status:  http.StatusBadRequest,
			Message: "invalid id param",
		})
	}

	err = h.Dashboard.Delete(c.Context(), id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Error{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.SendStatus(http.StatusNoContent)
}
