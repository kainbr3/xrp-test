package operation

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	fb "crypto-braza-tokens-api/clients/fireblocks"
	xrpn "crypto-braza-tokens-api/clients/ripple"
	binarycodec "crypto-braza-tokens-api/clients/ripple/utils/binary-codec"
	r "crypto-braza-tokens-api/repositories"
	l "crypto-braza-tokens-api/utils/logger"
	ow "crypto-braza-tokens-api/workers"

	"go.uber.org/zap"
)

type OperationService struct {
	repo      *r.Repository
	fbClient  *fb.FireblocksClient
	xrpClient *xrpn.RippleNodeClient
	worker    *ow.OperationsWorker
}

func NewOperationService(repo *r.Repository) *OperationService {
	fbCli, err := fb.NewFireblocksClient()
	if err != nil {
		l.Logger.Fatal("operation service: failed to create a new fireblocks client", zap.Error(err))
	}

	xrpCli, err := xrpn.NewRippleNodeClient()
	if err != nil {
		l.Logger.Fatal("operation service: failed to create a new xrp node client", zap.Error(err))
	}

	worker, err := ow.NewOperationsWorker(fbCli, xrpCli, repo)
	if err != nil {
		l.Logger.Fatal("operation service: failed to create a new worker", zap.Error(err))
	}

	return &OperationService{repo, fbCli, xrpCli, worker}
}

func (o *OperationService) GetOperations(ctx context.Context) ([]*r.Operation, error) {
	operations, err := o.repo.FindOperations(ctx)
	if err != nil {
		l.Logger.Error("operation service: failed to find operations", zap.Error(err))
		return nil, err
	}

	if len(operations) == 0 {
		l.Logger.Error("operation service: no operations found")
		return nil, errors.New("no operations found")
	}

	return operations, nil
}

func (o *OperationService) GetPaginatedOperations(ctx context.Context, params *r.QueryParams) (*r.PaginatedOperations, error) {
	operations, err := o.repo.FindFilteredAndPaginatedOperations(ctx, params)
	if err != nil {
		l.Logger.Error("operation service: failed to find operations", zap.Error(err))
		return nil, err
	}

	if len(operations.Data) == 0 {
		l.Logger.Error("operation service: no operations found")
		return nil, errors.New("no operations found")
	}

	return operations, nil
}

func (o *OperationService) GetOperationById(ctx context.Context, id string) (*OperationWithLogs, error) {
	operation, err := o.repo.FindOperationById(ctx, id)
	if err != nil {
		l.Logger.Error("operation service: failed to find operation", zap.Error(err))
		return nil, err
	}

	logs, err := o.repo.FindOperationLogsByOperationId(ctx, id)
	if err != nil {
		l.Logger.Error("operation service: failed to find operation logs", zap.Error(err))
		return nil, err
	}

	result := &OperationWithLogs{
		Operation: *operation,
		Logs:      logs,
	}

	return result, nil
}

func (o *OperationService) ValidateParams(ctx context.Context, opType, opDomain, tokenId, blockchainId string) error {
	errorMessage := "%s not found"

	if isValid := o.repo.OpTypeExists(ctx, opType); !isValid {
		return fmt.Errorf(errorMessage, "operation type")
	}

	if isValid := o.repo.DomainExists(ctx, opDomain); !isValid {
		return fmt.Errorf(errorMessage, "operation domain")
	}

	if isValid := o.repo.TokenExists(ctx, tokenId); !isValid {
		return fmt.Errorf("token with id %s not found or not supported", tokenId)
	}

	if isValid := o.repo.BlockchainExists(ctx, blockchainId); !isValid {
		return fmt.Errorf(errorMessage, fmt.Sprintf("blockchain with id %s", blockchainId))
	}

	l.Logger.Info("operation service: input params are valid",
		zap.String("operation type", opType),
		zap.String("operation domain", opDomain),
		zap.String("token id", tokenId),
		zap.String("blockchain id", blockchainId),
	)

	return nil
}

