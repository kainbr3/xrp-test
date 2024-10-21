package blockchains

import (
	"context"
	r "crypto-braza-tokens-api/repositories"
	l "crypto-braza-tokens-api/utils/logger"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type BlockchainService struct {
	repo *r.Repository
}

func NewBlockchainService(repo *r.Repository) *BlockchainService {
	return &BlockchainService{repo}
}

func (bs *BlockchainService) FindAll(ctx context.Context) ([]*Blockchain, error) {
	result := []*Blockchain{}

	blockchains, err := bs.repo.FindBlockchains(ctx)
	if err != nil {
		l.Logger.Error("blockchain service: error finding blockchains", zap.Error(err))
		return nil, err
	}

	for _, blockchain := range blockchains {
		result = append(result, &Blockchain{
			ID:        blockchain.ID.Hex(),
			Name:      blockchain.Name,
			Abbr:      blockchain.Abbr,
			MainToken: blockchain.MainToken,
			IsActive:  blockchain.IsActive,
			CreatedAt: blockchain.CreatedAt,
			UpdatedAt: blockchain.UpdatedAt,
		})
	}

	return result, nil
}

func (bs *BlockchainService) FindBlockchainByID(ctx context.Context, blockchainID string) (*Blockchain, error) {
	blockchainResult, err := bs.repo.FindBlockchainById(ctx, blockchainID)
	if err != nil {
		l.Logger.Error("blockchain service: error finding blockchain", zap.Error(err))
		return nil, err
	}

	blockchain := &Blockchain{
		ID:        blockchainResult.ID.Hex(),
		IsActive:  blockchainResult.IsActive,
		CreatedAt: blockchainResult.CreatedAt,
		UpdatedAt: blockchainResult.UpdatedAt,
		Name:      blockchainResult.Name,
		Abbr:      blockchainResult.Abbr,
		MainToken: blockchainResult.MainToken,
	}

	return blockchain, nil
}

func (bs *BlockchainService) SaveBlockchain(ctx context.Context, name, abbr, mainToken string, isActive bool) (primitive.ObjectID, error) {
	blockchain := &r.Blockchain{
		Name: name, Abbr: abbr, MainToken: mainToken, IsActive: isActive, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}

	if isValid := bs.repo.BlockchainExistsSave(ctx, name, abbr, mainToken); isValid {
		l.Logger.Error("blockchain service: blockchain already exists", zap.String("name", name), zap.String("abbr", abbr), zap.String("mainToken", mainToken))
		return primitive.ObjectID{}, errors.New("blockchain already exists")
	}

	id, err := bs.repo.SaveBlockchain(ctx, blockchain)
	if err != nil {
		l.Logger.Error("blockchain service: error saving blockchain", zap.Error(err))
		return primitive.ObjectID{}, err
	}

	return id, nil
}

func (bs *BlockchainService) DeleteBlockchain(ctx context.Context, blockchainId string) (string, error) {
	err := bs.repo.DeleteBlockchain(ctx, blockchainId)
	if err != nil {
		l.Logger.Error("blockchain service: error deleting blockchain", zap.Error(err))
		return "", err
	}

	return blockchainId, nil
}

func (bs *BlockchainService) EditBlockchain(ctx context.Context, name, abbr, mainToken string, isActive bool) (primitive.ObjectID, error) {
	blockchain := &r.Blockchain{
		Name: name, Abbr: abbr, MainToken: mainToken, IsActive: isActive, UpdatedAt: time.Now(),
	}

	id, err := bs.repo.EditBlockchain(ctx, blockchain)
	if err != nil {
		l.Logger.Error("blockchain service: error updating blockchain", zap.Error(err))
		return primitive.ObjectID{}, err
	}

	return id, nil
}
