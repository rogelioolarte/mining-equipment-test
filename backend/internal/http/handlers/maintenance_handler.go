package handlers

import (
	"minning/internal/app/usecases"
	"minning/internal/infra/logger"
	"net/http"
)

type MaintenanceHandler struct {
	logger             *logger.Logger
	maintenanceUseCase *usecases.MaintenanceUseCase
}

func NewnMaintenanceHandler(
	logger *logger.Logger,
	maintenanceUseCase *usecases.MaintenanceUseCase,
) *MaintenanceHandler {
	return &MaintenanceHandler{
		logger:             logger,
		maintenanceUseCase: maintenanceUseCase,
	}
}

func (h *MaintenanceHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Fetching all Maintenances")

	res, err := h.maintenanceUseCase.GetAll(r.Context())
	if err != nil {
		h.logger.Error("Failed while getting entities")
		InternalServerError(w, err.Error())
		return
	}

	Success(w, res)
}
