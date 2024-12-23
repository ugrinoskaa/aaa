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
}

type Error struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
