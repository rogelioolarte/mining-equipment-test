package entities

import "github.com/gofrs/uuid/v5"

type Equipment struct {
	ID                       uuid.UUID `db:"id"`
	Code                     string    `db:"code"`
	TypeID                   uuid.UUID `db:"type_id"`
	CurrentHourMeter         int32     `db:"current_h_meter"`
	LastMaintenanceHourMeter int32     `db:"last_mh_meter"`
}
