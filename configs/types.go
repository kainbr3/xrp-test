package configs

import (
	bs "crypto-braza-tokens-api/services/blockchain"
	fs "crypto-braza-tokens-api/services/fireblocks"
	hs "crypto-braza-tokens-api/services/health"
	ops "crypto-braza-tokens-api/services/operation"
	ts "crypto-braza-tokens-api/services/token"
	txs "crypto-braza-tokens-api/services/transaction"
	ws "crypto-braza-tokens-api/services/wallet"
	ow "crypto-braza-tokens-api/workers"
)

type Resources struct {
	HealthService      *hs.HealthService
	BlockchainService  *bs.BlockchainService
	FireblocksService  *fs.FireblocksService
	TokenService       *ts.TokenService
	WalletService      *ws.WalletService
	OperationService   *ops.OperationService
	TransactionService *txs.TransactionService
	Worker             *ow.OperationsWorker
}
