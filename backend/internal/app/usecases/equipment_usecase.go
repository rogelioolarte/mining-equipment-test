package usecases

import (
	"context"
	"minning/internal/app/dto"
	"minning/internal/domain/entities"
	"minning/internal/domain/ports"
	"minning/internal/domain/ports/repositories"
	"minning/internal/infra/logger"

	"github.com/gofrs/uuid/v5"
)

type EquipmentUseCase struct {
	logger              *logger.Logger
	manager             ports.TransactionManager
	EquipmentRepository repositories.EquipmentRepository
}

func NewEquipmentUseCase(
	l *logger.Logger,
	m ports.TransactionManager,
	EquipmentRepository repositories.EquipmentRepository,
) *EquipmentUseCase {
	return &EquipmentUseCase{
		logger:              l,
		manager:             m,
		EquipmentRepository: EquipmentRepository,
	}
}

func (u *EquipmentUseCase) GetAll(ctx context.Context) ([]*dto.EquipmentResponse, error) {
	equiptmentList, err := u.EquipmentRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]*dto.EquipmentResponse, 0, len(equiptmentList))

	for _, i := range equiptmentList {
		response = append(response, dto.ToResponseEquipment(i))
	}

	return response, nil
}

func (u *EquipmentUseCase) GetByID(ctx context.Context, id uuid.UUID) (*dto.EquipmentResponse, error) {
	equiptment, err := u.EquipmentRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return dto.ToResponseEquipment(equiptment), nil
}

func (u *EquipmentUseCase) Create(ctx context.Context, request *dto.EquipmentRequest) (*dto.EquipmentResponse, error) {
	entity, err := request.ToEntity()
	if err != nil {
		return nil, err
	}

	var created *entities.Equipment
	err = u.manager.Transaction(ctx, func(txCtx context.Context) error {
		created, err = u.EquipmentRepository.Create(ctx, entity)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return dto.ToResponseEquipment(created), nil
}

func (u *EquipmentUseCase) UpdateByID(ctx context.Context, id uuid.UUID, request *dto.EquipmentRequest) (*dto.EquipmentResponse, error) {
	entity, err := request.ToEntity()
	if err != nil {
		return nil, err
	}
	entity.ID = id

	var updated *entities.Equipment
	err = u.manager.Transaction(ctx, func(txCtx context.Context) error {
		updated, err = u.EquipmentRepository.UpdateByID(ctx, entity)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return dto.ToResponseEquipment(updated), nil
}

func (u *EquipmentUseCase) DeleteByID(ctx context.Context, id uuid.UUID) error {
	err := u.manager.Transaction(ctx, func(txCtx context.Context) error {
		err := u.EquipmentRepository.DeleteByID(ctx, id)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
