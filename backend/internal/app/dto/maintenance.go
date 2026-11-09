package dto

import (
	"minning/internal/domain/entities"
	"time"

	"github.com/gofrs/uuid/v5"
)

type MaintenanceRequest struct {
	EquipmentID            uuid.UUID `json:"equipmentId"`
	HourMeterAtMaintenance int32     `json:"hmMaintenance"`
	MaintenanceDate        time.Time `json:"maintenanceDate"`
	Description            *string   `json:"description,omitempty"`
}

type MaintenanceResponse struct {
	ID                     uuid.UUID `json:"id"`
	EquipmentID            uuid.UUID `json:"equipmentId"`
	EquipmentCode          *string   `json:"equipmentCode,omitempty"`
	EquipmentType          *string   `json:"equipmentType,omitempty"`
	HourMeterAtMaintenance int32     `json:"hmMaintenance"`
	MaintenanceDate        time.Time `json:"maintenanceDate"`
	Description            *string   `json:"description,omitempty"`
}

func (r *MaintenanceRequest) ToEntity() (*entities.Maintenance, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return &entities.Maintenance{
		ID:                     id,
		EquipmentID:            r.EquipmentID,
		HourMeterAtMaintenance: r.HourMeterAtMaintenance,
		MaintenanceDate:        r.MaintenanceDate,
		Description:            r.Description,
	}, nil
}

func ToResponseMaintenance(e *entities.Maintenance) *MaintenanceResponse {
	return &MaintenanceResponse{
		ID:                     e.ID,
		EquipmentID:            e.EquipmentID,
		EquipmentCode:          e.EquipmentCode,
		EquipmentType:          e.EquipmentType,
		HourMeterAtMaintenance: e.HourMeterAtMaintenance,
		MaintenanceDate:        e.MaintenanceDate,
		Description:            e.Description,
	}
}
