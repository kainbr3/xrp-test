package repositories

import (
	"context"
	"errors"

	l "crypto-braza-tokens-api/utils/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func (r *Repository) FindOperationsDomains(ctx context.Context) ([]*OperationDomain, error) {
	filter := bson.D{}
	findOptions := options.Find().SetSort(bson.M{"name": 1})

	cursor, err := r.operationsDomainsCollection.Find(ctx, filter, findOptions)
	if err != nil {
		l.Logger.Error("repository: error finding operations domains", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []*OperationDomain
	if err = cursor.All(ctx, &result); err != nil {
		l.Logger.Error("repository: error parsing operations domains result", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (r *Repository) FindOperationDomainById(ctx context.Context, id string) (*OperationDomain, error) {
	var result *OperationDomain

	filter := bson.M{"_id": id}
	err := r.operationsDomainsCollection.FindOne(ctx, filter, nil).Decode(&result)
	if err != nil {
		l.Logger.Error("repository: error operation domain", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (r *Repository) SaveOperationDomain(ctx context.Context, operationDomain *OperationDomain) (primitive.ObjectID, error) {
	// Ensure the operation has a valid ObjectID
	if operationDomain.ID.IsZero() {
		operationDomain.ID = primitive.NewObjectID()
	}

	result, err := r.operationsDomainsCollection.InsertOne(ctx, operationDomain)
	if err != nil {
		l.Logger.Error("repository: error saving operation domain", zap.Error(err))
		return primitive.NilObjectID, err
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		l.Logger.Error("repository: error converting inserted ID to ObjectID")
		return primitive.NilObjectID, errors.New("failed to retrieve inserted ID")
	}

	return id, nil
}

func (r *Repository) DeleteOperationDomain(ctx context.Context, operationDomainId string) error {
	objectID, err := primitive.ObjectIDFromHex(operationDomainId)
	if err != nil {
		l.Logger.Error("repository: error converting operationDomain Id to ObjectID", zap.Error(err))
		return err
	}
	filter := bson.M{"_id": objectID}

	_, err = r.operationsDomainsCollection.DeleteOne(ctx, filter)
	if err != nil {
		l.Logger.Error("repository: error deleting operationDomain", zap.Error(err))
		return err
	}

	return nil
}

func (r *Repository) EditOperationDomain(ctx context.Context, operationDomain *OperationDomain) (primitive.ObjectID, error) {
	filter := bson.M{"_id": operationDomain.ID}
	update := bson.M{
		"$set": bson.M{
			"name":      operationDomain.Name,
			"is_active": operationDomain.IsActive,
			"UpdatedAt": operationDomain.UpdatedAt,
		},
	}

	_, err := r.operationsDomainsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		l.Logger.Error("repository: error updating operation domain", zap.Error(err))
		return primitive.NilObjectID, err
	}

	return operationDomain.ID, nil
}

func (r *Repository) DomainExists(ctx context.Context, domain string) bool {
	filter := bson.M{"name": domain}
	count, err := r.operationsDomainsCollection.CountDocuments(ctx, filter)
	if err != nil {
		l.Logger.Error("repository: error checking domain existence", zap.Error(err))
		return false
	}

	return count > 0
}
