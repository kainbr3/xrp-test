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

func (t *TransactionService) ExecuteInternalTransaction(ctx context.Context, opType, sourceVaultId, destinaitionVaultId, assetId, amount, externalTxId string) (*fb.SubmittedTransactionResponse, error) {
	//get account by op type and domain
	// sypply / payments / off-ramp / on-ramp

	fbAccountFrom, err := t.repo.FindFireblocksAccountByVaultId(ctx, sourceVaultId)
	if err != nil {
		l.Logger.Error("transaction service: error finding fireblocks source account by vault id", zap.Error(err))
		return nil, fmt.Errorf("error finding fireblocks source account by vault id: %s", err)
	}

	fbAccountTo, err := t.repo.FindFireblocksAccountByVaultId(ctx, destinaitionVaultId)
	if err != nil {
		l.Logger.Error("transaction service: error finding fireblocks destination account by vault id", zap.Error(err))
		return nil, fmt.Errorf("error finding fireblocks destination account by vault id: %s", err)
	}

	sourceFbAccBalance, err := t.fbClient.GetVaultAccountAssetBalance(ctx, sourceVaultId, "123456")
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

	note := fmt.Sprintf("transfering %s %s from %s to %s with external id: %s", amount, assetId, fbAccountFrom.Name, fbAccountTo.Name, externalTxId)
	l.Logger.Info("transaction service: executing internal transaction", zap.String("note", note))

	internalTxRequest := t.fbClient.BuildInternalTransactionRequest(ctx, sourceVaultId, destinaitionVaultId, assetId, amount, note, externalTxId)

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
