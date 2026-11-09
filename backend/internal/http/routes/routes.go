package routes

import (
	"minning/internal/http/handlers"
	"minning/internal/infra/logger"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(
	logger *logger.Logger,
	router *chi.Mux,
	healthHandler *handlers.HealthHandler,
	equipmentTypeHandler *handlers.EquipmentTypeHandler,
	equipmentHandler *handlers.EquipmentHandler,
	maintenanceHanlder *handlers.MaintenanceHandler,
) {

	router.Route("/api", func(r chi.Router) {
		r.Route("/health", healthHandler.RegisterRoutes)
		r.Route("/v1", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Route("/equipment_type", func(r chi.Router) {
					r.Get("/", equipmentTypeHandler.GetAll)
				})

				r.Route("/maintenance", func(r chi.Router) {
					r.Post("/", maintenanceHanlder.Create)
				})

				r.Route("/equipment", func(r chi.Router) {
					r.Get("/", equipmentHandler.GetAll)
					r.Get("/{id}/maintenance", maintenanceHanlder.GetAllByEquipmentID)
					r.Post("/", equipmentHandler.Create)
					r.Put("/{id}", equipmentHandler.UpdateById)
					r.Delete("/{id}", equipmentHandler.DeleteById)
				})
			})
		})
	})

	router.Route("/api", func(r chi.Router) {
		r.Route("/health", healthHandler.RegisterRoutes)
		r.Route("/v1", func(r chi.Router) {

			r.Route("/equipment_type", func(r chi.Router) {
				r.Get("/", equipmentTypeHandler.GetAll)
			})

			r.Route("/maintenance", func(r chi.Router) {
				r.Post("/", maintenanceHanlder.Create)
			})

			r.Route("/equipment", func(r chi.Router) {
				r.Get("/", equipmentHandler.GetAll)
				r.Get("/{id}/maintenance", maintenanceHanlder.GetAllByEquipmentID)
				r.Post("/", equipmentHandler.Create)
				r.Put("/{id}", equipmentHandler.UpdateById)
				r.Delete("/{id}", equipmentHandler.DeleteById)
			})

		})
	})
}
