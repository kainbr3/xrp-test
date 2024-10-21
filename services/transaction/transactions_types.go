package transaction

import (
	"context"
	r "crypto-braza-tokens-api/repositories"
	l "crypto-braza-tokens-api/utils/logger"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func (t *TransactionService) FindAllTypesMap(ctx context.Context) (fiber.Map, error) {
	types, err := t.repo.FindTransactionTypes(ctx)
	if err != nil {
		l.Logger.Error("transaction service: failed to find transaction types", zap.Error(err))
		return nil, err
	}

	if len(types) == 0 {
		l.Logger.Error("transaction service: no operation types found")
		return nil, errors.New("no transaction types found")
	}

	list := []string{}

	for _, opType := range types {
		list = append(list, opType.Name)
	}

	result := fiber.Map{"transaction_types": list}

	return result, nil
}

func (t *TransactionService) FindAllTypes(ctx context.Context) ([]*TransactionType, error) {
	types, err := t.repo.FindTransactionTypes(ctx)
	if err != nil {
		l.Logger.Error("transaction service: failed to find transaction types", zap.Error(err))
		return nil, err
	}

	if len(types) == 0 {
		l.Logger.Error("transaction service: no transaction types found")
		return nil, errors.New("no transaction types found")
	}

	result := []*TransactionType{}

	for _, txType := range types {
		transactionType := &TransactionType{
			ID:        txType.ID.Hex(),
			Name:      txType.Name,
			IsActive:  txType.IsActive,
			CreatedAt: txType.CreatedAt,
			UpdatedAt: txType.UpdatedAt,
		}
		result = append(result, transactionType)
	}

	return result, nil
}

func (t *TransactionService) FindTransactionTypeById(ctx context.Context, id string) (*TransactionType, error) {
	transactionType, err := t.repo.FindTransactionTypeById(ctx, id)
	if err != nil {
		l.Logger.Error("transaction service: failed to find transaction type", zap.Error(err))
		return nil, err
	}

	if transactionType == nil {
		l.Logger.Error("transaction service: no transaction type found")
		return nil, errors.New("transaction type not found")
	}

	result := &TransactionType{
		ID:        transactionType.ID.Hex(),
		Name:      transactionType.Name,
		IsActive:  transactionType.IsActive,
		CreatedAt: transactionType.CreatedAt,
		UpdatedAt: transactionType.UpdatedAt,
	}

	return result, nil
}

func (t *TransactionService) SaveTransactionType(ctx context.Context, transactionType string) (primitive.ObjectID, error) {
	tansactionTypeParsed := &r.TransactionType{
		Name: transactionType, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}

	if isValid := t.repo.TxTypeExists(ctx, transactionType); isValid {
		l.Logger.Error("transaction service: transaction type already exists", zap.String("name", transactionType))
		return primitive.ObjectID{}, errors.New("transaction type already exists")
	}

	id, err := t.repo.SaveTransactionType(ctx, tansactionTypeParsed)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	return id, nil
}

func (t *TransactionService) DeleteTransactionType(ctx context.Context, transactionTypeId string) (string, error) {
	err := t.repo.DeleteTransactionType(ctx, transactionTypeId)
	if err != nil {
		l.Logger.Error("transaction service: failed to delete transaction type", zap.Error(err))
		return "", err
	}

	return transactionTypeId, nil
}

func (t *TransactionService) EditTransactionType(ctx context.Context, txTypeId, txType string) (primitive.ObjectID, error) {
	objectID, err := primitive.ObjectIDFromHex(txTypeId)
	if err != nil {
		l.Logger.Error("repository: error converting transaction type Id to ObjectID", zap.Error(err))
		return primitive.ObjectID{}, err
	}

	transactionTypeParsed := &r.TransactionType{
		ID:   objectID,
		Name: txType, UpdatedAt: time.Now(),
	}

	id, err := t.repo.EditTransactionType(ctx, transactionTypeParsed)
	if err != nil {
		l.Logger.Error("transaction service: failed to edit transaction type", zap.Error(err))
		return primitive.ObjectID{}, err
	}

	return id, nil
}
