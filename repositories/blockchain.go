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

func (r *Repository) FindBlockchains(ctx context.Context) ([]*Blockchain, error) {
	filter := bson.D{}
	findOptions := options.Find().SetSort(bson.M{"name": 1})

	cursor, err := r.blockchainsCollection.Find(ctx, filter, findOptions)
	if err != nil {
		l.Logger.Error("repository: error finding blockchains", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []*Blockchain
	if err = cursor.All(ctx, &result); err != nil {
		l.Logger.Error("repository: error parsing blockchains result", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (r *Repository) FindBlockchainTokens(ctx context.Context, blockchainId string) ([]*Token, error) {
	filter := bson.M{"blockchain": blockchainId, "is_active": true}
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

func (r *Repository) FindBlockchainById(ctx context.Context, blockchainId string) (*Blockchain, error) {
	objectID, err := primitive.ObjectIDFromHex(blockchainId)
	if err != nil {
		l.Logger.Error("repository: error converting blockchain Id to ObjectID", zap.Error(err))
		return nil, err
	}
	filter := bson.M{"_id": objectID}

	var result *Blockchain
	err = r.blockchainsCollection.FindOne(ctx, filter, nil).Decode(&result)
	if err != nil {
		l.Logger.Error("repository: error finding blockchain", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (r *Repository) FindBlockchainByAbbr(ctx context.Context, abbr string) (*Blockchain, error) {
	filter := bson.M{"abbr": abbr}

	var result *Blockchain
	err := r.blockchainsCollection.FindOne(ctx, filter, nil).Decode(&result)
	if err != nil {
		l.Logger.Error("repository: error finding blockchain", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (r *Repository) SaveBlockchain(ctx context.Context, blockchain *Blockchain) (primitive.ObjectID, error) {
	// Ensure the operation has a valid ObjectID
	if blockchain.ID.IsZero() {
		blockchain.ID = primitive.NewObjectID()
	}

	result, err := r.blockchainsCollection.InsertOne(ctx, blockchain)
	if err != nil {
		l.Logger.Error("repository: error saving blockchain", zap.Error(err))
		return primitive.NilObjectID, err
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		l.Logger.Error("repository: error converting inserted ID to ObjectID")
		return primitive.NilObjectID, errors.New("failed to retrieve inserted ID")
	}

	return id, nil
}

func (r *Repository) DeleteBlockchain(ctx context.Context, blockchainId string) error {
	objectID, err := primitive.ObjectIDFromHex(blockchainId)
	if err != nil {
		l.Logger.Error("repository: error converting blockchain Id to ObjectID", zap.Error(err))
		return err
	}
	filter := bson.M{"_id": objectID}

	_, err = r.blockchainsCollection.DeleteOne(ctx, filter)
	if err != nil {
		l.Logger.Error("repository: error deleting blockchain", zap.Error(err))
		return err
	}

	return nil
}

func (r *Repository) EditBlockchain(ctx context.Context, blockchain *Blockchain) (primitive.ObjectID, error) {
	filter := bson.M{"_id": blockchain.ID}
	update := bson.M{"$set": bson.M{
		"name":       blockchain.Name,
		"abbr":       blockchain.Abbr,
		"main_token": blockchain.MainToken,
		"is_active":  blockchain.IsActive,
		"updated_at": blockchain.UpdatedAt,
	}}

	_, err := r.blockchainsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		l.Logger.Error("repository: error updating blockchain", zap.Error(err))
		return primitive.NilObjectID, err
	}

	return blockchain.ID, nil
}

func (r *Repository) BlockchainExists(ctx context.Context, id string) bool {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		l.Logger.Error("repository: error converting blockchain Id to ObjectID", zap.Error(err))
		return false
	}
	filter := bson.M{"_id": objectID}

	count, err := r.blockchainsCollection.CountDocuments(ctx, filter)
	if err != nil {
		l.Logger.Error("repository: error checking blockchain existence", zap.Error(err))
		return false
	}

	return count > 0
}

func (r *Repository) BlockchainExistsSave(ctx context.Context, name, abbr, mainToken string) bool {
	filter := bson.M{"name": name, "abbr": abbr, "main_token": mainToken}

	count, err := r.blockchainsCollection.CountDocuments(ctx, filter)
	if err != nil {
		l.Logger.Error("repository: error checking blockchain existence", zap.Error(err))
		return false
	}

	return count > 0
}
