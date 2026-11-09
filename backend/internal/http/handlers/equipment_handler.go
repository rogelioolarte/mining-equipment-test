package handlers

import (
	"minning/internal/app/usecases"
	"minning/internal/infra/logger"
	"net/http"
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
		h.logger.Error("Failed while getting entities")
		InternalServerError(w, err.Error())
		return
	}

	Success(w, res)
}
