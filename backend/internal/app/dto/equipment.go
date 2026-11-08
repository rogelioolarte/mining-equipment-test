package dto

import "github.com/gofrs/uuid/v5"

type Equiptment struct {
	ID                       uuid.UUID `json:"id"`
	Code                     string    `json:"code"`
	Type                     string    `json:"type"`
	CurrentHourMeter         int32     `json:"current_hour_meter"`
	MaintenanceIntervalHours int32     `json:"maintenance_interval_hours"`
	LastMaintenanceHourMeter int32     `json:"last_maintenance_hour_meter"`
	MaintenanceStatus        string    `json:"maintenance_status"`
}
