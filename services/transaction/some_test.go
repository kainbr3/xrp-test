package transaction

import (
	"context"
	r "crypto-braza-tokens-api/repositories"
	keysvalues "crypto-braza-tokens-api/utils/keys-values"
	"crypto-braza-tokens-api/utils/logger"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSomething(t *testing.T) {
	keysvalues.Start()
	repo := r.NewRepository()
	logger.NewLogger()
	cli := NewTransactionService(repo)

	// test1, err := repo.FindTransactionTypes(context.Background())
	// fmt.Println(err)
	// fmt.Println(test1)

	// validation, err := cli.ValidateRippleTokenTrustSet(context.Background(), "66f74ad8ba6b56108cb3e80b", "rng5ZxeWue9pggAPuGHZKXkYQQBmspKTdZ")
	//validation, err := cli.ValidateRippleTokenTrustSet(context.Background(), "66ff66c897875b4fe72e173b", "rng5ZxeWue9pggAPuGHZKXkYQQBmspKTdZ")
	// fmt.Println(err)
	// fmt.Println(validation)
	//cli.ExecuteInternalTransaction(context.Background(), "ON-RAMP", "17", "18", "XRP_TEST", "3", "")
	externalId := uuid.New().String()
	//result, err := cli.ExecuteInternalTransaction(context.Background(), "GET-BRAZA", "OFF-RAMP", "66f6fe7eccc6398d39e981f9", "XRP_TEST", "3", externalId)
	result, err := cli.ExecuteInternalTransaction(context.Background(), "GET-BRAZA", "ON-RAMP", "66f6fe7eccc6398d39e981f9", "XRP_TEST", "3", externalId)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}
