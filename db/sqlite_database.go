package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Furkanberkay/ticket-booking-project-v1/config"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init(cfg *config.EnvSQLiteConfig, migrator func(db *gorm.DB) error) *gorm.DB {

	driver := strings.ToLower(strings.TrimSpace(cfg.DBDriver))

	var (
		dbConn *gorm.DB
		err    error
	)

	switch driver {
	case "sqlite":
		if dir := filepath.Dir(cfg.SQLitePath); dir != "." {
			_ = os.MkdirAll(dir, 0755)
		}
		dbConn, err = gorm.Open(sqlite.Open(cfg.SQLitePath), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

	case "postgres":
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSslMode,
		)
		dbConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

	default:
		log.Fatalf("unknown DB_DRIVER: %s (use sqlite or postgres)", cfg.DBDriver)
	}

	if err != nil {
		log.Fatalf("Unable to connect to the database: %v", err)
	}

	if err := migrator(dbConn); err != nil {
		log.Fatalf("Unable to migrate to the database: %v", err)
	}

	log.Printf("connected to the database (%s)", driver)
	return dbConn
}
