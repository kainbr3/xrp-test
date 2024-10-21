package wallet

import (
	"context"
	xsc "crypto-braza-tokens-api/clients/xrp-scan"
	r "crypto-braza-tokens-api/repositories"
	kvs "crypto-braza-tokens-api/utils/keys-values"
	l "crypto-braza-tokens-api/utils/logger"
	"errors"
	"strings"
	"time"

	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type WalletService struct {
	repo    *r.Repository
	xscCli  *xsc.XrpScanClient
	network string
}

func NewWalletService(repo *r.Repository) *WalletService {
	xrpScanCli, err := xsc.NewXrpScanClient()
	if err != nil {
		panic(err)
	}

	network, err := kvs.Get("XRP_NETWORK")
	if err != nil {
		panic(err)
	}

	return &WalletService{repo, xrpScanCli, network}
}

func (ws *WalletService) FindAll(ctx context.Context) ([]*Wallet, error) {
	result := []*Wallet{}

	blockchains, err := ws.repo.FindBlockchains(ctx)
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

	wallets, err := ws.repo.FindWallets(ctx)
	if err != nil {
		l.Logger.Error("service: error finding wallets", zap.Error(err))
		return nil, err
	}

	for _, wallet := range wallets {
		result = append(result, &Wallet{
			Base: Base{
				ID:        wallet.ID.Hex(),
				IsActive:  wallet.IsActive,
				CreatedAt: wallet.CreatedAt,
				UpdatedAt: wallet.UpdatedAt,
			},
			Name:         wallet.Name,
			Address:      wallet.Address,
			Type:         wallet.Type,
			BlockchainID: wallet.Blockchain,
			Blockchain:   blockchainsMap[wallet.Blockchain],
		})
	}

	return result, nil
}

func (ws *WalletService) FindWalletById(ctx context.Context, walletId string) (*Wallet, error) {
	walletResult, err := ws.repo.FindWalletById(ctx, walletId)
	if err != nil {
		l.Logger.Error("service: error finding wallet", zap.Error(err))
		return nil, err
	}

	result := &Wallet{
		Base: Base{
			ID:        walletResult.ID.Hex(),
			IsActive:  walletResult.IsActive,
			CreatedAt: walletResult.CreatedAt,
			UpdatedAt: walletResult.UpdatedAt,
		},
		Name:         walletResult.Name,
		Address:      walletResult.Address,
		Type:         walletResult.Type,
		BlockchainID: walletResult.Blockchain,
	}

	return result, nil
}

func (ws *WalletService) FindWalletByAddressAndBlockchain(ctx context.Context, address, blockchainId string) (*Wallet, error) {
	walletResult, err := ws.repo.FindWalletByAddressAndBlockchain(ctx, address, blockchainId)
	if err != nil {
		l.Logger.Error("service: error finding wallet", zap.Error(err))
		return nil, err
	}

	result := &Wallet{
		Base: Base{
			ID:        walletResult.ID.Hex(),
			IsActive:  walletResult.IsActive,
			CreatedAt: walletResult.CreatedAt,
			UpdatedAt: walletResult.UpdatedAt,
		},
		Name:         walletResult.Name,
		Address:      walletResult.Address,
		Type:         walletResult.Type,
		BlockchainID: walletResult.Blockchain,
	}

	return result, nil
}

func (ws *WalletService) FindWalletByBlockchainWalletTypeAndDomain(ctx context.Context, blockchainId, walletType, domain string) (*Wallet, error) {
	walletResult, err := ws.repo.FindWalletByBlockchainWalletTypeAndDomain(ctx, blockchainId, walletType, domain)
	if err != nil {
		l.Logger.Error("service: error finding wallet", zap.Error(err))
		return nil, err
	}

	result := &Wallet{
		Base: Base{
			ID:        walletResult.ID.Hex(),
			IsActive:  walletResult.IsActive,
			CreatedAt: walletResult.CreatedAt,
			UpdatedAt: walletResult.UpdatedAt,
		},
		Name:         walletResult.Name,
		Address:      walletResult.Address,
		Type:         walletResult.Type,
		BlockchainID: walletResult.Blockchain,
	}

	return result, nil
}

func (ws *WalletService) FindAllByBlockchainId(ctx context.Context, blockchainId string) ([]*Wallet, error) {
	result := []*Wallet{}

	blockchainResult, err := ws.repo.FindBlockchainById(ctx, blockchainId)
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

	wallets, err := ws.repo.FindWalletsByBlockchainId(ctx, blockchainId)
	if err != nil {
		l.Logger.Error("wallet service: error finding wallets", zap.Error(err))
		return nil, err
	}

	for _, wallet := range wallets {
		result = append(result, &Wallet{
			Base: Base{
				ID:        wallet.ID.Hex(),
				IsActive:  wallet.IsActive,
				CreatedAt: wallet.CreatedAt,
				UpdatedAt: wallet.UpdatedAt,
			},
			Name:         wallet.Name,
			Address:      wallet.Address,
			Type:         wallet.Type,
			BlockchainID: wallet.Blockchain,
			Blockchain:   blockchain,
		})
	}

	return result, nil
}

func (ws *WalletService) FindAllByBlockchainAndDomain(ctx context.Context, blockchainId, domain string) ([]*Wallet, error) {
	result := []*Wallet{}

	blockchainResult, err := ws.repo.FindBlockchainById(ctx, blockchainId)
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

	walletResult, err := ws.repo.FindWalletsByBlockchainAndDomain(ctx, blockchainId, domain)
	if err != nil {
		l.Logger.Error("service: error finding wallet", zap.Error(err))
		return nil, err
	}

	for _, wallet := range walletResult {
		result = append(result, &Wallet{
			Base: Base{
				ID:        wallet.ID.Hex(),
				IsActive:  wallet.IsActive,
				CreatedAt: wallet.CreatedAt,
				UpdatedAt: wallet.UpdatedAt,
			},
			Name:         wallet.Name,
			Address:      wallet.Address,
			Type:         wallet.Type,
			BlockchainID: wallet.Blockchain,
			Blockchain:   blockchain,
		})
	}

	return result, nil
}

func (ws *WalletService) GetMainBalance(ctx context.Context, address string) (*TokenBalance, error) {
	balance, err := ws.xscCli.GetAccountXrpBalance(ctx, address)
	if err != nil {
		errMsg := fmt.Errorf("failed to get xrp balance from address: %s with error: %v", address, err)
		l.Logger.Error("service: error get xrp balance", zap.Error(err))
		return nil, errMsg
	}

	return &TokenBalance{
		Address:  address,
		Contract: "XRP",
		Amount:   balance.XrpBalance,
	}, nil
}

func (ws *WalletService) GetTokensBalances(ctx context.Context, address string) ([]*TokenBalance, error) {
	balances, err := ws.xscCli.GetAccountTokensBalances(ctx, address)
	if err != nil {
		errMsg := fmt.Errorf("failed to get tokens balances from address: %s with error: %v", address, err)
		l.Logger.Error("service: error get tokens balances", zap.Error(err))
		return nil, errMsg
	}

	result := []*TokenBalance{}

	for _, balance := range balances {
		result = append(result, &TokenBalance{
			Address:  address,
			Contract: balance.Currency,
			Amount:   balance.Value,
		})
	}

	return result, nil
}

func (ws *WalletService) GetTokenObligations(ctx context.Context, issuerAddress string) (*TokenBalance, error) {
	obligations, err := ws.xscCli.GetTokenObligations(ctx, issuerAddress)
	if err != nil {
		errMsg := fmt.Errorf("failed to get token obligations from issuer address: %s with error: %v", issuerAddress, err)
		l.Logger.Error("service: error get token obligations", zap.Error(err))
		return nil, errMsg
	}

	result := &TokenBalance{
		Address:  issuerAddress,
		Contract: obligations.Currency,
		Amount:   obligations.Value,
	}

	return result, nil
}

func (ws *WalletService) GetAllBalances(ctx context.Context) ([]*TokenBalance, error) {
	if ws.network != "MAINNET" {
		errMsg := fmt.Errorf("wallet balances are only available on MAINNET network")
		l.Logger.Error("service: network != MAINNET", zap.Error(errMsg))
		return nil, errMsg
	}

	blockchain, err := ws.repo.FindBlockchainByAbbr(ctx, "XRP")
	if err != nil {
		l.Logger.Error("service: error finding blockchain", zap.Error(err))
		return nil, err
	}

	wallets, err := ws.FindAllByBlockchainId(ctx, blockchain.ID.Hex())
	if err != nil {
		l.Logger.Error("service: error finding blockchains", zap.Error(err))
		return nil, err
	}

	result := []*TokenBalance{}

	for _, wallet := range wallets {
		if strings.EqualFold(wallet.Type, "ISSUER") {
			balance, err := ws.GetTokenObligations(ctx, wallet.Address)
			if err != nil {
				l.Logger.Error("service: error getting token obligations", zap.Error(err))
				return nil, err
			}

			balance.Name = wallet.Name
			result = append(result, balance)
		} else {
			balanceMain, err := ws.GetMainBalance(ctx, wallet.Address)
			if err != nil {
				l.Logger.Error("service: error getting main balance", zap.Error(err))
				return nil, err
			}

			balanceMain.Name = wallet.Name
			result = append(result, balanceMain)

			balances, err := ws.GetTokensBalances(ctx, wallet.Address)
			if err != nil {
				l.Logger.Error("service: error getting tokens balances", zap.Error(err))
				return nil, err
			}

			for _, balance := range balances {
				for _, wallet := range wallets {
					if balance.Address == wallet.Address {
						balance.Name = wallet.Name
					}
				}
			}
			result = append(result, balances...)
		}
	}

	return result, nil
}

func (ws *WalletService) SaveWallet(ctx context.Context, blockchain, name, adress, walletType, domain string, isActive bool) (primitive.ObjectID, error) {
	if isValid := ws.repo.BlockchainExists(ctx, blockchain); !isValid {
		l.Logger.Error("service: blockchain already exists", zap.String("blockhain_id", blockchain))
		return primitive.ObjectID{}, errors.New("blockchain already exists")
	}

	wallet := &r.Wallet{
		Blockchain: blockchain, Name: name, Address: adress, Type: walletType, Domain: domain, IsActive: isActive, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}

	if isValid := ws.repo.WalletExistsSave(ctx, name, blockchain, adress, walletType, domain, isActive); isValid {
		l.Logger.Error("service: wallet already exists", zap.Bool("is_active", isValid))
		return primitive.ObjectID{}, errors.New("wallet already exists")
	}

	id, err := ws.repo.SaveWallet(ctx, wallet)
	if err != nil {
		l.Logger.Error("service: error saving wallet", zap.Error(err))
		return primitive.ObjectID{}, err
	}

	return id, nil
}

func (ws *WalletService) DeleteWallet(ctx context.Context, walletId string) (string, error) {
	err := ws.repo.DeleteWallet(ctx, walletId)
	if err != nil {
		l.Logger.Error("service: error deleting wallet", zap.Error(err))
		return "", err
	}

	return walletId, nil
}

func (ws *WalletService) EditWallet(ctx context.Context, blockchain, name, adress, walletType, domain string, isActive bool) (primitive.ObjectID, error) {
	wallet := &r.Wallet{
		Blockchain: blockchain, Name: name, Address: adress, Type: walletType, Domain: domain, IsActive: isActive, UpdatedAt: time.Now(),
	}

	id, err := ws.repo.EditWallet(ctx, wallet)
	if err != nil {
		l.Logger.Error("service: error updating wallet", zap.Error(err))
		return primitive.ObjectID{}, err
	}

	return id, nil
}
