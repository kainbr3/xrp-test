package api

import (
	cfg "crypto-braza-tokens-api/configs"
	r "crypto-braza-tokens-api/repositories"
	bs "crypto-braza-tokens-api/services/blockchain"
	fs "crypto-braza-tokens-api/services/fireblocks"
	hs "crypto-braza-tokens-api/services/health"
	ops "crypto-braza-tokens-api/services/operation"
	ts "crypto-braza-tokens-api/services/token"
	txs "crypto-braza-tokens-api/services/transaction"
	ws "crypto-braza-tokens-api/services/wallet"
	"fmt"
	"os"

	l "crypto-braza-tokens-api/utils/logger"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func Start() {
	// creates a new resources instance to be used on server execution
	repo := r.NewRepository()
	resources := &cfg.Resources{
		HealthService:      hs.NewHealthService(repo),
		FireblocksService:  fs.NewFireblocksService(repo),
		BlockchainService:  bs.NewBlockchainService(repo),
		TokenService:       ts.NewTokenService(repo),
		WalletService:      ws.NewWalletService(repo),
		OperationService:   ops.NewOperationService(repo),
		TransactionService: txs.NewTransactionService(repo),
	}

	// creates a new fiber instance
	app := fiber.New()

	// sets the api port to be used
	apiPort := "8000" // default port when not set
	if port := os.Getenv("API_PORT"); port != "" {
		apiPort = port
	}

	// adds middlewares to the api
	addMiddlewares(app)

	// creates all the endpoints to be served by the service
	buildRoutes(app, resources)

	// starts serving the api
	err := app.Listen(":" + apiPort)
	if err != nil {
		l.Logger.Fatal(fmt.Sprintf("failed serving api on port %s", apiPort), zap.Error(err))
	}
}
