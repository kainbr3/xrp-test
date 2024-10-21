package transaction

import (
	"context"
	r "crypto-braza-tokens-api/repositories"
	keysvalues "crypto-braza-tokens-api/utils/keys-values"
	"crypto-braza-tokens-api/utils/logger"
	"fmt"
	"testing"
)

func TestSomething(t *testing.T) {
	keysvalues.Start()
	repo := r.NewRepository()
	logger.NewLogger()
	cli := NewTransactionService(repo)

	test1, err := repo.FindTransactionTypes(context.Background())
	fmt.Println(err)
	fmt.Println(test1)

	validation, err := cli.ValidateRippleTokenTrustSet(context.Background(), "66f74ad8ba6b56108cb3e80b", "rng5ZxeWue9pggAPuGHZKXkYQQBmspKTdZ")
	//validation, err := cli.ValidateRippleTokenTrustSet(context.Background(), "66ff66c897875b4fe72e173b", "rng5ZxeWue9pggAPuGHZKXkYQQBmspKTdZ")
	fmt.Println(err)
	fmt.Println(validation)
	//cli.ExecuteInternalTransaction(context.Background(), "ON-RAMP", "17", "18", "XRP_TEST", "3", "")
	cli.ExecuteInternalTransaction(context.Background(), "OFF-RAMP", "18", "17", "XRP_TEST", "3", "")
}
