package repositories

import (
	"context"
	"errors"
	"fmt"
	"minning/internal/domain/entities"
	"minning/internal/domain/ports"
	"minning/internal/domain/ports/repositories"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type EquipmentPostgreRepository struct {
	m ports.TransactionManager
}

func NewEquipmentPostgreRepository(m ports.TransactionManager) repositories.EquipmentRepository {
	return &EquipmentPostgreRepository{m: m}
}

func (r *EquipmentPostgreRepository) Create(ctx context.Context, e *entities.Equipment) (*entities.Equipment, error) {
	query := `
	INSERT INTO app.equipment (id, type_id, code, current_h_meter)
	VALUES ($1, $2, $3, $4)
	RETURNING id, type_id, code, current_h_meter`

	var i entities.Equipment
	err := r.m.DB(ctx).QueryRow(ctx, query, e.ID, e.TypeID, e.Code, e.CurrentHourMeter).Scan(
		&i.ID, &i.TypeID, &i.Code, &i.CurrentHourMeter)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, entities.ErrAlreadyExistsEntity
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entities.ErrInvalidEntity
		}
		return nil, fmt.Errorf("failed to create equipment: %w", err)
	}

	return &i, nil
}

func (r *EquipmentPostgreRepository) UpdateByID(ctx context.Context, e *entities.Equipment) (*entities.Equipment, error) {
	query := `
	UPDATE app.equipment
	SET
		type_id = $2,
		code = $3,
		current_h_meter = $4
	WHERE id = $1
	RETURNING id, type_id, code, current_h_meter`

	var i entities.Equipment
	err := r.m.DB(ctx).QueryRow(ctx, query, e.ID, e.TypeID, e.Code,
		e.CurrentHourMeter).Scan(
		&i.ID, &i.TypeID, &i.Code, &i.CurrentHourMeter)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, entities.ErrAlreadyExistsEntity
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entities.ErrNotFoundEntity
		}
		return nil, fmt.Errorf("failed to update equipment: %w", err)
	}

	return &i, nil
}

func (r *EquipmentPostgreRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM app.equipment WHERE id = $1`

	_, err := r.m.DB(ctx).Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete equipment: %w", err)
	}

	return nil
}

func (r *EquipmentPostgreRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Equipment, error) {
	query := `
	SELECT
		e.id,
		e.type_id,
		e.code,
		e.current_h_meter,
		e.last_mh_meter,
		et.name AS equipment_type,
		et.maintenance_hours
	FROM app.equipment e
	INNER JOIN app.equipment_type et
		ON et.id = e.type_id
	WHERE e.id = $1`

	var i entities.Equipment
	err := r.m.DB(ctx).QueryRow(ctx, query, id).Scan(
		&i.ID, &i.TypeID, &i.Code, &i.CurrentHourMeter,
		&i.LastMaintenanceHourMeter, &i.EquipmentType, &i.MaintenanceHours)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entities.ErrNotFoundEntity
		}
		return nil, fmt.Errorf("failed to get equipment: %w", err)
	}

	return &i, nil
}

func (r *EquipmentPostgreRepository) GetAll(ctx context.Context) ([]*entities.Equipment, error) {
	query := `
	SELECT
		e.id,
		e.type_id,
		et.name AS equipment_type,
		e.code,
		e.current_h_meter,
		et.maintenance_hours,
		COALESCE(m.hm_maintenance, 0) AS last_mh_meter
	FROM app.equipment e
	INNER JOIN app.equipment_type et
		ON et.id = e.type_id
	LEFT JOIN LATERAL (
		SELECT m.hm_maintenance
		FROM app.maintenance m
		WHERE m.equipment_id = e.id
		ORDER BY m.maintenance_date DESC, m.id DESC
		LIMIT 1
	) m ON TRUE
	ORDER BY e.code ASC`

	is := make([]*entities.Equipment, 0)
	rows, err := r.m.DB(ctx).Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all equipment: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var i entities.Equipment
		if err := rows.Scan(&i.ID, &i.TypeID, &i.EquipmentType, &i.Code,
			&i.CurrentHourMeter, &i.MaintenanceHours, &i.LastMaintenanceHourMeter); err != nil {
			return nil, fmt.Errorf("failed to scan equipment: %w", err)
		}
		is = append(is, &i)
	}

	return is, nil
}
