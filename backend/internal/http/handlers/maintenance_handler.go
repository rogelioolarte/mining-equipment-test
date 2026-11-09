package handlers

import (
	"minning/internal/app/dto"
	"minning/internal/app/usecases"
	"minning/internal/domain/entities"
	"minning/internal/infra/logger"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
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

func (h *MaintenanceHandler) GetAllByEquipmentID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		BadRequest(w, "ID is required")
		return
	}
	ID, idErr := uuid.FromString(idStr)
	if idErr != nil {
		BadRequest(w, "Invalid ID")
		return
	}

	h.logger.Info("Fetching all Maintenances with Equipment ID", "ID", ID)

	res, err := h.maintenanceUseCase.GetAllByEquipmentID(r.Context(), ID)
	if err != nil {
		h.logger.Error("Failed while getting Maintenances")
		if err == entities.ErrNotFoundEntity {
			NotFound(w, "Maintenances Not Found")
			return
		}
		InternalServerError(w, err.Error())
		return
	}

	Success(w, res)
}

func (h *MaintenanceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.MaintenanceRequest
	if err := DecodeJSON(w, r, &req); err != nil {
		h.logger.Warn("Invalid Json Body for Create Maintenance")
		BadRequest(w, err.Error())
		return
	}

	res, err := h.maintenanceUseCase.Create(r.Context(), &req)
	if err != nil {
		h.logger.Error("Error while creating Maintenance")
		BadRequest(w, err.Error())
		return
	}
	Created(w, res)
}
