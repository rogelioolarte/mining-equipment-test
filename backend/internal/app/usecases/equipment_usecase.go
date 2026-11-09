package usecases

import (
	"context"
	"minning/internal/app/dto"
	"minning/internal/domain/ports"
	"minning/internal/domain/ports/repositories"
	"minning/internal/infra/logger"
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
