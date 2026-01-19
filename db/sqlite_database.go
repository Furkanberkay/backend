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

func Init(cfg *config.EnvConfig, migrator func(db *gorm.DB) error) *gorm.DB {

	driver := strings.ToLower(strings.TrimSpace(cfg.DBDriver))

	var (
		dbConn *gorm.DB
		err    error
	)

	switch driver {
	case "sqlite":
		dbPath := cfg.SQLitePath
		if dbPath == "" {
			dbPath = "data/ticket_booking.db"
		}

		if dir := filepath.Dir(dbPath); dir != "." {
			_ = os.MkdirAll(dir, 0755)
		}
		dbConn, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
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
		log.Fatalf("Bilinmeyen DB_DRIVER: %s (sqlite veya postgres kullanın)", cfg.DBDriver)
	}

	if err != nil {
		log.Fatalf("Veritabanına bağlanılamadı: %v", err)
	}

	if err := migrator(dbConn); err != nil {
		log.Fatalf("Migration hatası: %v", err)
	}

	log.Printf("Veritabanına bağlanıldı (%s)", driver)
	return dbConn
}
