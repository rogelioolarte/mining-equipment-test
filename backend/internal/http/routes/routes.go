package routes

import (
	"minning/internal/http/handlers"
	"minning/internal/infra/logger"

	"github.com/go-chi/chi"
)

func SetupRoutes(
	logger *logger.Logger,
	router *chi.Mux,
	equipmentTypeHandler *handlers.EquipmentTypeHandler,
	equipmentHandler *handlers.EquipmentHandler,
	maintenanceHanlder *handlers.MaintenanceHandler,
) {
	router.Route("/api/v1/", func(r chi.Router) {
		r.Group(func(r chi.Router) {

			r.Route("/equipment_type", func(r chi.Router) {
				r.Get("/", equipmentTypeHandler.GetAll)
			})

			r.Route("/maintenance", func(r chi.Router) {
				r.Get("/", maintenanceHanlder.GetAll)
			})

			r.Route("/equipment", func(r chi.Router) {
				r.Get("/", equipmentHandler.GetAll)
				/* r.Post("/", h.Organization.Create)
				r.Put("/{id}", h.Organization.UpdateById)
				r.Delete("/{id}", h.Organization.DeleteById)
				r.Get("/{id}", h.Organization.GetById) */
			})

		})
	})
}
