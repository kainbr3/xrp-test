package keysvalues

import (
	"context"
	"fmt"
	"os"
	"time"

	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var kvs *Kvs

type Kvs struct {
	data map[string]string
	repo *kvsRepository
}

func Start() {
	initRepo()
	Refresh()
}

func Get(key string) (string, error) {
	if kvs == nil {
		return "", fmt.Errorf("keys values not initialized")
	}

	if value, exists := kvs.data[key]; exists {
		return value, nil
	}

	return "", fmt.Errorf("key %s not found", key)
}

func Refresh() {
	if err := kvs.repo.find(); err != nil {
		log.Fatalf("failed to retrieve keys values: %v", err)
	}
}

type kvsRepository struct {
	collection *mongo.Collection
	namespace  string
}

type KeyValue struct {
	Key   string `bson:"key"`
	Value string `bson:"value"`
}

func initRepo() {
	if kvs == nil || kvs != nil && kvs.repo == nil {

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		connectionString := os.Getenv("KVS_MONGO_URI")
		clientOptions := options.Client().ApplyURI(connectionString)
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			panic(err)
		}

		database := client.Database(os.Getenv("KVS_DATABASE"))
		collection := database.Collection(os.Getenv("KVS_COLLECTION"))
		namespace := os.Getenv("NAMESPACE")

		if err := database.Client().Ping(ctx, nil); err != nil {
			panic(fmt.Errorf("failed to establish connection to mongo: %v", err))
		}

		kvs = &Kvs{
			data: make(map[string]string),
			repo: &kvsRepository{collection, namespace},
		}
	}
}

func (k *kvsRepository) find() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"namespace": k.namespace}
	findOptions := options.Find().SetSort(bson.M{"key": 1})

	cursor, err := k.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return fmt.Errorf("error finding keys values: %v", err)
	}
	defer cursor.Close(ctx)

	var result []*KeyValue
	if err = cursor.All(ctx, &result); err != nil {
		return fmt.Errorf("error parsing keys values result: %v", err)
	}

	if len(result) == 0 {
		log.Fatal("no key values found")
	}

	for _, kv := range result {
		kvs.data[kv.Key] = kv.Value
	}

	return nil
}
