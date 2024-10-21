package transaction

import (
	"context"
	fb "crypto-braza-tokens-api/clients/fireblocks"
	xrpn "crypto-braza-tokens-api/clients/ripple"
	r "crypto-braza-tokens-api/repositories"
	l "crypto-braza-tokens-api/utils/logger"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type TransactionService struct {
	repo      *r.Repository
	fbClient  *fb.FireblocksClient
	xrpClient *xrpn.RippleNodeClient
}

func NewTransactionService(repo *r.Repository) *TransactionService {
	fbCli, err := fb.NewFireblocksClient()
	if err != nil {
		l.Logger.Fatal("transaction service: failed to create a new fireblocks client", zap.Error(err))
	}

	xrpCli, err := xrpn.NewRippleNodeClient()
	if err != nil {
		l.Logger.Fatal("transaction service: failed to create a new xrp node client", zap.Error(err))
	}

	return &TransactionService{repo, fbCli, xrpCli}
}

func (t *TransactionService) ExecuteInternalTransaction(ctx context.Context, domain, txType, blockchainId, assetId, amount, externalTxId string) (*fb.SubmittedTransactionResponse, error) {
	fbAccountsList, err := t.repo.FindFireblocksAccountByDomain(ctx, domain)
	if err != nil {
		l.Logger.Error("transaction service: error finding fireblocks accounts by domain", zap.Error(err))
		return nil, fmt.Errorf("error finding fireblocks accounts by domain: %s", err)
	}

	if len(fbAccountsList) == 0 {
		l.Logger.Error("transaction service: no fireblocks accounts found by domain", zap.String("domain", domain))
	}

	walletList, err := t.repo.FindWalletsByBlockchainAndDomain(ctx, blockchainId, domain)
	if err != nil {
		l.Logger.Error("transaction service: error finding wallets by blockchain id and domain", zap.Error(err))
		return nil, fmt.Errorf("error finding wallets by blockchain id %s and domain: %s", blockchainId, domain)
	}

	if len(walletList) == 0 {
		l.Logger.Error("transaction service: no wallets found by blockchain id and domain", zap.String("blockchain_id", blockchainId), zap.String("domain", domain))
		return nil, fmt.Errorf("no wallets found by blockchain id %s and domain: %s", blockchainId, domain)
	}

	walletsMap := make(map[string]*r.Wallet)

	for _, wallet := range walletList {
		walletsMap[wallet.ID.Hex()] = wallet
	}

	// builds the wallet params for the transaction
	typeFrom := "SUPPLY"
	if strings.EqualFold(txType, "OFF-RAMP") {
		typeFrom = "PAYMENT"
	}

	var fromParams *InternalTransactionParams
	var toParams *InternalTransactionParams

	for _, fbAccount := range fbAccountsList {
		wallet := walletsMap[fbAccount.WalletID]

		txParams := &InternalTransactionParams{
			Domain:          fbAccount.Domain,
			BlockchinID:     fbAccount.Blockchain,
			WalletID:        wallet.ID.Hex(),
			FbAccVaultID:    fbAccount.VaultID,
			FbAccAssetID:    fbAccount.AssetID,
			FbAccName:       fbAccount.Name,
			FbAccAlias:      fbAccount.Alias,
			Flags:           fbAccount.Flags,
			WalletPublicKey: fbAccount.PublicKey,
			WalletName:      wallet.Name,
			Address:         wallet.Address,
			WalletType:      wallet.Type,
		}

		if wallet.Type == typeFrom {
			fromParams = txParams
		} else {
			toParams = txParams
		}
	}

	sourceFbAccBalance, err := t.fbClient.GetVaultAccountAssetBalance(ctx, fromParams.FbAccVaultID, fromParams.FbAccAssetID)
	if err != nil {
		if strings.Contains(err.Error(), fb.INVALID_ASSET_CODE) {
			l.Logger.Error("transaction service: unsupported asset", zap.String("asset_id", assetId))
			return nil, fmt.Errorf("unsupported asset: %s", assetId)
		}

		l.Logger.Error("transaction service: error getting source account balance", zap.Error(err))
		return nil, fmt.Errorf("error getting source account balance: %s", err)
	}

	parsedBalance, err := strconv.ParseFloat(sourceFbAccBalance.Available, 64)
	if err != nil {
		l.Logger.Error("transaction service: error parsing balance", zap.Error(err))
		parsedBalance = 0
	}

	parsedAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		l.Logger.Error("transaction service: error parsing amount", zap.Error(err))
		parsedAmount = 0
	}

	if parsedBalance < parsedAmount {
		l.Logger.Error("transaction service: insufficient balance", zap.Float64("balance", parsedBalance), zap.Float64("amount", parsedAmount))
		return nil, fmt.Errorf("insufficient balance: %f to transfer amount %f", parsedBalance, parsedAmount)
	}

	if externalTxId == "" {
		l.Logger.Info("transaction service: external tx id is empty, generating a new one for this transaction")
		externalTxId = uuid.New().String()
	}

	note := fmt.Sprintf("transfering %s %s from %s to %s with external id: %s", amount, assetId, fromParams.FbAccName, toParams.FbAccName, externalTxId)
	l.Logger.Info("transaction service: executing internal transaction", zap.String("note", note))

	internalTxRequest := t.fbClient.BuildInternalTransactionRequest(ctx, fromParams.FbAccVaultID, toParams.FbAccVaultID, assetId, amount, note, externalTxId)

	result, err := t.fbClient.SubmitTransaction(ctx, internalTxRequest)
	if err != nil {
		l.Logger.Error("transaction service: error submitting transaction", zap.Error(err))
		return nil, fmt.Errorf("error submitting transaction: %s", err)
	}

	return result, nil
}

func (t *TransactionService) ExecuteWhitelistedTransaction() {}

func (t *TransactionService) ExecuteOneTimeAddressTransaction() {}

func (t *TransactionService) ValidateRippleTokenTrustSet(ctx context.Context, tokenId, validatingAddress string) (bool, error) {
	token, err := t.repo.FindTokenById(ctx, tokenId)
	if err != nil {
		l.Logger.Error("transaction service: error finding token by id", zap.Error(err))
		return false, err
	}

	accLines, err := t.xrpClient.GetAccountLines(ctx, validatingAddress)
	if err != nil {
		l.Logger.Error("transaction service: error getting account lines", zap.Error(err))
		return false, err
	}

	for _, addressLine := range accLines.Lines {
		if strings.EqualFold(addressLine.Account, token.Address) {
			return true, nil
		}
	}

	return false, nil
}
