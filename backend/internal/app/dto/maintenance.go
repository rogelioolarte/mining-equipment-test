package dto

import (
	"minning/internal/domain/entities"
	"time"

	"github.com/gofrs/uuid/v5"
)

type MaintenanceRequest struct {
	EquipmentID            uuid.UUID `json:"equipment_id"`
	HourMeterAtMaintenance int32     `json:"hm_maintenance"`
	MaintenanceDate        time.Time `json:"maintenance_date"`
	Description            *string   `json:"description,omitempty"`
}

type MaintenanceResponse struct {
	ID                     uuid.UUID `json:"id"`
	EquipmentID            uuid.UUID `json:"equipment_id"`
	EquipmentCode          *string   `json:"equipment_code,omitempty"`
	EquipmentType          *string   `json:"equipment_type,omitempty"`
	HourMeterAtMaintenance int32     `json:"hm_maintenance"`
	MaintenanceDate        time.Time `json:"maintenance_date"`
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
