package models

import (
	"fmt"
)

const (
	TitleError             = "не указан заголовок задачи"
	DateError              = "не корректная дата"
	RepeatError            = "пустое правило повторения"
	RepeatRuleError        = "неподдерживаемый формат правила повторения"
	RuleDError             = "неверный формат для правила 'd'"
	RuleDDaysError         = "неверное количество дней"
	RuleDDateIntervalError = "недопустимый интервал в днях"
	RuleYError             = "неверный формат для правила 'y'"
	RuleWError             = "не верный формат правила w"
	RuleWWeekDayError      = "не верный день недели в парвиле w"
	RuleMError             = "не верный формат правила m"
	RuleMDayError          = "не верный день в парвиле m"
	RequestBodyError       = "ошибка чтения тела запроса"
	ConvertationError      = "ошибка конвертации"
	SaveTaskError          = "ошибка при сохранении задачи"
	UnmarshalError         = "ошибка анмаршаллинга"
	MarshalError           = "ошибка маршаллинга"
	RepeatParamError       = "обязателен параметр 'repeat'"
	GetTasksError          = "ошибка получения заявок"
	GetTaskError           = "ошибка получения заявки"
	IdFormatError          = "не верный формат id задачи"
	DeleteTaskError        = "ошибка удаления заявки"
)

type ResponseError struct {
	MyError string `json:"error"`
}

func (re ResponseError) Error() string {
	return fmt.Sprintf("%s", re.MyError)
}

var RespError map[string]string
