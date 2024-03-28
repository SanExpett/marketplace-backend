package main

import (
	"fmt"
	"github.com/SanExpett/marketplace-backend/internal/server"
	"github.com/SanExpett/marketplace-backend/pkg/config"
)

//	@title      MARKETPLACE project API
//	@version    1.0
//	@description  This is a server of MARKETPLACE server.
//
// @Schemes http
// @BasePath  /api/v1
func main() {
	configServer := config.New()

	srv := new(server.Server)
	if err := srv.Run(configServer); err != nil {
		fmt.Printf("Error in server: %s", err.Error())
	}
}
