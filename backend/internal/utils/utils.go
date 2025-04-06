package utils

import (
	"os"
	"strconv"
	"strings"
	"time"
	"errors"

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

// NextDate Вычисляет следующую дату для задачи
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("пустое правило повторения")
	}

	startDate, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", errors.New("некорректная дата начала")
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("неверный формат правила повторения")
	}

	rule := parts[0]
	switch rule {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("неверный формат для правила 'd'")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", errors.New("неверное количество дней")
		}
		if days <= 0 || days > 400 {
			return "", errors.New("недопустимый интервал в днях")
		}
		return nextDateDaily(now, startDate, days), nil
	case "y":
		if len(parts) != 1 {
			return "", errors.New("неверный формат для правила 'y'")
		}
		return nextDateYearly(now, startDate), nil
	case "w":
		if len(parts) != 2 {
			return "", errors.New("неверный формат для правила 'w'")
		}
		days, err := parseWeekDays(parts[1])
		if err != nil {
			return "", err
		}
		return nextDateWeekly(now, startDate, days)
	case "m":
		if len(parts) < 2 {
			return "", errors.New("неверный формат для правила 'm'")
		}
		monthDays, months, err := parseMonthRule(parts[1:])
		if err != nil {
			return "", err
		}
		return nextDateMonthly(now, startDate, monthDays, months)
	default:
		return "", errors.New("неподдерживаемый формат правила повторения")
	}
}

func nextDateDaily(now, startDate time.Time, days int) string {
	date := startDate
	for !date.After(now) {
		date = date.AddDate(0, 0, days)
	}
	return date.Format("20060102")
}

func nextDateYearly(now, startDate time.Time) string {
	date := startDate
	for !date.After(now) {
		date = date.AddDate(1, 0, 0)
	}
	return date.Format("20060102")
}

func parseWeekDays(s string) ([]int, error) {
	parts := strings.Split(s, ",")
	days := make([]int, 0, len(parts))
	for _, part := range parts {
		d, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil || d < 1 || d > 7 {
			return nil, errors.New("недопустимый день недели")
		}
		days = append(days, d)
	}
	return days, nil
}

func nextDateWeekly(now, startDate time.Time, weekDays []int) (string, error) {
	date := startDate
	for !date.After(now) {
		// Находим следующий день недели
		found := false
		for i := 1; i <= 7; i++ {
			next := date.AddDate(0, 0, i)
			wd := int(next.Weekday())
			if wd == 0 {
				wd = 7 // Воскресенье как 7
			}
			for _, d := range weekDays {
				if wd == d {
					date = next
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			return "", errors.New("не удалось найти следующий день недели")
		}
	}
	return date.Format("20060102"), nil
}

func parseMonthRule(parts []string) ([]int, []int, error) {
	if len(parts) == 0 || len(parts) > 2 {
		return nil, nil, errors.New("неверный формат для правила 'm'")
	}

	// Парсим дни месяца
	dayParts := strings.Split(parts[0], ",")
	days := make([]int, 0, len(dayParts))
	for _, part := range dayParts {
		d, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return nil, nil, errors.New("недопустимый день месяца")
		}
		if d != -1 && d != -2 && (d < 1 || d > 31) {
			return nil, nil, errors.New("недопустимый день месяца")
		}
		days = append(days, d)
	}

	// Парсим месяцы, если они есть
	var months []int
	if len(parts) == 2 {
		monthParts := strings.Split(parts[1], ",")
		months = make([]int, 0, len(monthParts))
		for _, part := range monthParts {
			m, err := strconv.Atoi(strings.TrimSpace(part))
			if err != nil || m < 1 || m > 12 {
				return nil, nil, errors.New("недопустимый месяц")
			}
			months = append(months, m)
		}
	}

	return days, months, nil
}

func nextDateMonthly(now, startDate time.Time, monthDays, months []int) (string, error) {
	date := startDate
	for !date.After(now) {
		// Получаем год и месяц текущей даты
		year, month, _ := date.Date()
		currentMonth := int(month)
		currentDay := date.Day()

		// Проверяем, есть ли ограничение по месяцам
		validMonths := months
		if len(validMonths) == 0 {
			validMonths = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
		}

		// Находим следующий месяц
		nextMonth := currentMonth
		nextYear := year
		found := false
		for {
			nextMonth++
			if nextMonth > 12 {
				nextMonth = 1
				nextYear++
			}
			for _, m := range validMonths {
				if nextMonth == m {
					found = true
					break
				}
			}
			if found {
				break
			}
			if nextMonth == currentMonth && nextYear > year+1 {
				return "", errors.New("не удалось найти следующий месяц")
			}
		}

		// Находим следующий день в этом месяце
		var nextDay int
		if len(monthDays) == 0 {
			nextDay = currentDay
		} else {
			// Находим минимальный день, который больше текущего
			minDay := 32
			for _, d := range monthDays {
				actualDay := d
				if d == -1 || d == -2 {
					// Вычисляем последний или предпоследний день месяца
					lastDay := daysInMonth(nextYear, nextMonth)
					if d == -1 {
						actualDay = lastDay
					} else {
						actualDay = lastDay - 1
					}
				}
				if actualDay < minDay && (nextYear > year || nextMonth > currentMonth || actualDay > currentDay) {
					minDay = actualDay
				}
			}
			if minDay == 32 {
				// Если не нашли в этом месяце, берем минимальный в следующем
				nextMonth++
				if nextMonth > 12 {
					nextMonth = 1
					nextYear++
				}
				minDay = 32
				for _, d := range monthDays {
					actualDay := d
					if d == -1 || d == -2 {
						lastDay := daysInMonth(nextYear, nextMonth)
						if d == -1 {
							actualDay = lastDay
						} else {
							actualDay = lastDay - 1
						}
					}
					if actualDay < minDay {
						minDay = actualDay
					}
				}
			}
			nextDay = minDay
		}

		// Проверяем, что день существует в месяце
		lastDay := daysInMonth(nextYear, nextMonth)
		if nextDay > lastDay {
			nextDay = lastDay
		}

		date = time.Date(nextYear, time.Month(nextMonth), nextDay, 0, 0, 0, 0, time.UTC)
	}

	return date.Format("20060102"), nil
}

func daysInMonth(year, month int) int {
	return time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()
}