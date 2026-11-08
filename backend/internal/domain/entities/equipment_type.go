package entities

import "github.com/gofrs/uuid/v5"

type EquipmentType struct {
	ID               uuid.UUID `db:"id"`
	Name             string    `db:"name"`
	MaintenanceHours int32     `db:"maintenance_hours"`
}
