package tokens

import (
	"context"
	r "crypto-braza-tokens-api/repositories"
	l "crypto-braza-tokens-api/utils/logger"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type TokenService struct {
	repo *r.Repository
}

func NewTokenService(repo *r.Repository) *TokenService {
	return &TokenService{repo}
}

func (ts *TokenService) FindAll(ctx context.Context) ([]*Token, error) {
	result := []*Token{}

	blockchains, err := ts.repo.FindBlockchains(ctx)
	if err != nil {
		l.Logger.Error("service: error finding blockchains", zap.Error(err))
		return nil, err
	}

	blockchainsMap := make(map[string]*Blockchain)
	for _, blockchain := range blockchains {
		blockchainsMap[blockchain.ID.Hex()] = &Blockchain{
			Base: Base{
				ID:        blockchain.ID.Hex(),
				IsActive:  blockchain.IsActive,
				CreatedAt: blockchain.CreatedAt,
				UpdatedAt: blockchain.UpdatedAt,
			},
			Name:      blockchain.Name,
			Abbr:      blockchain.Abbr,
			MainToken: blockchain.MainToken,
		}
	}

	tokens, err := ts.repo.FindTokensByBlockchainAndMintables(ctx, blockchains[0].ID.Hex())
	if err != nil {
		l.Logger.Error("service: error finding tokens", zap.Error(err))
		return nil, err
	}

	for _, token := range tokens {
		result = append(result, &Token{
			Base: Base{
				ID:        token.ID.Hex(),
				IsActive:  token.IsActive,
				CreatedAt: token.CreatedAt,
				UpdatedAt: token.UpdatedAt,
			},
			Name:         token.Name,
			Abbr:         token.Abbr,
			Contract:     token.Contract,
			Precision:    token.Precision,
			Type:         token.Type,
			BlockchainID: token.Blockchain,
			Blockchain:   blockchainsMap[token.Blockchain],
		})
	}

	return result, nil
}

func (ts *TokenService) FindTokenByID(ctx context.Context, tokenID string) (*Token, error) {
	tokenResult, err := ts.repo.FindTokenById(ctx, tokenID)
	if err != nil {
		l.Logger.Error("service: error finding token", zap.Error(err))
		return nil, err
	}

	return &Token{
		Base: Base{
			ID:        tokenResult.ID.Hex(),
			IsActive:  tokenResult.IsActive,
			CreatedAt: tokenResult.CreatedAt,
			UpdatedAt: tokenResult.UpdatedAt,
		},
		Name:         tokenResult.Name,
		Abbr:         tokenResult.Abbr,
		Contract:     tokenResult.Contract,
		Precision:    tokenResult.Precision,
		Type:         tokenResult.Type,
		BlockchainID: tokenResult.Blockchain,
	}, nil
}

func (ts *TokenService) FindTokensByBlockchainID(ctx context.Context, blockchainID string) ([]*Token, error) {
	result := []*Token{}

	blockchainResult, err := ts.repo.FindBlockchainById(ctx, blockchainID)
	if err != nil {
		l.Logger.Error("service: error finding blockchain", zap.Error(err))
		return nil, err
	}

	blockchain := &Blockchain{
		Base: Base{
			ID:        blockchainResult.ID.Hex(),
			IsActive:  blockchainResult.IsActive,
			CreatedAt: blockchainResult.CreatedAt,
			UpdatedAt: blockchainResult.UpdatedAt,
		},
		Name:      blockchainResult.Name,
		Abbr:      blockchainResult.Abbr,
		MainToken: blockchainResult.MainToken,
	}

	tokens, err := ts.repo.FindBlockchainTokens(ctx, blockchainID)
	if err != nil {
		l.Logger.Error("service: error finding tokens", zap.Error(err))
		return nil, err
	}

	for _, token := range tokens {
		result = append(result, &Token{
			Base: Base{
				ID:        token.ID.Hex(),
				IsActive:  token.IsActive,
				CreatedAt: token.CreatedAt,
				UpdatedAt: token.UpdatedAt,
			},
			Name:         token.Name,
			Abbr:         token.Abbr,
			Contract:     token.Contract,
			Precision:    token.Precision,
			Type:         token.Type,
			BlockchainID: token.Blockchain,
			Blockchain:   blockchain,
		})
	}

	return result, nil
}

func (ts *TokenService) SaveToken(ctx context.Context, blockchain, name, abbr, contract, tokenType string, precision int, isActive bool) (primitive.ObjectID, error) {
	if isValid := ts.repo.BlockchainExists(ctx, blockchain); !isValid {
		l.Logger.Error("service: blockchain already exists", zap.String("blockhain_id", blockchain))
		return primitive.ObjectID{}, errors.New("blockchain already exists")
	}

	token := &r.Token{
		Blockchain: blockchain, Name: name, Abbr: abbr, Contract: contract, Type: tokenType, Precision: precision, IsActive: isActive, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}

	if isValid := ts.repo.TokenExistsSave(ctx, name, abbr, contract, tokenType, precision, isActive); isValid {
		l.Logger.Error("service: token already exists", zap.Bool("is_active", isValid))
		return primitive.ObjectID{}, errors.New("token already exists")
	}

	id, err := ts.repo.SaveToken(ctx, token)
	if err != nil {
		l.Logger.Error("service: error saving token", zap.Error(err))
		return primitive.ObjectID{}, err
	}

	return id, nil
}

func (ts *TokenService) DeleteToken(ctx context.Context, tokenId string) (string, error) {
	err := ts.repo.DeleteToken(ctx, tokenId)
	if err != nil {
		l.Logger.Error("service: error deleting token", zap.Error(err))
		return "", err
	}

	return tokenId, nil
}

func (ts *TokenService) EditToken(ctx context.Context, blockchain, name, abbr, contract, tokenType string, precision int, isActive bool) (primitive.ObjectID, error) {
	token := &r.Token{
		Blockchain: blockchain, Name: name, Abbr: abbr, Contract: contract, Type: tokenType, Precision: precision, IsActive: isActive, UpdatedAt: time.Now(),
	}

	id, err := ts.repo.EditToken(ctx, token)
	if err != nil {
		l.Logger.Error("service: error updating token", zap.Error(err))
		return primitive.ObjectID{}, err
	}

	return id, nil
}
