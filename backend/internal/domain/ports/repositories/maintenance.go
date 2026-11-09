package repositories

import (
	"context"
	"minning/internal/domain/entities"

	"github.com/gofrs/uuid/v5"
)

type MaintenanceRepository interface {
	Create(ctx context.Context, entity *entities.Maintenance) (*entities.Maintenance, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Maintenance, error)
	GetAllByEquipmentID(ctx context.Context, id uuid.UUID) ([]*entities.Maintenance, error)
	GetAll(ctx context.Context) ([]*entities.Maintenance, error)
}
