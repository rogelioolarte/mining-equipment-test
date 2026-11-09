package repositories

import (
	"context"
	"minning/internal/domain/entities"

	"github.com/gofrs/uuid/v5"
)

type EquipmentRepository interface {
	Create(ctx context.Context, entity *entities.Equipment) (*entities.Equipment, error)
	UpdateByID(ctx context.Context, entity *entities.Equipment) (*entities.Equipment, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Equipment, error)
	GetAll(ctx context.Context) ([]*entities.Equipment, error)
}
