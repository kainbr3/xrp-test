package repositories

import (
	"context"
	l "crypto-braza-tokens-api/utils/logger"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func (r *Repository) FindWallets(ctx context.Context) ([]*Wallet, error) {
	filter := bson.M{"is_active": true}
	findOptions := options.Find().SetSort(bson.M{"name": 1})

	cursor, err := r.walletsCollection.Find(ctx, filter, findOptions)
	if err != nil {
		l.Logger.Error("repository: error finding wallets", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []*Wallet
	if err = cursor.All(ctx, &result); err != nil {
		l.Logger.Error("repository: error parsing wallets result", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (r *Repository) FindWalletsByBlockchainId(ctx context.Context, blockchainId string) ([]*Wallet, error) {
	filter := bson.M{"blockchain": blockchainId}
	findOptions := options.Find().SetSort(bson.M{"name": 1})

	cursor, err := r.walletsCollection.Find(ctx, filter, findOptions)
	if err != nil {
		l.Logger.Error("repository: error finding wallets", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []*Wallet
	if err = cursor.All(ctx, &result); err != nil {
		l.Logger.Error("repository: error parsing wallets result", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (r *Repository) FindWalletsByBlockchainAndDomain(ctx context.Context, blockchainId, domain string) ([]*Wallet, error) {
	filter := bson.M{"blockchain": blockchainId, "domain": domain}
	findOptions := options.Find().SetSort(bson.M{"name": 1})

	cursor, err := r.walletsCollection.Find(ctx, filter, findOptions)
	if err != nil {
		l.Logger.Error("repository: error finding wallets", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []*Wallet
	if err = cursor.All(ctx, &result); err != nil {
		l.Logger.Error(fmt.Sprintf("repository: error finding wallets with blockchain %s and domain %s", blockchainId, domain), zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (r *Repository) FindWalletById(ctx context.Context, walletId string) (*Wallet, error) {
	objectID, err := primitive.ObjectIDFromHex(walletId)
	if err != nil {
		l.Logger.Error("repository: error converting wallet Id to ObjectID", zap.Error(err))
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	var result *Wallet

	err = r.walletsCollection.FindOne(ctx, filter, nil).Decode(&result)
	if err != nil {
		l.Logger.Error(fmt.Sprintf("repository: error finding wallet for id %s", walletId), zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (r *Repository) FindWalletByAddressAndBlockchain(ctx context.Context, address string, blockchainId string) (*Wallet, error) {
	filter := bson.M{"address": address, "blockchain": blockchainId}

	var result *Wallet

	err := r.walletsCollection.FindOne(ctx, filter, nil).Decode(&result)
	if err != nil {
		l.Logger.Error(fmt.Sprintf("repository: error finding wallet for address %s on blockchain %s", address, blockchainId), zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (r *Repository) FindWalletByBlockchainWalletTypeAndDomain(ctx context.Context, blockchainId, walletType, domain string) (*Wallet, error) {
	filter := bson.M{"blockchain": blockchainId, "type": walletType, "domain": domain}

	var result *Wallet

	err := r.walletsCollection.FindOne(ctx, filter, nil).Decode(&result)
	if err != nil {
		l.Logger.Error(fmt.Sprintf("repository: error finding wallet for blockchain %s, wallet type %s and domain %s", blockchainId, walletType, domain), zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (r *Repository) SaveWallet(ctx context.Context, wallet *Wallet) (primitive.ObjectID, error) {
	// Ensure the operation has a valid ObjectID
	if wallet.ID.IsZero() {
		wallet.ID = primitive.NewObjectID()
	}

	result, err := r.walletsCollection.InsertOne(ctx, wallet)
	if err != nil {
		l.Logger.Error("repository: error saving wallet", zap.Error(err))
		return primitive.NilObjectID, err
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		l.Logger.Error("repository: error converting inserted ID to ObjectID")
		return primitive.NilObjectID, errors.New("error converting inserted ID to ObjectID")
	}

	return id, nil
}

func (r *Repository) EditWallet(ctx context.Context, wallet *Wallet) (primitive.ObjectID, error) {
	filter := bson.M{"_id": wallet.ID}
	update := bson.M{
		"$set": bson.M{
			"Blockchain": wallet.Blockchain,
			"Name":       wallet.Name,
			"Address":    wallet.Address,
			"Type":       wallet.Type,
			"Domain":     wallet.Domain,
			"IsActive":   wallet.IsActive,
			"UpdatedAt":  wallet.UpdatedAt,
		},
	}

	_, err := r.walletsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		l.Logger.Error("repository: error updating wallet", zap.Error(err))
		return primitive.NilObjectID, err
	}

	return wallet.ID, nil
}

func (r *Repository) DeleteWallet(ctx context.Context, walletId string) error {
	objectID, err := primitive.ObjectIDFromHex(walletId)
	if err != nil {
		l.Logger.Error("repository: error converting wallet Id to ObjectID", zap.Error(err))
		return err
	}

	filter := bson.M{"_id": objectID}
	_, err = r.walletsCollection.DeleteOne(ctx, filter)
	if err != nil {
		l.Logger.Error("repository: error deleting wallet", zap.Error(err))
		return err
	}

	return nil
}

func (r *Repository) WalletExistsSave(ctx context.Context, blockchain, name, address, walleType, domain string, isActive bool) bool {
	filter := bson.M{
		"Blockchain": blockchain,
		"Name":       name,
		"Address":    address,
		"Type":       walleType,
		"Domain":     domain,
		"IsActive":   isActive,
	}

	count, err := r.walletsCollection.CountDocuments(ctx, filter)
	if err != nil {
		l.Logger.Error("repository: error checking wallet existence", zap.Error(err))
		return false
	}

	return count > 0
}
