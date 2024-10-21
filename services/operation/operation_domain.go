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

func (o *OperationService) GetOperationDomainsNames(ctx context.Context) (fiber.Map, error) {
	domains, err := o.repo.FindOperationsDomains(ctx)
	if err != nil {
		l.Logger.Error("operation service: error getting operation domains", zap.Error(err))
		return nil, err
	}

	if len(domains) == 0 {
		l.Logger.Error("operation service: no operation domains found")
		return nil, errors.New("no operation domains found")
	}

	list := []string{}

	for _, opdomain := range domains {
		list = append(list, opdomain.Name)
	}

	result := fiber.Map{"operation_domains": list}

	return result, nil
}

func (o *OperationService) FindAllDomains(ctx context.Context) ([]*OperationDomain, error) {
	domains, err := o.repo.FindOperationsDomains(ctx)
	if err != nil {
		l.Logger.Error("operation service: failed to find operation domains", zap.Error(err))
		return nil, err
	}

	if len(domains) == 0 {
		l.Logger.Error("operation service: no operation domains found")
		return nil, errors.New("no operation domains found")
	}

	result := []*OperationDomain{}

	for _, opDomain := range domains {
		operationDomain := &OperationDomain{
			Base: Base{
				ID:        opDomain.ID.Hex(),
				IsActive:  opDomain.IsActive,
				CreatedAt: opDomain.CreatedAt,
				UpdatedAt: opDomain.UpdatedAt,
			},
			Name: opDomain.Name,
		}
		result = append(result, operationDomain)
	}

	return result, nil
}

func (o *OperationService) FindOperationDomainById(ctx context.Context, id string) (*OperationDomain, error) {
	operationDomain, err := o.repo.FindOperationDomainById(ctx, id)
	if err != nil {
		l.Logger.Error("operation service: failed to find operation domain", zap.Error(err))
		return nil, err
	}

	if operationDomain == nil {
		l.Logger.Error("operation service: no operation domain found")
		return nil, errors.New("operation domain not found")
	}

	result := &OperationDomain{
		Base: Base{
			ID:        operationDomain.ID.Hex(),
			IsActive:  operationDomain.IsActive,
			CreatedAt: operationDomain.CreatedAt,
			UpdatedAt: operationDomain.UpdatedAt,
		},
		Name: operationDomain.Name,
	}

	return result, nil
}

func (o *OperationService) SaveOperationDomain(ctx context.Context, name string, isActive bool) (primitive.ObjectID, error) {
	operationDomain := &r.OperationDomain{
		Name: name, IsActive: isActive, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}

	if isValid := o.repo.DomainExists(ctx, name); isValid {
		l.Logger.Error("operation service: domain already exists", zap.String("name", name))
		return primitive.ObjectID{}, errors.New("domain already exists")
	}

	id, err := o.repo.SaveOperationDomain(ctx, operationDomain)
	if err != nil {
		l.Logger.Error("operation service: error saving operation domain", zap.Error(err))
		return primitive.ObjectID{}, err
	}

	return id, nil
}

func (o *OperationService) DeleteOperationDomain(ctx context.Context, operationDomainId string) (string, error) {
	err := o.repo.DeleteOperationDomain(ctx, operationDomainId)
	if err != nil {
		l.Logger.Error("operation service: error deleting operation domain", zap.Error(err))
		return "", err
	}

	return operationDomainId, nil
}

func (o *OperationService) EditOperationDomain(ctx context.Context, name string, isActive bool) (primitive.ObjectID, error) {
	operationDomain := &r.OperationDomain{
		Name: name, IsActive: isActive, UpdatedAt: time.Now(),
	}

	id, err := o.repo.EditOperationDomain(ctx, operationDomain)
	if err != nil {
		l.Logger.Error("operation service: error editing operation domain", zap.Error(err))
		return primitive.ObjectID{}, err
	}

	return id, nil
}
