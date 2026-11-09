package repositories

import (
	"context"
	"minning/internal/domain/entities"

	"github.com/gofrs/uuid/v5"
)

type EquipmentTypeRepository interface {
	Create(ctx context.Context, entity *entities.EquipmentType) (*entities.EquipmentType, error)
	UpdateByID(ctx context.Context, entity *entities.EquipmentType) (*entities.EquipmentType, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.EquipmentType, error)
	GetAll(ctx context.Context) ([]*entities.EquipmentType, error)
}
