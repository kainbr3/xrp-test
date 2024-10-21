package repositories

import (
	"context"
	l "crypto-braza-tokens-api/utils/logger"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func (r *Repository) SaveOperation(ctx context.Context, operation *Operation) (primitive.ObjectID, error) {
	// Ensure the operation has a valid ObjectID
	if operation.ID.IsZero() {
		operation.ID = primitive.NewObjectID()
	}

	result, err := r.operationsCollection.InsertOne(ctx, operation)
	if err != nil {
		l.Logger.Error("error saving operation", zap.Error(err))
		return primitive.NilObjectID, err
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		l.Logger.Error("error converting inserted ID to ObjectID")
		return primitive.NilObjectID, errors.New("failed to retrieve inserted ID")
	}

	return id, nil
}

func (r *Repository) SaveOperationLog(ctx context.Context, log *OperationLog) error {
	// Ensure the operation log has a valid ObjectID
	if log.ID.IsZero() {
		log.ID = primitive.NewObjectID()
	}

	_, err := r.operationsLogsCollection.InsertOne(ctx, log)
	if err != nil {
		l.Logger.Error("error saving operation log", zap.Error(err))
		return err
	}

	return nil
}

func (r *Repository) UpdateOperationFireblocksIdAndStatus(ctx context.Context, operationId, fireblocksId, status string) error {
	objectID, err := primitive.ObjectIDFromHex(operationId)
	if err != nil {
		l.Logger.Error("error converting operation Id to ObjectID", zap.Error(err))
		return err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"fireblocks_id":     fireblocksId,
			"fireblocks_status": status,
			"updated_at":        time.Now(),
		},
	}

	_, err = r.operationsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		l.Logger.Error("error updating operation fireblocks id and status", zap.Error(err))
		return err
	}

	return nil
}

func (r *Repository) UpdateOperationFireblocksStatus(ctx context.Context, operationId, status string) error {
	objectID, err := primitive.ObjectIDFromHex(operationId)
	if err != nil {
		l.Logger.Error("error converting operation Id to ObjectID", zap.Error(err))
		return err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"fireblocks_status": status,
			"updated_at":        time.Now(),
		},
	}

	_, err = r.operationsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		l.Logger.Error("error updating operation fireblocks status", zap.Error(err))
		return err
	}

	return nil
}

func (r *Repository) UpdateOperationBlockchainStatus(ctx context.Context, operationId, newStatus, txHash, txLink string) error {
	objectID, err := primitive.ObjectIDFromHex(operationId)
	if err != nil {
		l.Logger.Error("error converting operation Id to ObjectID", zap.Error(err))
		return err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"blockchain_status": newStatus,
			"transaction_hash":  txHash,
			"transaction_link":  txLink,
			"updated_at":        time.Now(),
		},
	}

	_, err = r.operationsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		l.Logger.Error("error updating operation blockchain status", zap.Error(err))
		return err
	}

	return nil
}

func (r *Repository) FindOperationById(ctx context.Context, operationId string) (*Operation, error) {
	objectID, err := primitive.ObjectIDFromHex(operationId)
	if err != nil {
		l.Logger.Error("error converting operation Id to ObjectID", zap.Error(err))
		return nil, err
	}
	filter := bson.M{"_id": objectID}

	var result *Operation
	err = r.operationsCollection.FindOne(ctx, filter, nil).Decode(&result)
	if err != nil {
		l.Logger.Error("error finding operation", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (r *Repository) FindOperationLogsByOperationId(ctx context.Context, operationId string) ([]*OperationLog, error) {
	filter := bson.M{"operation_id": operationId}
	findOptions := options.Find().SetSort(bson.M{"created_at": -1})

	cursor, err := r.operationsLogsCollection.Find(ctx, filter, findOptions)
	if err != nil {
		l.Logger.Error("error finding operation logs", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []*OperationLog
	if err = cursor.All(ctx, &result); err != nil {
		l.Logger.Error("error parsing operation logs result", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (r *Repository) FindOperations(ctx context.Context) ([]*Operation, error) {
	filter := bson.D{}
	findOptions := options.Find().SetSort(bson.M{"updated_at": -1})

	cursor, err := r.operationsCollection.Find(ctx, filter, findOptions)
	if err != nil {
		l.Logger.Error("error finding operations", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []*Operation
	if err = cursor.All(ctx, &result); err != nil {
		l.Logger.Error("error parsing operations result", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (r *Repository) FindFilteredAndPaginatedOperations(ctx context.Context, params *QueryParams) (*PaginatedOperations, error) {
	// Default filter and sort
	filter := bson.M{}
	sort := bson.D{{Key: "updated_at", Value: -1}}
	page := 1
	limit := 10

	// Override defaults with provided params
	if params.FilterParam != "" && params.FilterValue != "" {
		filter[params.FilterParam] = params.FilterValue
	}
	if params.SortField != "" {
		order := 1
		if strings.EqualFold(params.SortOrder, "desc") {
			order = -1
		}
		sort = bson.D{{Key: params.SortField, Value: order}}
	}
	if params.Page > 0 {
		page = params.Page
	}
	if params.Limit > 0 {
		limit = params.Limit
	}

	// Limit the maximum page size to 100
	if limit > 100 {
		limit = 100
	}

	// Set pagination options
	findOptions := options.Find()
	if page > 0 && limit > 0 {
		findOptions.SetSkip(int64((page - 1) * limit))
		findOptions.SetLimit(int64(limit))
	}

	// Set sorting options
	findOptions.SetSort(sort)

	// Perform the query with filter, pagination, and sorting
	cursor, err := r.operationsCollection.Find(ctx, filter, findOptions)
	if err != nil {
		l.Logger.Error("error finding operations", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []*Operation
	for cursor.Next(ctx) {
		var operation Operation
		if err := cursor.Decode(&operation); err != nil {
			l.Logger.Error("error parsing operation result", zap.Error(err))
			return nil, err
		}
		result = append(result, &operation)
	}

	if err := cursor.Err(); err != nil {
		l.Logger.Error("cursor error", zap.Error(err))
		return nil, err
	}

	// Get the total count of documents matching the filter
	totalCount, err := r.operationsCollection.CountDocuments(ctx, filter)
	if err != nil {
		l.Logger.Error("error counting documents", zap.Error(err))
		return nil, err
	}

	// Calculate total pages
	totalPages := int((totalCount + int64(limit) - 1) / int64(limit))

	// Prepare the paginated result
	paginatedResult := &PaginatedOperations{
		TotalCount:   int(totalCount),
		TotalPages:   int(totalPages),
		CurrentPage:  page,
		NextPage:     page + 1,
		PreviousPage: page - 1,
		Data:         result,
	}

	// Adjust next and previous page values
	if paginatedResult.NextPage > paginatedResult.TotalPages {
		paginatedResult.NextPage = paginatedResult.TotalPages
	}
	if paginatedResult.PreviousPage < 1 {
		paginatedResult.PreviousPage = 1
	}

	return paginatedResult, nil
}
