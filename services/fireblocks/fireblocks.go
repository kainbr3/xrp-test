package fireblocks

import (
	"context"
	fb "crypto-braza-tokens-api/clients/fireblocks"
	r "crypto-braza-tokens-api/repositories"
	l "crypto-braza-tokens-api/utils/logger"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type FireblocksService struct {
	repo  *r.Repository
	fbCli *fb.FireblocksClient
}

func NewFireblocksService(repo *r.Repository) *FireblocksService {
	fbCli, err := fb.NewFireblocksClient()
	if err != nil {
		l.Logger.Fatal("fireblocks service: failed to create a new fireblocks client", zap.Error(err))
	}

	return &FireblocksService{repo, fbCli}
}

func (fb *FireblocksService) FindAllFireblocksAccounts(ctx context.Context) ([]*FireblocksAccount, error) {
	result := []*FireblocksAccount{}

	blockchains, err := fb.repo.FindBlockchains(ctx)
	if err != nil {
		l.Logger.Error("fireblocks service: error finding blockchains", zap.Error(err))
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

	wallets, err := fb.repo.FindWallets(ctx)
	if err != nil {
		l.Logger.Error("fireblocks service: error finding wallets", zap.Error(err))
		return nil, err
	}

	walletsMap := make(map[string]*Wallet)
	for _, wallet := range wallets {
		walletsMap[wallet.ID.Hex()] = &Wallet{
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
		}
	}

	fireblocksAccounts, err := fb.repo.FindFireblocksAccounts(ctx)
	if err != nil {
		l.Logger.Error("fireblocks service: error finding fireblocks accounts", zap.Error(err))
		return nil, err
	}

	for _, account := range fireblocksAccounts {
		// Get the public key info from the Fireblocks API
		fbPubKey := ""

		fbAccResult, err := fb.fbCli.GetPublicKeyInfoFromVaultAccount(ctx, account.VaultID, account.AssetID, 0, 0)
		if err != nil {
			l.Logger.Error("fireblocks service: error getting public key info from Fireblocks API", zap.Error(err))
		}

		fbPubKey = fbAccResult.PublicKey

		result = append(result, &FireblocksAccount{
			Base: Base{
				ID:        account.ID.Hex(),
				IsActive:  account.IsActive,
				CreatedAt: account.CreatedAt,
				UpdatedAt: account.UpdatedAt,
			},
			VaultID:           account.VaultID,
			AssetID:           account.AssetID,
			Name:              account.Name,
			Alias:             account.Alias,
			PublicKey:         fbPubKey,
			PublicKeyFallback: account.PublicKey,
			WalletID:          account.WalletID,
			Wallet:            walletsMap[account.WalletID],
		})
	}

	return result, nil
}

func (fb *FireblocksService) FindFireblocksAccountByID(ctx context.Context, fireblocksAccountId string) (*FireblocksAccount, error) {
	fireblocksAccount, err := fb.repo.FindFireblocksAccountById(ctx, fireblocksAccountId)
	if err != nil {
		l.Logger.Error("fireblocks service: error finding fireblocks account", zap.Error(err))
		return nil, err
	}

	blockchains, err := fb.repo.FindBlockchains(ctx)
	if err != nil {
		l.Logger.Error("fireblocks service: error finding blockchains", zap.Error(err))
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

	wallets, err := fb.repo.FindWallets(ctx)
	if err != nil {
		l.Logger.Error("fireblocks service: error finding wallets", zap.Error(err))
		return nil, err
	}

	walletsMap := make(map[string]*Wallet)
	for _, wallet := range wallets {
		walletsMap[wallet.ID.Hex()] = &Wallet{
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
		}
	}

	// Get the public key info from the Fireblocks API
	fbPubKey := ""

	fbAccResult, err := fb.fbCli.GetPublicKeyInfoFromVaultAccount(ctx, fireblocksAccount.VaultID, fireblocksAccount.AssetID, 0, 0)
	if err != nil {
		l.Logger.Error("fireblocks error getting public key info from Fireblocks API", zap.Error(err))
	}

	fbPubKey = fbAccResult.PublicKey

	return &FireblocksAccount{
		Base: Base{
			ID:        fireblocksAccount.ID.Hex(),
			IsActive:  fireblocksAccount.IsActive,
			CreatedAt: fireblocksAccount.CreatedAt,
			UpdatedAt: fireblocksAccount.UpdatedAt,
		},
		VaultID:           fireblocksAccount.VaultID,
		AssetID:           fireblocksAccount.AssetID,
		Name:              fireblocksAccount.Name,
		Alias:             fireblocksAccount.Alias,
		PublicKey:         fbPubKey,
		PublicKeyFallback: fireblocksAccount.PublicKey,
		WalletID:          fireblocksAccount.WalletID,
		Wallet:            walletsMap[fireblocksAccount.WalletID],
	}, nil
}

func (fb *FireblocksService) FindFireblocksAccountByVaultID(ctx context.Context, vaultId string) (*FireblocksAccount, error) {
	fireblocksAccount, err := fb.repo.FindFireblocksAccountByVaultId(ctx, vaultId)
	if err != nil {
		l.Logger.Error("fireblocks service: error finding fireblocks account", zap.Error(err))
		return nil, err
	}

	blockchains, err := fb.repo.FindBlockchains(ctx)
	if err != nil {
		l.Logger.Error("fireblocks service: error finding blockchains", zap.Error(err))
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

	wallets, err := fb.repo.FindWallets(ctx)
	if err != nil {
		l.Logger.Error("fireblocks service: error finding wallets", zap.Error(err))
		return nil, err
	}

	walletsMap := make(map[string]*Wallet)
	for _, wallet := range wallets {
		walletsMap[wallet.ID.Hex()] = &Wallet{
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
		}
	}

	// Get the public key info from the Fireblocks API
	fbPubKey := ""

	fbAccResult, err := fb.fbCli.GetPublicKeyInfoFromVaultAccount(ctx, fireblocksAccount.VaultID, fireblocksAccount.AssetID, 0, 0)
	if err != nil {
		l.Logger.Error("fireblocks error getting public key info from Fireblocks API", zap.Error(err))
	}

	fbPubKey = fbAccResult.PublicKey

	return &FireblocksAccount{
		Base: Base{
			ID:        fireblocksAccount.ID.Hex(),
			IsActive:  fireblocksAccount.IsActive,
			CreatedAt: fireblocksAccount.CreatedAt,
			UpdatedAt: fireblocksAccount.UpdatedAt,
		},
		VaultID:           fireblocksAccount.VaultID,
		AssetID:           fireblocksAccount.AssetID,
		Name:              fireblocksAccount.Name,
		Alias:             fireblocksAccount.Alias,
		PublicKey:         fbPubKey,
		PublicKeyFallback: fireblocksAccount.PublicKey,
		WalletID:          fireblocksAccount.WalletID,
		Wallet:            walletsMap[fireblocksAccount.WalletID],
	}, nil
}

func (fb *FireblocksService) SaveFireblocksAccount(ctx context.Context, vaultID, assetID, walletID, name, alias, domain, publicKey string, accFlags int, isActive bool) (primitive.ObjectID, error) {
	if isValid := fb.repo.FireblocksAccountExists(ctx, vaultID, walletID, assetID, name, alias, publicKey, domain, accFlags, isActive); isValid {
		l.Logger.Error("fireblocks service: fireblocks account already exists", zap.Bool("is_active", isValid))
		return primitive.ObjectID{}, errors.New("fireblocks account already exists")
	}

	fireblocksAccount := &r.FireblocksAccount{
		VaultID: vaultID, AssetID: assetID, Name: name, Alias: alias, Domain: domain, PublicKey: publicKey, Flags: accFlags, IsActive: isActive, WalletID: walletID, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}

	id, err := fb.repo.SaveFireblocksAccount(ctx, fireblocksAccount)
	if err != nil {
		l.Logger.Error("fireblocks service: error saving fireblocks account", zap.Error(err))
		return primitive.ObjectID{}, err
	}

	return id, nil
}

func (fb *FireblocksService) DeleteFireblocks(ctx context.Context, fireblocksAccountId string) (string, error) {
	err := fb.repo.DeleteFireblocksAccount(ctx, fireblocksAccountId)
	if err != nil {
		l.Logger.Error("fireblocks service: error deleting fireblocksAccount", zap.Error(err))
		return "", err
	}

	return fireblocksAccountId, nil
}

func (fb *FireblocksService) EditFireblocksAccount(ctx context.Context, vaultID, assetID, walletID, name, alias, domain, publicKey string, accFlags int, isActive bool) (primitive.ObjectID, error) {
	fireblocksAccount := &r.FireblocksAccount{
		VaultID: vaultID, AssetID: assetID, Name: name, Alias: alias, Domain: domain, PublicKey: publicKey, Flags: accFlags, IsActive: isActive, WalletID: walletID, UpdatedAt: time.Now(),
	}

	id, err := fb.repo.EditFireblocksAccount(ctx, fireblocksAccount)
	if err != nil {
		l.Logger.Error("fireblocks service: error updating fireblocks account", zap.Error(err))
		return primitive.ObjectID{}, err
	}

	return id, nil
}
