package store

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	UserID        uint       `gorm:"column:user_id" json:"user_id"`
	Author        User       `gorm:"foreignKey:id;references:user_id" json:"author"`
	PlayUrl       string     `gorm:"column:play_url" json:"play_url"`
	CoverUrl      string     `gorm:"column:cover_url" json:"cover_url"`
	FavoriteCount int64      `gorm:"column:favorite_count" json:"favorite_count"`
	CommentCount  int64      `gorm:"column:comment_count" json:"comment_count"`
	IsFavorite    bool       `gorm:"column:is_favorite" json:"is_favorite"`
	Title         string     `gorm:"column:title" json:"title"`
	CommentList   []*Comment `gorm:"foreignKey:video_id;references:id" json:"comment_list"`
}

func (v Video) TableName() string {
	return "video"
}
