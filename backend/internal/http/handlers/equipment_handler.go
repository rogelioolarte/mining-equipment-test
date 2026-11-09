package handlers

import (
	"minning/internal/app/dto"
	"minning/internal/app/usecases"
	"minning/internal/infra/logger"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid/v5"
)

type EquipmentHandler struct {
	logger           *logger.Logger
	equipmentUseCase *usecases.EquipmentUseCase
}

func NewEquipmentHandler(
	logger *logger.Logger,
	equipmentUseCase *usecases.EquipmentUseCase,
) *EquipmentHandler {
	return &EquipmentHandler{
		logger:           logger,
		equipmentUseCase: equipmentUseCase,
	}
}

func (h *EquipmentHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Fetching all Equipments")

	res, err := h.equipmentUseCase.GetAll(r.Context())
	if err != nil {
		h.logger.Error("Failed while getting Equipments")
		InternalServerError(w, err.Error())
		return
	}

	Success(w, res)
}

func (h *EquipmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.EquipmentRequest
	if err := DecodeJSON(w, r, &req); err != nil {
		h.logger.Warn("Invalid Json Body for Create Equipment")
		BadRequest(w, err.Error())
		return
	}

	res, err := h.equipmentUseCase.Create(r.Context(), &req)
	if err != nil {
		h.logger.Error("Error while creating Equipment")
		BadRequest(w, err.Error())
		return
	}
	Created(w, res)
}

func (h *EquipmentHandler) UpdateById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		h.logger.Warn("ID param is empty")
		BadRequest(w, "ID is required")
		return
	}
	ID, err := uuid.FromString(idStr)
	if err != nil {
		BadRequest(w, "Invalid ID")
		return
	}

	var req dto.EquipmentRequest
	if err := DecodeJSON(w, r, &req); err != nil {
		h.logger.Warn("Invalid Json Body for Update Equipment")
		BadRequest(w, err.Error())
		return
	}

	h.logger.Info("Updating Equipment with ID", "ID", ID)

	res, err := h.equipmentUseCase.UpdateByID(r.Context(), ID, &req)
	if err != nil {
		h.logger.Error("Error while updating Equipment")
		BadRequest(w, err.Error())
		return
	}
	Success(w, res)
}

func (h *EquipmentHandler) DeleteById(w http.ResponseWriter, r *http.Request) {

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

	h.logger.Info("Deleting Equipment with ID", "ID", ID)

	err := h.equipmentUseCase.DeleteByID(r.Context(), ID)
	if err != nil {
		h.logger.Error("Error while deleting Equipment")
		InternalServerError(w, err.Error())
		return
	}
	NoContent(w)
}
