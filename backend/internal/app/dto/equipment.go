package dto

import (
	"minning/internal/domain/entities"

	"github.com/gofrs/uuid/v5"
)

type EquipmentRequest struct {
	Code             string    `json:"code"`
	TypeID           uuid.UUID `json:"typeId"`
	CurrentHourMeter int32     `json:"current_hour_meter"`
}

type EquipmentResponse struct {
	ID                       uuid.UUID `json:"id"`
	Code                     string    `json:"code"`
	Type                     *string   `json:"type"`
	CurrentHourMeter         int32     `json:"currentHourMeter"`
	MaintenanceIntervalHours *int32    `json:"maintenanceIntervalHours"`
	HoursSinceMaintenance    int32     `json:"hoursSinceMaintenance"`
	LastMaintenanceHourMeter *int32    `json:"lastMaintenanceHourMeter"`
	MaintenanceStatus        string    `json:"maintenanceStatus"`
}

func (e *EquipmentRequest) ToEntity() (*entities.Equipment, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return &entities.Equipment{
		ID:               id,
		TypeID:           e.TypeID,
		Code:             e.Code,
		CurrentHourMeter: e.CurrentHourMeter,
	}, nil
}

func ToResponseEquipment(e *entities.Equipment) *EquipmentResponse {
	return &EquipmentResponse{
		ID:                       e.ID,
		Code:                     e.Code,
		Type:                     e.EquipmentType,
		CurrentHourMeter:         e.CurrentHourMeter,
		MaintenanceIntervalHours: e.MaintenanceHours,
		HoursSinceMaintenance:    e.HoursSinceMaintenance(),
		LastMaintenanceHourMeter: e.LastMaintenanceHourMeter,
		MaintenanceStatus:        string(e.MaintenanceStatus()),
	}
}
