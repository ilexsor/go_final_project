package database

import (
	"fmt"
	"os"
	"time"

	"github.com/ilexsor/internal/models"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Миграция структуры в БД
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.Scheduler{})
}

// ConfigureDB функция для конфигурации соединений к БД
func ConfigureDB(dataBase *gorm.DB) {
	sqliteDB, _ := dataBase.DB()
	sqliteDB.SetMaxOpenConns(1)
	sqliteDB.SetMaxIdleConns(0)
	sqliteDB.SetConnMaxLifetime(time.Minute * 5)
}

// NewSqliteDB функция для создания БД и миграции
func NewSqliteDB(dbPath string) (*gorm.DB, error) {

	_, err := os.Stat(dbPath)
	// Если БД не создана, то создаем, делаем миграцию и возвращаем объект БД
	if err != nil {
		db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if err != nil {

			log.WithFields(log.Fields{
				"OpenDB": "error",
			}).Errorf("can`t open DB %v", dbPath)

			return nil, fmt.Errorf("failed to connect database: %w", err)
		}

		log.WithFields(log.Fields{
			"migrtation": "started",
		}).Infof("starting migration %v", dbPath)

		err = Migrate(db)
		if err != nil {

			log.WithFields(log.Fields{
				"MigrationDB": "error",
			}).Errorf("can`t migrate DB %v", dbPath)

			return nil, fmt.Errorf("failed to migrate database: %w", err)
		}

		ConfigureDB(db)
		return db, nil
	}

	// Если файл БД существует, то открываем его и возвращаем
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {

		log.WithFields(log.Fields{
			"OpenDB": "error",
		}).Errorf("can`t open DB %v", dbPath)

		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	ConfigureDB(db)
	return db, nil

}
