package usecases

import (
	"context"
	"minning/internal/app/dto"
	"minning/internal/domain/ports"
	"minning/internal/domain/ports/repositories"
	"minning/internal/infra/logger"
)

type MaintenanceUseCase struct {
	logger                *logger.Logger
	manager               ports.TransactionManager
	MaintenanceRepository repositories.MaintenanceRepository
}

func NewMaintenanceUseCase(
	l *logger.Logger,
	m ports.TransactionManager,
	MaintenanceRepository repositories.MaintenanceRepository,
) *MaintenanceUseCase {
	return &MaintenanceUseCase{
		logger:                l,
		manager:               m,
		MaintenanceRepository: MaintenanceRepository,
	}
}

func (u *MaintenanceUseCase) GetAll(ctx context.Context) ([]*dto.MaintenanceResponse, error) {
	equiptmentList, err := u.MaintenanceRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]*dto.MaintenanceResponse, 0, len(equiptmentList))

	for _, i := range equiptmentList {
		response = append(response, dto.ToResponseMaintenance(i))
	}

	return response, nil
}
