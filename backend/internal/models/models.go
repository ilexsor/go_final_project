package models

const (
	TaskLimit = 50
)

// Структура для миграции и корректного добавления таски в БД
// ID int в остальных случаях string
type Scheduler struct {
	ID      int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Date    string `gorm:"index:idx_date" json:"date"`
	Title   string `gorm:"type:varchar(255);not null" json:"title"`
	Comment string `gorm:"type:text" json:"comment"`
	Repeat  string `gorm:"type:varchar(100)" json:"repeat"`
}

func (Scheduler) TableName() string {
	return "scheduler"
}

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func (Task) TableName() string {
	return "scheduler"
}

// Структура для отета с номером таски
type TaskResponse struct {
	ID string `json:"id"`
}

// Структура для ответа списком тасок
type TasksResponse struct {
	Tasks []Task `json:"tasks"`
}
