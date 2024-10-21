package repositories

import (
	"context"
	"fmt"

	l "crypto-braza-tokens-api/utils/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func (r *Repository) FindOperationsTypes(ctx context.Context) ([]*OperationType, error) {
	filter := bson.D{}
	findOptions := options.Find().SetSort(bson.M{"name": 1})

	cursor, err := r.operationsTypesCollection.Find(ctx, filter, findOptions)
	if err != nil {
		l.Logger.Error("repository: error finding operations types", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []*OperationType
	if err = cursor.All(ctx, &result); err != nil {
		l.Logger.Error("repository: error parsing operations types result", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (r *Repository) FindOperationTypeById(ctx context.Context, id string) (*OperationType, error) {
	var result *OperationType

	filter := bson.M{"_id": id}
	err := r.operationsTypesCollection.FindOne(ctx, filter, nil).Decode(&result)
	if err != nil {
		l.Logger.Error("repository: error operation type", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (r *Repository) SaveOperationType(ctx context.Context, operationType *OperationType) (primitive.ObjectID, error) {
	// Ensure the operation type has a valid ObjectID
	if operationType.ID.IsZero() {
		operationType.ID = primitive.NewObjectID()
	}

	result, err := r.operationsTypesCollection.InsertOne(ctx, operationType)
	if err != nil {
		l.Logger.Error("repository: error saving operation type", zap.Error(err))
		return primitive.NilObjectID, err
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		l.Logger.Error("repository: error converting inserted ID to ObjectID")
		return primitive.NilObjectID, fmt.Errorf("failed to retrieve inserted ID")
	}

	return id, nil
}

func (r *Repository) DeleteOperationType(ctx context.Context, operationTypeId string) error {
	objectID, err := primitive.ObjectIDFromHex(operationTypeId)
	if err != nil {
		l.Logger.Error("repository: error converting operation type Id to ObjectID", zap.Error(err))
		return err
	}
	filter := bson.M{"_id": objectID}

	_, err = r.operationsTypesCollection.DeleteOne(ctx, filter)
	if err != nil {
		l.Logger.Error("repository: error deleting operation type", zap.Error(err))
		return err
	}

	return nil
}

func (r *Repository) EditOperationType(ctx context.Context, operationType *OperationType) (primitive.ObjectID, error) {
	filter := bson.M{"_id": operationType.ID}
	update := bson.M{
		"$set": bson.M{
			"name":      operationType.Name,
			"is_active": operationType.IsActive,
			"UpdatedAt": operationType.UpdatedAt,
		},
	}

	_, err := r.operationsTypesCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		l.Logger.Error("repository: error updating operation type", zap.Error(err))
		return primitive.NilObjectID, err
	}

	return operationType.ID, nil
}

func (r *Repository) OpTypeExists(ctx context.Context, opType string) bool {
	filter := bson.M{"name": opType}
	count, err := r.operationsTypesCollection.CountDocuments(ctx, filter)
	if err != nil {
		l.Logger.Error("repository: error checking operation type existence", zap.Error(err))
		return false
	}

	return count > 0
}
