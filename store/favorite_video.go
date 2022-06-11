package store

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type FavoriteVideo struct {
	gorm.Model
	UserId  uint `gorm:"column:user_id" json:"user_id"`
	VideoId uint `gorm:"column:video_id" json:"video_id"`
}

func AddFavoriteVideo(userId uint, videoId uint) error {
	data := FavoriteVideo{
		UserId:  userId,
		VideoId: videoId,
	}

	var count int64
	Db.Model(&FavoriteVideo{}).Where("user_id = ? and video_id = ?", userId, videoId).Count(&count)
	if count > 0 {
		return errors.New(fmt.Sprintf("点赞记录已存在, user_id = %v, video_id = %v", userId, videoId))
	}

	err := Db.Model(&FavoriteVideo{}).Create(&data).Error
	if err != nil {
		return err
	}
	//修改video 的 favorite_count
	Db.Model(&Video{}).Where("id = ?", videoId).
		UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1))
	return nil
}

func DeleteFavoriteVideo(userId uint, videoId uint) error {
	err := Db.Where("user_id = ? and video_id = ?", userId, videoId).Delete(&FavoriteVideo{}).Error
	if err != nil {
		return err
	}
	//修改video 的 favorite_count
	Db.Model(&Video{}).Where("id = ?", videoId).
		UpdateColumn("favorite_count", gorm.Expr("favorite_count - ?", 1))
	return nil
}

func GetFavoriteVideoList(userId uint) (videoList []Video) {
	Db.Model(Video{}).Preload("Author").Where("id in ( select video_id from favorite_video where user_id = ? and deleted_at is null)", userId).
		Find(&videoList)
	return
}
