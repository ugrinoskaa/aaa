package api

import (
	"log"

	"github.com/amukoski/aaa/service"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Logger    *log.Logger
	Sources   *service.SourceService
	Datasets  *service.DatasetService
	Charts    *service.ChartService
	Dashboard *service.DashboardService
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Get("/sources", h.SourceAll)
	router.Get("/sources/:id", h.SourceGet)
	router.Post("/sources", h.SourceCreate)
	router.Post("/sources/discovery", h.SourceDiscovery)
	router.Delete("/sources/:id", h.SourceDelete)

	router.Get("/datasets", h.DatasetAll)
	router.Get("/datasets/:id", h.DatasetGet)
	router.Post("/datasets", h.DatasetCreate)
	router.Delete("/datasets/:id", h.DatasetDelete)

	router.Get("/chart-types", h.ChartTypesAll)
	router.Get("/chart-types/:type", h.ChartGetType)

	router.Get("/charts", h.ChartAll)
	router.Get("/charts/:id", h.ChartGet)
	router.Post("/charts", h.ChartCreate)
	router.Post("/charts/validate", h.ChartValidate)
	router.Post("/charts/:id/data", h.ChartRun)
	router.Delete("/charts/:id", h.ChartDelete)

	router.Get("/dashboards", h.DashboardAll)
	router.Get("/dashboards/:id", h.DashboardGet)
	router.Post("/dashboards", h.DashboardCreate)
	router.Delete("/dashboards/:id", h.DashboardDelete)
}

type Error struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
