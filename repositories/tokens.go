package repositories

import (
	"context"
	l "crypto-braza-tokens-api/utils/logger"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func (r *Repository) FindTokens(ctx context.Context) ([]*Token, error) {
	filter := bson.M{"is_active": true}
	findOptions := options.Find().SetSort(bson.D{{Key: "type", Value: 1}, {Key: "name", Value: 1}})

	cursor, err := r.tokensCollection.Find(ctx, filter, findOptions)
	if err != nil {
		l.Logger.Error("repository: error finding tokens", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []*Token
	if err = cursor.All(ctx, &result); err != nil {
		l.Logger.Error("repository: error parsing tokens result", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (r *Repository) FindTokensByBlockchainAndMintables(ctx context.Context, blockchainId string) ([]*Token, error) {
	filter := bson.M{"blockchain": blockchainId, "type": "ISSUED_CURRENCY", "is_active": true}
	findOptions := options.Find().SetSort(bson.D{{Key: "type", Value: 1}, {Key: "name", Value: 1}})

	cursor, err := r.tokensCollection.Find(ctx, filter, findOptions)
	if err != nil {
		l.Logger.Error("repository: error finding tokens", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []*Token
	if err = cursor.All(ctx, &result); err != nil {
		l.Logger.Error("repository: error parsing tokens result", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (r *Repository) FindTokenById(ctx context.Context, tokenId string) (*Token, error) {
	objectID, err := primitive.ObjectIDFromHex(tokenId)
	if err != nil {
		l.Logger.Error("repository: error converting token Id to ObjectID", zap.Error(err))
		return nil, err
	}
	filter := bson.M{"_id": objectID}

	var result *Token
	err = r.tokensCollection.FindOne(ctx, filter, nil).Decode(&result)
	if err != nil {
		l.Logger.Error("repository: error finding token", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (r *Repository) SaveToken(ctx context.Context, token *Token) (primitive.ObjectID, error) {
	// Ensure the operation has a valid ObjectID
	if token.ID.IsZero() {
		token.ID = primitive.NewObjectID()
	}

	result, err := r.tokensCollection.InsertOne(ctx, token)
	if err != nil {
		l.Logger.Error("repository: error saving token", zap.Error(err))
		return primitive.NilObjectID, err
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		l.Logger.Error("repository: error converting inserted ID to ObjectID")
		return primitive.NilObjectID, errors.New("error converting inserted ID to ObjectID")
	}

	return id, nil
}

func (r *Repository) DeleteToken(ctx context.Context, tokenId string) error {
	objectID, err := primitive.ObjectIDFromHex(tokenId)
	if err != nil {
		l.Logger.Error("repository: error converting token Id to ObjectID", zap.Error(err))
		return err
	}
	filter := bson.M{"_id": objectID}

	_, err = r.tokensCollection.DeleteOne(ctx, filter)
	if err != nil {
		l.Logger.Error("repository: error deleting token", zap.Error(err))
		return err
	}

	return nil
}

func (r *Repository) EditToken(ctx context.Context, token *Token) (primitive.ObjectID, error) {
	filter := bson.M{"_id": token.ID}
	update := bson.M{
		"$set": bson.M{
			"name":      token.Name,
			"abbr":      token.Abbr,
			"contract":  token.Contract,
			"type":      token.Type,
			"precision": token.Precision,
			"is_active": token.IsActive,
			"UpdatedAt": token.UpdatedAt,
		},
	}

	_, err := r.tokensCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		l.Logger.Error("repository: error updating token", zap.Error(err))
		return primitive.NilObjectID, err
	}

	return token.ID, nil
}

func (r *Repository) TokenExists(ctx context.Context, id string) bool {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		l.Logger.Error("repository: error converting token Id to ObjectID", zap.Error(err))
		return false
	}

	filter := bson.M{
		"_id":  objectID,
		"type": "ISSUED_CURRENCY",
	}

	count, err := r.tokensCollection.CountDocuments(ctx, filter)
	if err != nil {
		l.Logger.Error("repository: error checking token existence", zap.Error(err))
		return false
	}

	return count > 0
}

func (r *Repository) TokenExistsSave(ctx context.Context, name, abbr, contract, tokenType string, precision int, isActive bool) bool {
	filter := bson.M{
		"name":      name,
		"abbr":      abbr,
		"contract":  contract,
		"type":      tokenType,
		"precision": precision,
		"is_active": isActive,
	}

	count, err := r.tokensCollection.CountDocuments(ctx, filter)
	if err != nil {
		l.Logger.Error("repository: error checking token existence", zap.Error(err))
		return false
	}

	return count > 0
}
