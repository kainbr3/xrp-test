package api

import (
	h "crypto-braza-tokens-api/api/handlers"
	cfg "crypto-braza-tokens-api/configs"

	s "github.com/gofiber/swagger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// buildRoutes - setup api and creates all routes
func buildRoutes(app *fiber.App, resources *cfg.Resources) {
	// Default path validation and redirect
	app.Get("/", h.DefaultPath)

	// Sets the api base route
	api := app.Group("/api")

	baseRoutes(api, resources)
	v1Endpoints(app, resources)
}

// baseRoutes - creates base/root endpoints (health, swagger, etc)
func baseRoutes(api fiber.Router, resources *cfg.Resources) {
	// Open API Documentation Routes
	api.Get("/docs/*", s.HandlerDefault)

	// default health path to check if the application is running
	api.Get("/health", h.HealthHandler{Resources: resources}.Health)

	// liveness probe path to kubernetes check if application is alive (it always returns StatusCode 200)
	api.Get("/health/liveness", h.HealthHandler{Resources: resources}.HealthLiveness)

	// readiness probe path to check if the application is running with all of its dependencies
	api.Get("/health/readiness", h.HealthHandler{Resources: resources}.HealthReadiness)
}

func v1Endpoints(api fiber.Router, resources *cfg.Resources) {
	// Sets v1 route group
	v1 := api.Group("api/v1")

	// Blockchains
	v1.Get("/blockchains", h.BlockchainHandler{Resources: resources}.GetBlockchains)
	v1.Get("/blockchains/:id", h.BlockchainHandler{Resources: resources}.GetBlockchainByID)
	v1.Get("/blockchains/:id/tokens", h.BlockchainHandler{Resources: resources}.GetBlockchainTokens)
	v1.Post("/blockchains", h.BlockchainHandler{Resources: resources}.PostBlockchain)
	v1.Delete("/blockchains/:id", h.BlockchainHandler{Resources: resources}.DeleteBlockchain)
	v1.Patch("/blockchains", h.BlockchainHandler{Resources: resources}.PatchBlockchain)

	// Tokens
	v1.Get("/tokens", h.TokensHandler{Resources: resources}.GetTokens)
	v1.Get("/tokens/:id", h.TokensHandler{Resources: resources}.GetTokenByID)
	v1.Post("/tokens", h.TokensHandler{Resources: resources}.PostToken)
	v1.Delete("/tokens/:id", h.TokensHandler{Resources: resources}.DeleteToken)
	v1.Patch("/tokens", h.TokensHandler{Resources: resources}.PatchToken)

	// Wallets
	v1.Get("/wallets", h.WalletsHandler{Resources: resources}.GetWallets)
	v1.Get("/wallets/:id", h.WalletsHandler{Resources: resources}.GetWalletByID)
	v1.Get("/wallets/address/:address/blockchain/:blockchain_id", h.WalletsHandler{Resources: resources}.GetWalletByAddressAndBlockchain)
	v1.Get("/wallets/blockchain/:blockchain_id/type/:wallet_type/domain/:domain", h.WalletsHandler{Resources: resources}.GetWalletByBlockchainWalletTypeAndDomain)
	v1.Get("/wallets/blockchain/:blockchain_id/domain/:domain", h.WalletsHandler{Resources: resources}.GetWalletsByBlockchainAndDomain)
	v1.Get("/wallets/balances", h.WalletsHandler{Resources: resources}.GetWalletsBalances)
	v1.Post("/wallets", h.WalletsHandler{Resources: resources}.PostWallet)
	v1.Patch("/wallets", h.WalletsHandler{Resources: resources}.PatchWallet)
	v1.Delete("/wallets/:id", h.WalletsHandler{Resources: resources}.DeleteWallet)

	// Fireblocks Accounts
	v1.Get("/fireblocks-accounts", h.FireblocksAccountsHandler{Resources: resources}.GetFireblocksAccounts)
	v1.Get("/fireblocks-accounts/:id", h.FireblocksAccountsHandler{Resources: resources}.GetFireblocksAccountById)
	v1.Get("/fireblocks-accounts/vault/:vault_id", h.FireblocksAccountsHandler{Resources: resources}.GetFireblocksAccountByVaultId)
	v1.Post("/fireblocks-accounts", h.FireblocksAccountsHandler{Resources: resources}.PostFireblocksAccounts)
	v1.Patch("/fireblocks-accounts", h.FireblocksAccountsHandler{Resources: resources}.PatchFireblocksAccounts)
	v1.Delete("/fireblocks-accounts/:id", h.FireblocksAccountsHandler{Resources: resources}.DeleteFireblocksAccount)

	// Operations
	v1.Get("/operations", h.OperationsHandler{Resources: resources}.GetOperations)
	v1.Get("/operations/:id", h.OperationsHandler{Resources: resources}.GetOperationById)
	v1.Post("/operations", h.OperationsHandler{Resources: resources}.PostOperation)

	// Operation Types
	v1.Get("/operations-types/list", h.OperationsHandler{Resources: resources}.GetOperationTypesNames)
	v1.Get("/operations-types", h.OperationsHandler{Resources: resources}.GetOperationTypes)
	v1.Get("/operations-types/:id", h.OperationsHandler{Resources: resources}.GetOperationTypeById)
	v1.Post("/operations-types", h.OperationsHandler{Resources: resources}.PostOperationType)
	v1.Patch("/operations-types", h.OperationsHandler{Resources: resources}.PatchOperationType)
	v1.Delete("/operations-types/:id", h.OperationsHandler{Resources: resources}.DeleteOperationType)

	// Operation Domains
	v1.Get("/operations-domains/list", h.OperationsHandler{Resources: resources}.GetOperationDomainsNames)
	v1.Get("/operations-domains", h.OperationsHandler{Resources: resources}.GetOperationDomains)
	v1.Get("/operations-domains/:id", h.OperationsHandler{Resources: resources}.GetOperationDomainById)
	v1.Post("/operations-domains", h.OperationsHandler{Resources: resources}.PostOperationDomain)
	v1.Patch("/operations-domains", h.OperationsHandler{Resources: resources}.PatchOperationDomain)
	v1.Delete("/operations-domains/:id", h.OperationsHandler{Resources: resources}.DeleteOperationDomain)

	// Transactions
	v1.Get("/transactions", h.TransactionsHandler{Resources: resources}.GetTransactions)
	v1.Get("/transactions/:id", h.TransactionsHandler{Resources: resources}.GetTransactions)
	v1.Post("/transactions", h.TransactionsHandler{Resources: resources}.GetTransactions)
	v1.Post("/transactions/webhook", h.TransactionsHandler{Resources: resources}.PostWebhook)
	v1.Post("/transfers/webhook", h.TransactionsHandler{Resources: resources}.PostWebhook)

	// Transactions Assets
	v1.Get("/transactions-assets", h.TransactionsHandler{Resources: resources}.GetTransactions)
	v1.Get("/transactions-assets/:id", h.TransactionsHandler{Resources: resources}.GetTransactions)
	v1.Post("/transactions-assets", h.TransactionsHandler{Resources: resources}.GetTransactions)
	v1.Patch("/transactions-assets", h.TransactionsHandler{Resources: resources}.GetTransactions)
	v1.Delete("/transactions-assets", h.TransactionsHandler{Resources: resources}.GetTransactions)

	// Transactions Types
	v1.Get("/transactions-types/list", h.TransactionsTypesHandler{Resources: resources}.GetTransactionsTypesNames)
	v1.Get("/transactions-types", h.TransactionsTypesHandler{Resources: resources}.GetTransactionsTypes)
	v1.Get("/transactions-types/:id", h.TransactionsTypesHandler{Resources: resources}.GetTransactionTypeById)
	v1.Post("/transactions-types", h.TransactionsTypesHandler{Resources: resources}.PostTransactionType)
	v1.Patch("/transactions-types/:id", h.TransactionsTypesHandler{Resources: resources}.PatchTransactionType)
	v1.Delete("/transactions-types/:id", h.TransactionsTypesHandler{Resources: resources}.DeleteTransactionType)
}

func addMiddlewares(app *fiber.App) {
	setupCors(app)
	setupAuthorization(app)
}

func setupCors(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowHeaders: "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))
}

func setupAuthorization(app *fiber.App) {
	app.Use(h.AuthorizerHandler)
}
