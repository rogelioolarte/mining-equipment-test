package dto

import (
	"minning/internal/domain/entities"

	"github.com/gofrs/uuid/v5"
)

type EquipmentTypeResponse struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	MaintenanceHours int32     `json:"maintenanceHours"`
}

func ToResponseEquipmentType(e *entities.EquipmentType) *EquipmentTypeResponse {
	return &EquipmentTypeResponse{
		ID:               e.ID,
		Name:             e.Name,
		MaintenanceHours: e.MaintenanceHours,
	}
}
