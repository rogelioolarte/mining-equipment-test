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

type MaintenancePostgreRepository struct {
	m ports.TransactionManager
}

func NewMaintenancePostgreRepository(m ports.TransactionManager) repositories.MaintenanceRepository {
	return &MaintenancePostgreRepository{m: m}
}

func (r *MaintenancePostgreRepository) Create(ctx context.Context, e *entities.Maintenance) (*entities.Maintenance, error) {
	query := `
        INSERT INTO app.maintenance (id, equipment_id, hm_maintenance, maintenance_date, description)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, equipment_id, hm_maintenance, maintenance_date, description`

	var i entities.Maintenance
	err := r.m.DB(ctx).QueryRow(ctx, query, e.ID, e.EquipmentID, e.HourMeterAtMaintenance,
		e.MaintenanceDate, e.Description).Scan(
		&i.ID, &i.EquipmentID, &i.HourMeterAtMaintenance,
		&i.MaintenanceDate, &i.Description)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, entities.ErrAlreadyExistsEntity
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entities.ErrInvalidEntity
		}
		return nil, fmt.Errorf("failed to create maintenance: %w", err)
	}

	return &i, nil
}

func (r *MaintenancePostgreRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Maintenance, error) {
	query := `
	SELECT
		id,
		equipment_id,
		hm_maintenance,
		maintenance_date,
		description
	FROM app.maintenance
	WHERE id = $1`

	var i entities.Maintenance
	err := r.m.DB(ctx).QueryRow(ctx, query, id).Scan(&i.ID, &i.EquipmentID,
		&i.HourMeterAtMaintenance, &i.MaintenanceDate, &i.Description)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, entities.ErrAlreadyExistsEntity
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entities.ErrNotFoundEntity
		}
		return nil, fmt.Errorf("failed to get maintenance: %w", err)
	}

	return &i, nil
}

func (r *MaintenancePostgreRepository) GetAllByEquipmentID(ctx context.Context, id uuid.UUID) ([]*entities.Maintenance, error) {
	query := `
	SELECT
		m.id,
		m.equipment_id,
		m.hm_maintenance,
		m.maintenance_date,
		m.description
	FROM app.maintenance m
	WHERE m.equipment_id = $1
	ORDER BY
		m.maintenance_date DESC,
		m.id DESC`

	is := make([]*entities.Maintenance, 0)
	rows, err := r.m.DB(ctx).Query(ctx, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entities.ErrNotFoundEntity
		}
		return nil, fmt.Errorf("failed to get all maintenance: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var i entities.Maintenance
		if err := rows.Scan(&i.ID, &i.EquipmentID, &i.HourMeterAtMaintenance,
			&i.MaintenanceDate, &i.Description); err != nil {
			return nil, fmt.Errorf("failed to scan maintenance: %w", err)
		}
		is = append(is, &i)
	}

	return is, nil
}

func (r *MaintenancePostgreRepository) GetAll(ctx context.Context) ([]*entities.Maintenance, error) {
	query := `
	SELECT
		m.id,
		m.equipment_id,
		e.code AS equipment_code,
		et.name AS equipment_type,
		m.hm_maintenance,
		m.maintenance_date,
		m.description
	FROM app.maintenance m
	INNER JOIN app.equipment e
		ON e.id = m.equipment_id
	INNER JOIN app.equipment_type et
		ON et.id = e.type_id
	ORDER BY m.maintenance_date DESC`

	is := make([]*entities.Maintenance, 0)
	rows, err := r.m.DB(ctx).Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all maintenance: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var i entities.Maintenance
		if err := rows.Scan(&i.ID, &i.EquipmentID, &i.EquipmentCode, &i.EquipmentType,
			&i.HourMeterAtMaintenance, &i.MaintenanceDate, &i.Description); err != nil {
			return nil, fmt.Errorf("failed to scan maintenance: %w", err)
		}
		is = append(is, &i)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate maintenance: %w", err)
	}

	return is, nil
}
