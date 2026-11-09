package dto

import (
	"minning/internal/domain/entities"
	"time"

	"github.com/gofrs/uuid/v5"
)

type MaintenanceResponse struct {
	ID                     uuid.UUID `json:"id"`
	EquipmentID            uuid.UUID `json:"equipment_id"`
	EquipmentCode          *string   `json:"equipment_code"`
	EquipmentType          *string   `json:"equipment_type"`
	HourMeterAtMaintenance int32     `json:"hm_maintenance"`
	MaintenanceDate        time.Time `json:"maintenance_date"`
	Description            *string   `json:"description"`
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
