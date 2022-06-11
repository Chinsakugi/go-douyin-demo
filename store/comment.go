package store

import (
	"gorm.io/gorm"
)

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

func AddComment(videoId, userId uint, commentText string) (Comment, error) {
	data := Comment{
		UserID:  userId,
		VideoID: videoId,
		Content: commentText,
	}
	err := Db.Model(&Comment{}).Create(&data).Error
	if err != nil {
		return data, err
	}
	//修改video 的 comment_count
	Db.Model(&Video{}).Where("id = ?", videoId).
		UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1))
	return data, nil
}

func DeleteComment(commentId uint, videoId uint) error {
	err := Db.Where("id = ?", commentId).Delete(&Comment{}).Error
	if err != nil {
		return err
	}
	//修改video 的 comment_count
	Db.Model(&Video{}).Where("id = ?", videoId).
		UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1))
	return nil
}

func GetCommentList(videoId uint) (commentList []Comment) {
	Db.Model(&Comment{}).Preload("User").Where("video_id = ?", videoId).
		Find(&commentList)
	return
}
