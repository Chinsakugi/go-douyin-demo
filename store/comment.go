package store

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	UserID  uint   `gorm:"column:user_id" json:"user_id"`
	User    User   `gorm:"foreignKey:id;references:user_id" json:"user"`
	VideoID uint   `gorm:"column:video_id" json:"video_id"`
	Content string `gorm:"column:content" json:"content"`
}

func (c Comment) TableName() string {
	return "comment"
}
