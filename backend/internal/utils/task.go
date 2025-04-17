package utils

import (
	"strconv"
	"strings"
	"time"

	"github.com/ilexsor/internal/models"
)

// ConvertTaskToSchedule функция для конвертации полученной задачи через АПИ в структуру для БД
func ConvertTaskToSchedule(task *models.Task, sched *models.Scheduler) (*models.Scheduler, error) {

	sched.ID, _ = strconv.Atoi(task.ID)
	sched.Date = task.Date
	sched.Title = task.Title
	sched.Comment = task.Comment
	sched.Repeat = task.Repeat

	return sched, nil
}

// CheckTask функция для проверки корректности заполнения задачи
func CheckTask(task *models.Task) (*models.Task, error) {

	date := task.Date
	title := task.Title
	repeat := task.Repeat

	now := time.Now().Format(DateFormat)

	if title == "" || len(title) == 0 {

		return nil, models.ResponseError{
			MyError: models.TitleError,
		}

	}

	if date == "" || len(date) == 0 {
		task.Date = now
		return task, nil

	}

	//dateTime дата из структуры Task переведенная в формат time.Time для проверки на коррктность
	_, err := time.Parse(DateFormat, date)
	if err != nil {
		return nil, models.ResponseError{
			MyError: models.DateError,
		}
	}

	if date < now {
		if repeat == "" || len(repeat) == 0 {
			date = now
			task.Date = date
		} else {
			date, err := NextDate(time.Now(), date, repeat)
			if err != nil {
				return nil, err
			}
			task.Date = date
		}
	}

	if repeat != "" {
		parts := strings.Fields(repeat)
		if len(parts) == 0 {
			return nil, models.ResponseError{
				MyError: models.RepeatRuleError,
			}
		}
		rule := parts[0]
		switch rule {
		case "d":
			if len(parts) != 2 {
				return nil, models.ResponseError{
					MyError: models.RuleDError,
				}
			}
			days, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, models.ResponseError{
					MyError: models.RuleDDateIntervalError,
				}
			}
			if days <= 0 || days > 400 {
				return nil, models.ResponseError{
					MyError: models.RuleDDateIntervalError,
				}
			}
		case "y":
			if len(parts) != 1 {
				return nil, models.ResponseError{
					MyError: models.RuleYError,
				}
			}
		case "w":
			if len(parts) != 2 {
				return nil, models.ResponseError{
					MyError: models.RuleWError,
				}
			}

		case "m":
			if len(parts) != 2 {
				return nil, models.ResponseError{
					MyError: models.RuleMError,
				}
			}
		default:
			return nil, models.ResponseError{
				MyError: models.RepeatRuleError,
			}
		}
	}

	return task, nil
}
