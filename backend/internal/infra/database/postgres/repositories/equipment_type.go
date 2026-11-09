package repositories

import (
	"context"
	"fmt"
	"minning/internal/domain/entities"
	"minning/internal/domain/ports"
	"minning/internal/domain/ports/repositories"

	"github.com/gofrs/uuid/v5"
)

type EquipmentTypePostgreRepository struct {
	m ports.TransactionManager
}

func NewEquipmentTypePostgreRepository(m ports.TransactionManager) repositories.EquipmentTypeRepository {
	return &EquipmentTypePostgreRepository{m: m}
}

func (r *EquipmentTypePostgreRepository) Create(ctx context.Context, e *entities.EquipmentType) (*entities.EquipmentType, error) {
	query := `
    INSERT INTO app.equipment_type (id, name, maintenance_hours)
    VALUES ($1, $2, $3)
    RETURNING id, name, maintenance_hours`

	var i entities.EquipmentType
	err := r.m.DB(ctx).QueryRow(ctx, query, e.ID, e.Name, e.MaintenanceHours).Scan(
		&i.ID, &i.Name, &i.MaintenanceHours)
	if err != nil {
		return nil, fmt.Errorf("failed to create equipment type: %w", err)
	}

	return &i, nil
}

func (r *EquipmentTypePostgreRepository) UpdateByID(ctx context.Context, e *entities.EquipmentType) (*entities.EquipmentType, error) {
	query := `
	UPDATE app.equipment_type
	SET name = $2, maintenance_hours = $3
	WHERE id = $1
	RETURNING id, name, maintenance_hours`

	var i entities.EquipmentType
	err := r.m.DB(ctx).QueryRow(ctx, query, e.ID, e.Name, e.MaintenanceHours).Scan(
		&i.ID, &i.Name, &i.MaintenanceHours)
	if err != nil {
		return nil, fmt.Errorf("failed to update equipment type: %w", err)
	}

	return &i, nil
}

func (r *EquipmentTypePostgreRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM app.equipment_type WHERE id = $1`

	_, err := r.m.DB(ctx).Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete equipment type: %w", err)
	}

	return nil
}

func (r *EquipmentTypePostgreRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.EquipmentType, error) {
	query := `
	SELECT
		id,
		name,
		maintenance_hours
	FROM app.equipment_type
	WHERE id = $1`

	var i entities.EquipmentType
	err := r.m.DB(ctx).QueryRow(ctx, query, id).Scan(&i.ID, &i.Name, &i.MaintenanceHours)
	if err != nil {
		return nil, fmt.Errorf("failed to get equipment type: %w", err)
	}

	return &i, nil
}

func (r *EquipmentTypePostgreRepository) GetAll(ctx context.Context) ([]*entities.EquipmentType, error) {
	query := `
	SELECT
		id,
		name,
		maintenance_hours
	FROM app.equipment_type
	ORDER BY name ASC`

	is := make([]*entities.EquipmentType, 0)
	rows, err := r.m.DB(ctx).Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all equipment type: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var i entities.EquipmentType
		if err := rows.Scan(&i.ID, &i.Name, &i.MaintenanceHours); err != nil {
			return nil, fmt.Errorf("failed to scan equipment type: %w", err)
		}
		is = append(is, &i)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate equipment types: %w", err)
	}

	return is, nil
}
