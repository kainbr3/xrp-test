package repositories

import (
	"context"
	l "crypto-braza-tokens-api/utils/logger"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func (r *Repository) FindTransactionTypes(ctx context.Context) ([]*TransactionType, error) {
	filter := bson.D{}
	findOptions := options.Find().SetSort(bson.M{"name": 1})

	cursor, err := r.transactionsTypesCollection.Find(ctx, filter, findOptions)
	if err != nil {
		l.Logger.Error("repository: error finding transactions types", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []*TransactionType
	if err = cursor.All(ctx, &result); err != nil {
		l.Logger.Error("repository: error parsing transactions types result", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (r *Repository) FindTransactionTypeById(ctx context.Context, transactionTypeId string) (*TransactionType, error) {

	objectID, err := primitive.ObjectIDFromHex(transactionTypeId)
	if err != nil {
		l.Logger.Error("repository: error converting transaction type Id to ObjectID", zap.Error(err))
		return nil, err
	}
	filter := bson.M{"_id": objectID}

	var result *TransactionType
	err = r.transactionsTypesCollection.FindOne(ctx, filter, nil).Decode(&result)
	if err != nil {
		l.Logger.Error("repository: error transaction type", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (r *Repository) SaveTransactionType(ctx context.Context, transactionType *TransactionType) (primitive.ObjectID, error) {
	// Ensure the transaction type has a valid ObjectID
	if transactionType.ID.IsZero() {
		transactionType.ID = primitive.NewObjectID()
	}

	result, err := r.transactionsTypesCollection.InsertOne(ctx, transactionType)
	if err != nil {
		l.Logger.Error("repository: error saving transaction type", zap.Error(err))
		return primitive.NilObjectID, err
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		l.Logger.Error("repository: error converting inserted ID to ObjectID")
		return primitive.NilObjectID, fmt.Errorf("failed to retrieve inserted ID")
	}

	return id, nil
}

func (r *Repository) DeleteTransactionType(ctx context.Context, transactionTypeId string) error {
	objectID, err := primitive.ObjectIDFromHex(transactionTypeId)
	if err != nil {
		l.Logger.Error("repository: error converting transaction type Id to ObjectID", zap.Error(err))
		return err
	}
	filter := bson.M{"_id": objectID}

	_, err = r.transactionsTypesCollection.DeleteOne(ctx, filter)
	if err != nil {
		l.Logger.Error("repository: error deleting transaction type", zap.Error(err))
		return err
	}

	return nil
}

func (r *Repository) EditTransactionType(ctx context.Context, transactionType *TransactionType) (primitive.ObjectID, error) {
	filter := bson.M{"_id": transactionType.ID}
	update := bson.M{
		"$set": bson.M{
			"name":      transactionType.Name,
			"is_active": transactionType.IsActive,
			"UpdatedAt": transactionType.UpdatedAt,
		},
	}

	_, err := r.transactionsTypesCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		l.Logger.Error("repository: error updating transaction type", zap.Error(err))
		return primitive.NilObjectID, err
	}

	return transactionType.ID, nil
}

func (r *Repository) TxTypeExists(ctx context.Context, txType string) bool {
	filter := bson.M{"name": txType}
	count, err := r.transactionsTypesCollection.CountDocuments(ctx, filter)
	if err != nil {
		l.Logger.Error("repository: error checking transaction type existence", zap.Error(err))
		return false
	}

	return count > 0
}
