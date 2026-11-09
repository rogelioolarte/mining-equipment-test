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
	maintenanceList, err := u.MaintenanceRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]*dto.MaintenanceResponse, 0, len(maintenanceList))

	for _, i := range maintenanceList {
		response = append(response, dto.ToResponseMaintenance(i))
	}

	return response, nil
}

func (u *MaintenanceUseCase) GetAllByEquipmentID(ctx context.Context, id uuid.UUID) ([]*dto.MaintenanceResponse, error) {
	maintenanceList, err := u.MaintenanceRepository.GetAllByEquipmentID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := make([]*dto.MaintenanceResponse, 0, len(maintenanceList))

	for _, i := range maintenanceList {
		response = append(response, dto.ToResponseMaintenance(i))
	}

	return response, nil
}

func (u *MaintenanceUseCase) Create(ctx context.Context, request *dto.MaintenanceRequest) (*dto.MaintenanceResponse, error) {
	entity, err := request.ToEntity()
	if err != nil {
		return nil, err
	}

	var created *entities.Maintenance
	err = u.manager.Transaction(ctx, func(txCtx context.Context) error {
		created, err = u.MaintenanceRepository.Create(ctx, entity)
		if err != nil {
			return err
		}
		return nil
	})

	return dto.ToResponseMaintenance(created), nil
}
