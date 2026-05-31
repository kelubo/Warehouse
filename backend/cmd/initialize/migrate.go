package initialize

import (
	"flag"
	"log"

	"warehouse-management/backend/config"
	"warehouse-management/backend/database"
	"warehouse-management/backend/utils"
)

func Migrate() {
	flag.Parse()

	cfg := config.LoadConfig()

	if err := database.InitDB(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	utils.Info("Migration completed successfully")
}