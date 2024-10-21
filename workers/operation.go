package worker

import (
	"context"
	fb "crypto-braza-tokens-api/clients/fireblocks"
	xrpn "crypto-braza-tokens-api/clients/ripple"
	binarycodec "crypto-braza-tokens-api/clients/ripple/utils/binary-codec"
	r "crypto-braza-tokens-api/repositories"
	l "crypto-braza-tokens-api/utils/logger"
	"fmt"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

type OperationsWorker struct {
	fbCli  *fb.FireblocksClient
	XrpCli *xrpn.RippleNodeClient
	repo   *r.Repository
	wg     sync.WaitGroup
}

func NewOperationsWorker(fbClient *fb.FireblocksClient, xrpClient *xrpn.RippleNodeClient, repository *r.Repository) (*OperationsWorker, error) {
	return &OperationsWorker{
		fbCli:  fbClient,
		XrpCli: xrpClient,
		repo:   repository,
	}, nil
}

func (o *OperationsWorker) Start(operationId string, rawTransaction map[string]any, callback func()) {
	o.wg.Add(1)

	go func() {
		defer func() {
			o.wg.Done()
		}()

		ctx := context.Background()
		operation, err := o.repo.FindOperationById(ctx, operationId)
		if err != nil {
			l.Logger.Error("operation worker: failed to find operation", zap.Error(err))
			return
		}

		var signedTx *fb.TransactionByIdResponse
		for {
			// retrieve the signed transaction status from fireblocks
			signedTx, err = o.fbCli.GetTransactionByID(ctx, operation.FireblocksId)
			if err != nil {
				l.Logger.Error("operation worker: failed to get transaction status from fireblocks", zap.Error(err))
				return
			}

			errLog := o.repo.SaveOperationLog(ctx, &r.OperationLog{
				Event:        "Fireblocks Raw Transaction Status Update",
				Description:  "Fireblocks Raw Transaction Status Response",
				OperationID:  operationId,
				FireblocksID: signedTx.ID,
				Payload:      "",
				Response:     signedTx,
				Error:        err,
			})
			if errLog != nil {
				l.Logger.Error("operation worker: failed to save operation log", zap.Error(errLog))
				return
			}

			// updates the operation status
			err = o.repo.UpdateOperationFireblocksStatus(ctx, operationId, signedTx.Status)
			if err != nil {
				l.Logger.Error("operation worker: failed to update operation status", zap.Error(err))
				return
			}

			if strings.EqualFold(signedTx.Status, "COMPLETED") {
				break
			}

			if strings.EqualFold(signedTx.Status, "FAILED") {
				l.Logger.Error("operation worker: transaction failed", zap.String("status", signedTx.Status))
				return
			}

			// Optionally, add a sleep to avoid hammering the API
			time.Sleep(5 * time.Second)
		}

		o.processOperation(ctx, operationId, signedTx, rawTransaction, callback)
	}()
}

func (o *OperationsWorker) processOperation(ctx context.Context, operationId string, signedTx *fb.TransactionByIdResponse, rawTransaction map[string]any, callback func()) {
	defer callback()

	rHex := signedTx.SignedMessages[0].Signature.R
	sHex := signedTx.SignedMessages[0].Signature.S
	derEncoded, err := xrpn.EncodeDER(rHex, sHex)
	if err != nil {
		l.Logger.Error("operation worker: failed to create a DER-encoded hexadecimal", zap.Error(err))
		return
	}

	// add the der encoded signature to the RAW tx payload
	rawTransaction["TxnSignature"] = derEncoded

	// encode the signed RAW transaction into a blob
	signedTxBlob, err := binarycodec.Encode(rawTransaction)
	if err != nil {
		l.Logger.Error("operation worker: failed to encode xrp tx into blob", zap.Error(err))
		return
	}

	contactedPrefixWithSignedTxBlob := xrpn.ConcactPrefixWithTxBlob(xrpn.PREFIX_SIGNED, signedTxBlob)
	hashedSignedTx, err := xrpn.Sha512Half(xrpn.HASH_SIZE, contactedPrefixWithSignedTxBlob)
	if err != nil {
		l.Logger.Error("operation worker: failed to computes the SHA-512 hash of the input hex string", zap.Error(err))
		return
	}

	txJsonRquest := o.XrpCli.BuildRawTransactionRequest(signedTxBlob)
	submitedTx, err := o.XrpCli.SubmitSignedTransaction(ctx, signedTxBlob, txJsonRquest)
	if err != nil {
		l.Logger.Error("operation worker: failed to submit signed transaction", zap.Error(err))
		return
	}

	// change to bclockchain things
	errLog := o.repo.SaveOperationLog(ctx, &r.OperationLog{
		Event:        "Submit XRP Signed Transaction",
		Description:  fmt.Sprintf("Submited XRP Signed Transaction to Ripple Node for Operation ID %s and Hash %s", operationId, hashedSignedTx),
		OperationID:  operationId,
		FireblocksID: signedTx.ID,
		Payload:      txJsonRquest,
		Response:     submitedTx,
		Error:        err,
	})
	if errLog != nil {
		l.Logger.Error("operation worker: failed to save operation log", zap.Error(errLog))
		return
	}

	hash := hashedSignedTx
	link := ""
	status := "FAILED"
	if strings.EqualFold(submitedTx.Result.EngineResult, "tesSUCCESS") {
		link = o.XrpCli.GetTransactionLink(hash)
		status = "COMPLETED"
	}

	err = o.repo.UpdateOperationBlockchainStatus(ctx, operationId, status, hash, link)
	if err != nil {
		l.Logger.Error("operation worker: failed to update operation status", zap.Error(err))
		return
	}

	l.Logger.Info(fmt.Sprintf("operation worker: operation %s completed with hash %s", operationId, hash), zap.String("details at:", link))
}
