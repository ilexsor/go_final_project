package utils

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ilexsor/internal/models"
	"gorm.io/gorm"
)

// GetServerPort функция получения номера порта из переменной окружения TODO_PORT
// Значение по-усолчанию :7540
func GetServerPort() string {
	defaultPort := ":7540"

	port := os.Getenv("TODO_PORT")
	if port == "" {
		return defaultPort
	}

	if strings.HasPrefix(port, ":") {
		portPart := port[1:] // убираем двоеточие для проверки номера порта
		if _, err := strconv.Atoi(portPart); err != nil {
			return defaultPort
		}
		return port
	}

	// Если нет ":" в начале, проверяем что это просто число
	if _, err := strconv.Atoi(port); err != nil {
		return defaultPort
	}

	// Если порт задан как "8080", добавляем ":" в начало
	return ":" + port
}

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
