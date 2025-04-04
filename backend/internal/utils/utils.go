package utils

import (
	"os"
	"strconv"
	"strings"
)

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
