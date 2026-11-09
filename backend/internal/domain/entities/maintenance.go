package entities

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Maintenance struct {
	ID                     uuid.UUID `db:"id"`
	EquipmentID            uuid.UUID `db:"equipment_id"`
	EquipmentCode          *string   `db:"equipment_code"`
	EquipmentType          *string   `db:"equipment_type"`
	HourMeterAtMaintenance int32     `db:"hm_maintenance"`
	MaintenanceDate        time.Time `db:"maintenance_date"`
	Description            *string   `db:"description"`
}
