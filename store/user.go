package store

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username      string `gorm:"column:username" json:"username"`
	Password      string `gorm:"column:password" json:"password"`
	FollowCount   int64  `gorm:"column:follow_count" json:"follow_count"`
	FollowerCount int64  `gorm:"column:follower_count" json:"follower_count"`
	IsFollow      bool   `gorm:"column:is_follow" json:"is_follow"`
}

func (u User) TableName() string {
	return "user"
}

func CreateUser(username, password string) *User {
	user := User{
		Username: username,
		Password: password,
		IsFollow: false,
	}
	Db.Model(&User{}).Create(&user)
	return &user
}

func LoginCheck(username, password string) User {
	var user User
	Db.Model(&User{}).Where("username = ? and password = ?", username, password).First(&user)
	return user
}

func GetUser(userId uint) (user User, err error) {
	err = Db.Model(&User{}).First(&user, userId).Error
	return user, err
}
