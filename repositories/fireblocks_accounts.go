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

func (r *Repository) FindFireblocksAccounts(ctx context.Context) ([]*FireblocksAccount, error) {
	filter := bson.M{"is_active": true}
	findOptions := options.Find().SetSort(bson.M{"name": 1})

	cursor, err := r.fireblocksAccountsCollection.Find(ctx, filter, findOptions)
	if err != nil {
		l.Logger.Error("repository: error finding fireblocks accounts", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []*FireblocksAccount
	if err = cursor.All(ctx, &result); err != nil {
		l.Logger.Error("repository: error parsing fireblocks accounts result", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (r *Repository) FindFireblocksAccountByWalletId(ctx context.Context, walletId string) (*FireblocksAccount, error) {
	filter := bson.M{"wallet_id": walletId}

	var result *FireblocksAccount
	err := r.fireblocksAccountsCollection.FindOne(ctx, filter, nil).Decode(&result)
	if err != nil {
		l.Logger.Error("repository: error finding fireblocks account", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (r *Repository) FindFireblocksAccountById(ctx context.Context, fireblocksAccountId string) (*FireblocksAccount, error) {
	var result *FireblocksAccount

	filter := bson.M{"_id": fireblocksAccountId}
	err := r.fireblocksAccountsCollection.FindOne(ctx, filter, nil).Decode(&result)
	if err != nil {
		l.Logger.Error("repository: error finding fireblocks account", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (r *Repository) FindFireblocksAccountByVaultId(ctx context.Context, vaultId string) (*FireblocksAccount, error) {
	var result *FireblocksAccount

	filter := bson.M{"vault_id": vaultId}
	err := r.fireblocksAccountsCollection.FindOne(ctx, filter, nil).Decode(&result)
	if err != nil {
		l.Logger.Error("repository: error finding fireblocks account", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (r *Repository) SaveFireblocksAccount(ctx context.Context, fireblocksAccount *FireblocksAccount) (primitive.ObjectID, error) {
	// Ensure the operation has a valid ObjectID
	if fireblocksAccount.ID.IsZero() {
		fireblocksAccount.ID = primitive.NewObjectID()
	}

	result, err := r.fireblocksAccountsCollection.InsertOne(ctx, fireblocksAccount)
	if err != nil {
		l.Logger.Error("repository: error saving fireblocks account", zap.Error(err))
		return primitive.NilObjectID, err
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		l.Logger.Error("repository: error converting inserted ID to ObjectID")
		return primitive.NilObjectID, errors.New("failed to retrieve inserted ID")
	}

	return id, nil
}

func (r *Repository) DeleteFireblocksAccount(ctx context.Context, fireblocksAccountId string) error {
	objectID, err := primitive.ObjectIDFromHex(fireblocksAccountId)
	if err != nil {
		l.Logger.Error("repository: error converting fireblocksAccount Id to ObjectID", zap.Error(err))
		return err
	}
	filter := bson.M{"_id": objectID}

	_, err = r.fireblocksAccountsCollection.DeleteOne(ctx, filter)
	if err != nil {
		l.Logger.Error("repository: error deleting fireblocksAccount", zap.Error(err))
		return err
	}

	return nil
}

func (r *Repository) EditFireblocksAccount(ctx context.Context, fireblocksAccount *FireblocksAccount) (primitive.ObjectID, error) {
	filter := bson.M{"_id": fireblocksAccount.ID}
	update := bson.M{
		"$set": bson.M{
			"vault_id":   fireblocksAccount.VaultID,
			"wallet_id":  fireblocksAccount.WalletID,
			"asset_id":   fireblocksAccount.AssetID,
			"name":       fireblocksAccount.Name,
			"alias":      fireblocksAccount.Alias,
			"domain":     fireblocksAccount.Domain,
			"public_key": fireblocksAccount.PublicKey,
			"flags":      fireblocksAccount.Flags,
			"is_active":  fireblocksAccount.IsActive,
			"UpdatedAt":  fireblocksAccount.UpdatedAt,
		},
	}

	_, err := r.fireblocksAccountsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		l.Logger.Error("repository: error updating fireblocks account", zap.Error(err))
		return primitive.NilObjectID, err
	}

	return fireblocksAccount.ID, nil
}

func (r *Repository) FireblocksAccountExists(ctx context.Context, vaultID, walletID, assetID, name, alias, publicKey, domain string, flags int, isActive bool) bool {
	filter := bson.M{
		"vault_id":   vaultID,
		"wallet_id":  walletID,
		"asset_id":   assetID,
		"name":       name,
		"alias":      alias,
		"domain":     domain,
		"public_key": publicKey,
		"flags":      flags,
		"is_active":  isActive,
	}

	count, err := r.operationsTypesCollection.CountDocuments(ctx, filter)
	if err != nil {
		l.Logger.Error("repository: error checking fireblocks account existence", zap.Error(err))
		return false
	}

	return count > 0
}
