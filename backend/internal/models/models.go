package models

import "time"

type Scheduler struct {
	ID      int       `gorm:"primaryKey;autoIncrement"`
	Date    time.Time `gorm:"index:idx_date"`
	Title   string    `gorm:"type:varchar(255);not null"`
	Comment string    `gorm:"type:text"`
	Repeat  string    `gorm:"type:varchar(100)"`
}
