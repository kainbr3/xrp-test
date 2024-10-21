package configs

import (
	k "crypto-braza-tokens-api/utils/keys-values"
	l "crypto-braza-tokens-api/utils/logger"
)

func Startup() {
	l.NewLogger()
	validateRequiredEnvs()
	k.Start()
}
