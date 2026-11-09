package handlers

import (
	"minning/internal/app/usecases"
	"minning/internal/infra/logger"
	"net/http"
)

type EquipmentTypeHandler struct {
	logger               *logger.Logger
	equipmentTypeUseCase *usecases.EquipmentTypeUseCase
}

func NewnEquipmentTypeHandler(
	logger *logger.Logger,
	equipmentTypeUseCase *usecases.EquipmentTypeUseCase,
) *EquipmentTypeHandler {
	return &EquipmentTypeHandler{
		logger:               logger,
		equipmentTypeUseCase: equipmentTypeUseCase,
	}
}

func (h *EquipmentTypeHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Fetching All EquipmentTypes")

	res, err := h.equipmentTypeUseCase.GetAll(r.Context())
	if err != nil {
		h.logger.Error("Failed while getting EquipmentTypes")
		InternalServerError(w, err.Error())
		return
	}

	Success(w, res)
}
