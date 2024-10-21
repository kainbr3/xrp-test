package repositories

import (
	"context"
	kvs "crypto-braza-tokens-api/utils/keys-values"
	l "crypto-braza-tokens-api/utils/logger"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var repo *Repository

type Repository struct {
	database                     *mongo.Database
	blockchainsCollection        *mongo.Collection
	tokensCollection             *mongo.Collection
	walletsCollection            *mongo.Collection
	fireblocksAccountsCollection *mongo.Collection
	operationsCollection         *mongo.Collection
	operationsTypesCollection    *mongo.Collection
	operationsDomainsCollection  *mongo.Collection
	operationsLogsCollection     *mongo.Collection
	transactionsCollection       *mongo.Collection
	transactionsTypesCollection  *mongo.Collection
}

func NewRepository() *Repository {
	if repo != nil {
		return repo
	}

	connectionString := os.Getenv("MONGO_URI")
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		l.Logger.Fatal("repository: failed to connect to mongo instance", zap.Error(err))
	}

	brazaTokensDatabase, err := kvs.Get("MONGO_BRAZA_TOKENS_DATABASE")
	if err != nil {
		l.Logger.Fatal("repository: " + err.Error())
	}
	database := client.Database(brazaTokensDatabase)

	blockchainsCollection, err := kvs.Get("MONGO_BLOCKCHAINS_COLLECTION")
	if err != nil {
		l.Logger.Fatal("repository: " + err.Error())
	}
	blockchains := database.Collection(blockchainsCollection)

	tokensCollection, err := kvs.Get("MONGO_TOKENS_COLLECTION")
	if err != nil {
		l.Logger.Fatal("repository: " + err.Error())
	}
	tokens := database.Collection(tokensCollection)

	walletsCollection, err := kvs.Get("MONGO_WALLETS_COLLECTION")
	if err != nil {
		l.Logger.Fatal("repository: " + err.Error())
	}
	wallets := database.Collection(walletsCollection)

	fireblocksAccountsCollection, err := kvs.Get("MONGO_FIREBLOCKS_ACCOUNTS_COLLECTION")
	if err != nil {
		l.Logger.Fatal("repository: " + err.Error())
	}
	fireblocksAccounts := database.Collection(fireblocksAccountsCollection)

	operationsCollection, err := kvs.Get("MONGO_OPERATIONS_COLLECTION")
	if err != nil {
		l.Logger.Fatal("repository: " + err.Error())
	}
	operations := database.Collection(operationsCollection)

	operationsTypesCollection, err := kvs.Get("MONGO_OPERATIONS_TYPES_COLLECTION")
	if err != nil {
		l.Logger.Fatal("repository: " + err.Error())
	}
	operationsTypes := database.Collection(operationsTypesCollection)

	operationsDomainsCollection, err := kvs.Get("MONGO_OPERATIONS_DOMAINS_COLLECTION")
	if err != nil {
		l.Logger.Fatal("repository: " + err.Error())
	}
	operationsDomains := database.Collection(operationsDomainsCollection)

	operationsLogsCollection, err := kvs.Get("MONGO_OPERATIONS_LOGS_COLLECTION")
	if err != nil {
		l.Logger.Fatal("repository: " + err.Error())
	}
	operationsLogs := database.Collection(operationsLogsCollection)

	transactionsCollection, err := kvs.Get("MONGO_TRANSACTIONS_COLLECTION")
	if err != nil {
		l.Logger.Fatal("repository: " + err.Error())
	}
	transactions := database.Collection(transactionsCollection)

	transactionsTypesCollection, err := kvs.Get("MONGO_TRANSACTIONS_TYPES_COLLECTION")
	if err != nil {
		l.Logger.Fatal("repository: " + err.Error())
	}
	transactionsTypes := database.Collection(transactionsTypesCollection)

	repo = &Repository{
		database,
		blockchains,
		tokens,
		wallets,
		fireblocksAccounts,
		operations,
		operationsTypes,
		operationsDomains,
		operationsLogs,
		transactions,
		transactionsTypes,
	}

	return repo
}

func (r *Repository) CheckHealth(ctx context.Context) error {
	return r.database.Client().Ping(ctx, nil)
}
