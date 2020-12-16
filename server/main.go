package main

import (
	"github.com/kuops/go-example-app/server/cmd"
	_ "github.com/kuops/go-example-app/server/docs"
	"os"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-token

// @host localhost:8080
// @BasePath /
func main() {
	command := cmd.NewServerCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}