package operation

import (
	"context"
	"errors"
	"time"

	r "crypto-braza-tokens-api/repositories"
	l "crypto-braza-tokens-api/utils/logger"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func (o *OperationService) FindAllTypesMap(ctx context.Context) (fiber.Map, error) {
	types, err := o.repo.FindOperationsTypes(ctx)
	if err != nil {
		l.Logger.Error("operation service: failed to find operation types", zap.Error(err))
		return nil, err
	}

	if len(types) == 0 {
		l.Logger.Error("operation service: no operation types found")
		return nil, errors.New("no operation types found")
	}

	list := []string{}

	for _, opType := range types {
		list = append(list, opType.Name)
	}

	result := fiber.Map{"operation_types": list}

	return result, nil
}

func (o *OperationService) FindAllTypes(ctx context.Context) ([]*OperationType, error) {
	types, err := o.repo.FindOperationsTypes(ctx)
	if err != nil {
		l.Logger.Error("operation service: failed to find operation types", zap.Error(err))
		return nil, err
	}

	if len(types) == 0 {
		l.Logger.Error("operation service: no operation types found")
		return nil, errors.New("no operation types found")
	}

	result := []*OperationType{}

	for _, opType := range types {
		operationType := &OperationType{
			Base: Base{
				ID:        opType.ID.Hex(),
				IsActive:  opType.IsActive,
				CreatedAt: opType.CreatedAt,
				UpdatedAt: opType.UpdatedAt,
			},
			Name: opType.Name,
		}
		result = append(result, operationType)
	}

	return result, nil
}

func (o *OperationService) FindOperationTypeById(ctx context.Context, id string) (*OperationType, error) {
	operationType, err := o.repo.FindOperationTypeById(ctx, id)
	if err != nil {
		l.Logger.Error("operation service: failed to find operation type", zap.Error(err))
		return nil, err
	}

	if operationType == nil {
		l.Logger.Error("operation service: no operation type found")
		return nil, errors.New("operation type not found")
	}

	result := &OperationType{
		Base: Base{
			ID:        operationType.ID.Hex(),
			IsActive:  operationType.IsActive,
			CreatedAt: operationType.CreatedAt,
			UpdatedAt: operationType.UpdatedAt,
		},
		Name: operationType.Name,
	}

	return result, nil
}

func (o *OperationService) SaveOperationType(ctx context.Context, operationType string) (primitive.ObjectID, error) {
	operationTypeParsed := &r.OperationType{
		Name: operationType, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}

	if isValid := o.repo.OpTypeExists(ctx, operationType); isValid {
		l.Logger.Error("operation service: operation type already exists", zap.String("name", operationType))
		return primitive.ObjectID{}, errors.New("operation type already exists")
	}

	id, err := o.repo.SaveOperationType(ctx, operationTypeParsed)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	return id, nil
}

func (o *OperationService) DeleteOperationType(ctx context.Context, operationTypeId string) (string, error) {
	err := o.repo.DeleteOperationType(ctx, operationTypeId)
	if err != nil {
		l.Logger.Error("operation service: failed to delete operation type", zap.Error(err))
		return "", err
	}

	return operationTypeId, nil
}

func (o *OperationService) EditOperationType(ctx context.Context, operationType string) (primitive.ObjectID, error) {
	operationTypeParsed := &r.OperationType{
		Name: operationType, UpdatedAt: time.Now(),
	}

	id, err := o.repo.EditOperationType(ctx, operationTypeParsed)
	if err != nil {
		l.Logger.Error("operation service: failed to edit operation type", zap.Error(err))
		return primitive.ObjectID{}, err
	}

	return id, nil
}
