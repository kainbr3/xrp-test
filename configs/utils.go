package configs

import (
	"log"
	"os"
)

func validateRequiredEnvs() {
	requiredEnvs := []string{
		"NAMESPACE",
		"API_PORT",
		"KVS_MONGO_URI",
		"KVS_DATABASE",
		"KVS_COLLECTION",
		"MONGO_URI",
		"FIREBLOCKS_API_KEY",
		"FIREBLOCKS_API_SECRET",
	}

	for _, env := range requiredEnvs {
		if os.Getenv(env) == "" {
			log.Fatalf("Required environment variable %s is not set", env)
		}
	}
}
