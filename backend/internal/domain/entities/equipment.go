package entities

import "github.com/gofrs/uuid/v5"

type MaintenanceStatus string

const (
	MaintenanceStatusOperational MaintenanceStatus = "OPERATIONAL"
	MaintenanceStatusRequired    MaintenanceStatus = "REQUIRES_MAINTENANCE"
	MaintenanceStatusUnkown      MaintenanceStatus = "UNKOWN"
)

type Equipment struct {
	ID                       uuid.UUID `db:"id"`
	Code                     string    `db:"code"`
	TypeID                   uuid.UUID `db:"type_id"`
	EquipmentType            *string   `db:"equipment_type"`
	CurrentHourMeter         int32     `db:"current_h_meter"`
	MaintenanceHours         *int32    `db:"maintenance_hours"`
	LastMaintenanceHourMeter *int32    `db:"last_mh_meter"`
}

func (e *Equipment) MaintenanceStatus() MaintenanceStatus {
	if e == nil || e.MaintenanceHours == nil {
		return MaintenanceStatusUnkown
	}

	if e.LastMaintenanceHourMeter == nil {
		*e.LastMaintenanceHourMeter = 0
	}

	hoursSinceMaintenance := e.CurrentHourMeter - *e.LastMaintenanceHourMeter

	if hoursSinceMaintenance >= *e.MaintenanceHours {
		return MaintenanceStatusRequired
	}

	return MaintenanceStatusOperational
}

func (e *Equipment) HoursSinceMaintenance() int32 {
	if e == nil || e.LastMaintenanceHourMeter == nil {
		return 0
	}
	return e.CurrentHourMeter - *e.LastMaintenanceHourMeter
}
