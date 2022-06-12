package store

import "gorm.io/gorm"

type UserRelation struct {
	gorm.Model
	UserId   uint `gorm:"column:user_id" json:"user_id"`
	ToUserId uint `gorm:"column:to_user_id" json:"to_user_id"`
	Status   int  `gorm:"column:status" json:"status"`
}

func AddUserRelation(userId uint, toUserId uint) error {
	var relationCount int64
	Db.Model(&UserRelation{}).Where("user_id = ? and to_user_id = ?", userId, toUserId).
		Count(&relationCount)
	if relationCount > 0 {
		Db.Model(&UserRelation{}).Where("user_id = ? and to_user_id = ?", userId, toUserId).
			UpdateColumn("status", 1)
	} else {
		data := UserRelation{
			UserId:   userId,
			ToUserId: toUserId,
			Status:   1,
		}
		err := Db.Model(&UserRelation{}).Create(&data).Error
		if err != nil {
			return err
		}
	}
	//修改用户和被关注用户的响应记录
	Db.Model(&User{}).Where("id = ?", userId).
		UpdateColumn("follow_count", gorm.Expr("follow_count + ?", 1))
	Db.Model(&User{}).Where("id = ?", toUserId).
		UpdateColumn("follower_count", gorm.Expr("follower_count + ?", 1))
	return nil
}

func DeleteUserRelation(userId uint, toUserId uint) error {
	err := Db.Model(&UserRelation{}).Where("user_id = ? and to_user_id = ?", userId, toUserId).
		UpdateColumn("status", 0).Error
	if err != nil {
		return err
	}

	//修改用户和被关注用户的响应记录
	Db.Model(&User{}).Where("id = ?", userId).
		UpdateColumn("follow_count", gorm.Expr("follow_count - ?", 1))
	Db.Model(&User{}).Where("id = ?", toUserId).
		UpdateColumn("follower_count", gorm.Expr("follower_count - ?", 1))
	return nil
}

func GetFollowList(userId uint) (followList []User) {
	var toUserIds []uint
	Db.Model(&UserRelation{}).Select("to_user_id").Where("user_id = ? and status = 1", userId).Find(&toUserIds)
	if len(toUserIds) == 0 {
		return
	}
	Db.Model(&User{}).Find(&followList, toUserIds)
	return
}

func GetFollowerList(userId uint) (followerList []User) {
	var userIds []uint
	Db.Model(&UserRelation{}).Select("user_id").Where("to_user_id = ?", userId).Find(&userIds)
	if len(userIds) == 0 {
		return
	}
	Db.Model(&User{}).Find(&followerList, userIds)
	return
}
