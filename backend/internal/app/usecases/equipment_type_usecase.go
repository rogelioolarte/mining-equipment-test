package usecases

import (
	"context"
	"minning/internal/app/dto"
	"minning/internal/domain/ports"
	"minning/internal/domain/ports/repositories"
	"minning/internal/infra/logger"
)

type EquipmentTypeUseCase struct {
	logger                  *logger.Logger
	manager                 ports.TransactionManager
	EquipmentTypeRepository repositories.EquipmentTypeRepository
}

func NewEquipmentTypeUseCase(
	l *logger.Logger,
	m ports.TransactionManager,
	EquipmentTypeRepository repositories.EquipmentTypeRepository,
) *EquipmentTypeUseCase {
	return &EquipmentTypeUseCase{
		logger:                  l,
		manager:                 m,
		EquipmentTypeRepository: EquipmentTypeRepository,
	}
}

func (u *EquipmentTypeUseCase) GetAll(ctx context.Context) ([]*dto.EquipmentTypeResponse, error) {
	equiptmentList, err := u.EquipmentTypeRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]*dto.EquipmentTypeResponse, 0, len(equiptmentList))

	for _, i := range equiptmentList {
		response = append(response, dto.ToResponseEquipmentType(i))
	}

	return response, nil
}
