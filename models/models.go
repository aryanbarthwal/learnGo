package models

type Model struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string
	Email    string
	Password string
}