package db

import (
	"fmt"
	"log"

	"github.com/Furkanberkay/ticket-booking-project-v1/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init(cfg *config.EnvConfig, Dbmigrator func(db *gorm.DB) error) *gorm.DB {
	uri := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.ServerPort, cfg.DBSslMode,
	)

	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v", err)
	}

	if err := DBMigrator(db); err != nil {
		log.Fatalf("Unable to migrate to the database: %v", err)
	}

	log.Println("connected to the database")
	return db
}