func (o *OperationService) ExecuteOperation(ctx context.Context, opType, opDomain, tokenId, blockchainId, amount, operator string, callback func()) (string, error) {
	// retrieve blockchain info for the operation
	blockchain, err := o.repo.FindBlockchainById(ctx, blockchainId)
	if err != nil {
		l.Logger.Error("operation service: failed to find blockchain", zap.Error(err))
		return "", err
	}

	// retrieve token info for the operation
	token, err := o.repo.FindTokenById(ctx, tokenId)
	if err != nil {
		l.Logger.Error("operation service: failed to find token", zap.Error(err))
		return "", err
	}

	// builds the wallet params for the operation
	domainFrom := token.Abbr
	domainTo := opDomain
	typeFrom := "ISSUER"
	typeTo := "SUPPLY"

	if strings.EqualFold(opType, "BURN") {
		domainFrom = opDomain
		domainTo = token.Abbr
		typeFrom = "SUPPLY"
		typeTo = "ISSUER"
	}

	// retrieve origin wallet for the operation
	walletFrom, err := o.repo.FindWalletByBlockchainWalletTypeAndDomain(ctx, blockchain.ID.Hex(), typeFrom, domainFrom)
	if err != nil {
		l.Logger.Error("operation service: failed to find wallet", zap.Error(err))
		return "", err
	}

	// retrieve destination wallet for the operation
	walletTo, err := o.repo.FindWalletByBlockchainWalletTypeAndDomain(ctx, blockchain.ID.Hex(), typeTo, domainTo)
	if err != nil {
		l.Logger.Error("operation service: failed to find wallet", zap.Error(err))
		return "", err
	}

	issuerAddress := walletFrom.Address
	if strings.EqualFold(opType, "BURN") {
		issuerAddress = walletTo.Address
	}

	// retrieve fireblocks account for the origin wallet
	fbAccountFrom, err := o.repo.FindFireblocksAccountByWalletId(ctx, walletFrom.ID.Hex())
	if err != nil {
		l.Logger.Error("operation service: failed to find fireblocks account", zap.Error(err))
		return "", err
	}

	// create the operation object and store it to futher update and trackings
	operation := &r.Operation{
		Type:             opType,
		Domain:           opDomain,
		Amount:           amount,
		Operator:         operator,
		FireblocksStatus: "",
		BlockchainStatus: "",
		FireblocksId:     "",
		TransactionHash:  "",
		TransactionLink:  "",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	operationId, err := o.repo.SaveOperation(ctx, operation)
	if err != nil {
		l.Logger.Error("operation service: failed to save operation", zap.Error(err))
		return "", err
	}

	msg := fmt.Sprintf("New %s Operation of %s %s tokens from %s to %s", opType, amount, token.Abbr, walletFrom.Name, walletTo.Name)
	l.Logger.Info(msg)

	operationLog := &r.OperationLog{
		Event:        "Operation Started",
		Description:  msg,
		OperationID:  operationId.Hex(),
		FireblocksID: "",
		Payload:      parseStructToJson(operation),
		Response:     "",
		Error:        parseStructToJson(err),
		CreatedAt:    time.Now(),
	}

	if err := o.repo.SaveOperationLog(ctx, operationLog); err != nil {
		l.Logger.Error("operation service: failed to save operation log", zap.Error(err))
		return "", err
	}

	// retrieve fireblocks account pubkey for the origin wallet
	fbAccResult, err := o.fbClient.GetPublicKeyInfoFromVaultAccount(ctx, fbAccountFrom.VaultID, fbAccountFrom.AssetID, 0, 0)

	operationLog = &r.OperationLog{
		Event:        "Retrieve Fireblocks Account Public Key",
		Description:  fmt.Sprintf("Retrieve Fireblocks Acc PubKey for wallet %s and address %s of domain %s", walletFrom.Name, walletFrom.Address, domainFrom),
		OperationID:  operationId.Hex(),
		FireblocksID: "",
		Payload:      fmt.Sprintf("Fireblocks Account ID: %s, Asset ID: %s Change: %d Address Index: %d", fbAccountFrom.VaultID, fbAccountFrom.AssetID, 0, 0),
		Response:     fbAccResult,
		Error:        parseStructToJson(err),
		CreatedAt:    time.Now(),
	}

	if errLog := o.repo.SaveOperationLog(ctx, operationLog); errLog != nil {
		return "", errLog
	}

	if err != nil {
		l.Logger.Error("operation service: failed to get public key info from fireblocks", zap.Error(err))
		return "", err
	}

	// replace the public key for the updated one retrieved from fireblocks if it is not empty
	if fbAccResult != nil && fbAccResult.PublicKey != "" {
		fbAccountFrom.PublicKey = fbAccResult.PublicKey
	}

	// retrieve xrp account info for the origin wallet
	accNodeInfo, err := o.xrpClient.GetAccountInfo(ctx, walletFrom.Address)

	operationLog = &r.OperationLog{
		Event:        "Retrieve Account Params XRP Blockchain",
		Description:  fmt.Sprintf("Retrieve Account Params Info for address %s from XRP Blockchain Node API", walletFrom.Address),
		OperationID:  operationId.Hex(),
		FireblocksID: "",
		Payload:      fmt.Sprintf("Wallet %s, Address %s", walletFrom.Name, walletFrom.Address),
		Response:     parseStructToJson(accNodeInfo),
		Error:        parseStructToJson(err),
		CreatedAt:    time.Now(),
	}

	if errLog := o.repo.SaveOperationLog(ctx, operationLog); errLog != nil {
		return "", errLog
	}

	if err != nil {
		l.Logger.Error("operation service: failed to get account info from xrp node", zap.Error(err))
		return "", err
	}

	// builds the note message to be sent to fireblocks authorizers who will sign the RAW transaction
	note := fmt.Sprintf("%s %s %s tokens from %s to %s", opType, amount, token.Abbr, walletFrom.Name, walletTo.Name)
	l.Logger.Info(note)

	// builds the base payload for the RAW transaction
	rawTransactionBasePayload := buildRippleRawTransactionPayload(walletFrom.Address, walletTo.Address, token.Abbr, issuerAddress, amount, fbAccountFrom.PublicKey, fbAccountFrom.Flags, accNodeInfo.Result.AccountData.Sequence, accNodeInfo.Result.LedgerCurrentIndex)

	// encode the unsigned RAW transaction into a blob
	unsignTxBlob, err := binarycodec.Encode(rawTransactionBasePayload)
	if err != nil {
		l.Logger.Error("operation service: failed to encode xrp tx into blob", zap.Error(err))
		return "", err
	}

	// hashes the tx blob into 32 bytes message content for fireblocks raw sign
	contactedPrefixWithUnsignedTxBlob := xrpn.ConcactPrefixWithTxBlob(xrpn.PREFIX_UNSIGNED, unsignTxBlob)

	hasheUnsignedTx, err := xrpn.Sha512Half(xrpn.HASH_SIZE, contactedPrefixWithUnsignedTxBlob)
	if err != nil {
		l.Logger.Error("operation service: failed to computes the SHA-512 hash of the input hex string", zap.Error(err))
		return "", err
	}

	// build the raw transaction request to be submitted to fireblocks
	rawTxRequest := o.fbClient.BuildRawTransactionRequest(ctx, fbAccountFrom.VaultID, fbAccountFrom.AssetID, note, hasheUnsignedTx)

	// submit the raw transaction to fireblocks to be signed
	createRawTxResult, err := o.fbClient.SubmitTransaction(ctx, rawTxRequest)

	errLog := o.repo.SaveOperationLog(ctx, &r.OperationLog{
		Event:        "Fireblocks Raw Transaction Submitted",
		Description:  "Fireblocks Raw Transaction Submitted to be signed by authorizers",
		OperationID:  operationId.Hex(),
		FireblocksID: "",
		Payload:      parseStructToJson(rawTxRequest),
		Response:     "",
		Error:        parseStructToJson(err),
	})
	if errLog != nil {
		return "", errLog
	}

	if err != nil {
		l.Logger.Error("operation service: failed to submit raw transaction to fireblocks", zap.Error(err))
		return "", err
	}

	operationToUpdate, err := o.repo.FindOperationById(ctx, operationId.Hex())
	if err != nil {
		l.Logger.Error("operation service: failed to find operation", zap.Error(err))
		return "", err
	}

	operationToUpdate.FireblocksId = createRawTxResult.ID
	operationToUpdate.FireblocksStatus = createRawTxResult.Status
	operationToUpdate.UpdatedAt = time.Now()

	err = o.repo.UpdateOperationFireblocksIdAndStatus(ctx, operationToUpdate.ID.Hex(), operationToUpdate.FireblocksId, operationToUpdate.FireblocksStatus)
	if err != nil {
		l.Logger.Error("operation service: failed to update operation", zap.Error(err))
		return "", err
	}

	// start a worker to check the signed transaction status and submit it to the ripple network
	go o.worker.Start(operationId.Hex(), rawTransactionBasePayload, callback)

	l.Logger.Info(fmt.Sprintf("operation service: starting worker for operation %s", operationId))

	return operationId.Hex(), nil
}
