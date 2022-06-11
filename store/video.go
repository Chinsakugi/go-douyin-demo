package store

import (
	"gorm.io/gorm"
	"time"
)

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

func GetVideoList(userId uint) (err error, videoList []Video) {
	err = Db.Model(&Video{}).Preload("Author").Where("user_id = ?", userId).Find(&videoList).Error
	return
}

func GetVideoFeed(latestTime int) (err error, videoList []Video) {
	err = Db.Model(&Video{}).Preload("Author").
		Where("created_at < ?", time.UnixMilli(int64(latestTime)).Format("2006-01-02 15:04:05")).
		Order("created_at desc").Limit(30).Find(&videoList).Error
	return
}
