package utils

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ilexsor/internal/models"
	"gorm.io/gorm"
)

const (
	DateFormat = "20060102"
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

	startDate, err := time.Parse(DateFormat, dstart)
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
			return "", errors.New("не верный формат правила w")
		}
		daysStr := strings.Split(parts[1], ",")
		weekdays := make([]int, 0, len(daysStr))
		for _, dayStr := range daysStr {
			day, err := strconv.Atoi(dayStr)
			if err != nil || day < 1 || day > 7 {
				return "", errors.New("не верный день недели в парвиле w")
			}
			weekdays = append(weekdays, day)
		}
		return findNextWeekday(startDate, now, weekdays).Format(DateFormat), nil
	case "m":
		if len(parts) != 2 {
			return "", errors.New("не верный формат правила m")
		}
		daysStr := strings.Split(parts[1], ",")
		monthDays := make([]int, 0, len(daysStr))
		for _, dayStr := range daysStr {
			day, err := strconv.Atoi(dayStr)
			if err != nil {
				return "", errors.New("не верный день в парвиле m")
			}
			if day < -31 || day == 0 || day > 31 {
				return "", errors.New("не верный день в парвиле m")
			}
			monthDays = append(monthDays, day)
		}
		nextDate := findNextMonthDay(startDate, now, monthDays)
		return nextDate.Format(DateFormat), nil
	default:
		return "", errors.New("неподдерживаемый формат правила повторения")
	}
}

// Case "d"  задача переносится на указанное число дней
func nextDateDaily(now, startDate time.Time, days int) string {

	if startDate.Compare(now) == 0 {
		return startDate.AddDate(0, 0, days).Format(DateFormat)
	}

	if startDate.After(now) {
		return startDate.AddDate(0, 0, days).Format(DateFormat)
	}
	for !startDate.After(now) {
		startDate = startDate.AddDate(0, 0, days)
	}
	return startDate.Format(DateFormat)
}

// Case "y" задача выполняется ежегодно
func nextDateYearly(now, startDate time.Time) string {

	if startDate.Format("20060102") == now.Format(DateFormat) {
		return startDate.AddDate(1, 0, 0).Format(DateFormat)
	}

	if startDate.After(now) {
		return startDate.AddDate(1, 0, 0).Format(DateFormat)
	}

	for !startDate.After(now) {
		startDate = startDate.AddDate(1, 0, 0)
	}
	return startDate.Format(DateFormat)
}

// Case "w" задача выполняется в указанные дни месяца
func findNextWeekday(startDate, now time.Time, weekdays []int) time.Time {
	current := startDate
	for current.Before(now) || current.Equal(now) {
		// Check if current day is one of the target weekdays and after startDate
		if current.After(startDate) {
			currentWeekday := int(current.Weekday())
			if currentWeekday == 0 {
				currentWeekday = 7 // Sunday is 7 in our system
			}
			for _, wd := range weekdays {
				if currentWeekday == wd {
					return current
				}
			}
		}
		current = current.AddDate(0, 0, 1)
	}

	// If we passed now, find the next weekday
	for {
		currentWeekday := int(current.Weekday())
		if currentWeekday == 0 {
			currentWeekday = 7
		}
		for _, wd := range weekdays {
			if currentWeekday == wd {
				return current
			}
		}
		current = current.AddDate(0, 0, 1)
	}
}

// Case "m" задача назначается на указанные дни месяца
func findNextMonthDay(startDate, now time.Time, monthDays []int) time.Time {
	current := startDate
	for current.Before(now) || current.Equal(now) {
		// Check if current day is one of the target days and after startDate
		if current.After(startDate) {
			day := current.Day()
			for _, md := range monthDays {
				if md > 0 && day == md {
					return current
				} else if md < 0 {
					// Handle negative days (last days of month)
					lastDay := daysInMonth(current.Year(), current.Month())
					if day == lastDay+md+1 {
						return current
					}
				}
			}
		}
		// Move to the next candidate day
		current = nextMonthDayCandidate(current, monthDays)
	}

	// If we passed now, find the next valid day
	for {
		day := current.Day()
		for _, md := range monthDays {
			if md > 0 && day == md {
				return current
			} else if md < 0 {
				lastDay := daysInMonth(current.Year(), current.Month())
				if day == lastDay+md+1 {
					return current
				}
			}
		}
		current = nextMonthDayCandidate(current, monthDays)
	}
}

func nextMonthDayCandidate(date time.Time, monthDays []int) time.Time {
	// Find the smallest day in monthDays that is greater than current day
	currentDay := date.Day()
	year, month, _ := date.Date()

	var found bool
	var minPositive = 32
	var maxNegative = -32

	for _, md := range monthDays {
		if md > 0 {
			if md > currentDay && md < minPositive {
				minPositive = md
				found = true
			}
		} else {
			if md > maxNegative {
				maxNegative = md
			}
		}
	}

	if found {
		daysInMonth := daysInMonth(year, month)
		if minPositive > daysInMonth {
			// Move to next month and find the first suitable day
			return time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
		}
		return time.Date(year, month, minPositive, 0, 0, 0, 0, time.UTC)
	}

	// If no positive days found, use the smallest positive or negative day in next month
	if len(monthDays) > 0 {
		nextMonth := month + 1
		nextYear := year
		if nextMonth > 12 {
			nextMonth = 1
			nextYear += 1
		}

		// Try positive days first
		for _, md := range monthDays {
			if md > 0 {
				daysInNextMonth := daysInMonth(nextYear, nextMonth)
				if md <= daysInNextMonth {
					return time.Date(nextYear, nextMonth, md, 0, 0, 0, 0, time.UTC)
				}
			}
		}

		// If no positive days work, use negative days
		for _, md := range monthDays {
			if md < 0 {
				daysInNextMonth := daysInMonth(nextYear, nextMonth)
				day := daysInNextMonth + md + 1
				if day > 0 {
					return time.Date(nextYear, nextMonth, day, 0, 0, 0, 0, time.UTC)
				}
			}
		}
	}

	// Fallback (shouldn't happen with proper validation)
	return date.AddDate(0, 1, 0)
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
		currentMonth := month
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
				if nextMonth == time.Month(m) {
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

	return date.Format(DateFormat), nil
}

func daysInMonth(year int, month time.Month) int {
	return time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
