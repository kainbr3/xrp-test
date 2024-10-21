package main

import (
	"crypto-braza-tokens-api/api"
	_ "crypto-braza-tokens-api/api/docs"
	cfg "crypto-braza-tokens-api/configs"
)

// @title           Braza Tokens API
// @version         1.0
// @description     A service to provide MINT/BURN operations for BRAZA tokens (BBRL and USDB)
// @termsOfService  https://www.brazaon.com.br
// @BasePath /

// @contact.name   Any member of getBRAZA Backend Team
// @contact.email  getbraza@braza.com.br
// @contact.url    https://www.brazaon.com.br

// @host      localhost:6000
// @BasePath  /

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	cfg.Startup()
	api.Start()
}
