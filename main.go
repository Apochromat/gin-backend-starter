package main

import (
	"log"
	"os"
	"time"

	"github.com/imshawan/gin-backend-starter/configs"
	"github.com/imshawan/gin-backend-starter/infra/database"
	"github.com/imshawan/gin-backend-starter/infra/metrics"
	"github.com/imshawan/gin-backend-starter/routers"
	"github.com/spf13/viper"
)

func main() {
	// Pre-start, setup intial configurations before the server starts up
	if err := configs.SetupConfig(); err != nil {
		log.Printf("Error while setting up server config, %s", err)
	}

	viper.SetDefault("SERVER_TIMEZONE", "Asia/Kolkata")
	loc, _ := time.LoadLocation(viper.GetString("SERVER_TIMEZONE"))
	time.Local = loc

	// Initialize the database connection before the server starts
	if err := database.SetupDbConnection(); err != nil {
		log.Printf("Error while initializing database connection, %s", err)
		os.Exit(1)
	}

	registry := metrics.SetupMetrics()

	router := routers.SetupRouters(registry)

	// Run the server with the configurations
	err := router.Run(configs.ServerConfig())

	if err != nil {
		log.Printf("Error while starting up the server, %s", err)
		os.Exit(1)
	}
}
