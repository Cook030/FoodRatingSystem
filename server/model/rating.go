package model

import "time"

type Rating struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         string    `gorm:"type:varchar(100)" json:"user_id"`
	UserName       string    `gorm:"type:varchar(255)" json:"user_name"`
	RestaurantID   uint      `gorm:"index" json:"restaurant_id"`
	RestaurantName string    `gorm:"type:varchar(255)" json:"restaurant_name"`
	Stars          float64   `gorm:"type:decimal(2,1)" json:"stars"`
	Comment        string    `gorm:"type:text" json:"comment"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
}
