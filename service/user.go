package service

import (
	"github.com/gin-gonic/gin"
	jwtHelper "go-douyin-demo/middleware/jwt"
	"go-douyin-demo/store"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type UserRegisterResponse struct {
	Response
	UserId uint   `json:"user_id,omitempty"`
	Token  string `json:"token,omitempty"`
}

type UserLoginResponse struct {
	Response
	UserId uint   `json:"user_id,omitempty"`
	Token  string `json:"token,omitempty"`
}

type UserResponse struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

type UserInfoResponse struct {
	Response
	User UserResponse `json:"user,omitempty"`
}

// Register
// @Tags         基础接口
// @Summary      用户注册
// @Param        username query string true "注册用户名，最长32个字符"
// @Param        password query string true "密码，最长32个字符"
// @Success      200  {string}  json "{"status_code":"200","status_msg":"", "user_id":"", "token":""}"
// @Router       /douyin/user/register/ [post]
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	if username == "" || password == "" {
		c.JSON(http.StatusOK, UserRegisterResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "缺少username或password",
			},
		})
		return
	}

	user := store.CreateUser(username, password)
	if user.ID == 0 {
		c.JSON(http.StatusOK, UserRegisterResponse{
			Response: Response{StatusCode: -1, StatusMsg: "注册失败"},
		})
		return
	}

	//生成token
	token, err := jwtHelper.GenToken(user.ID, username)
	if err != nil {
		c.JSON(http.StatusOK, UserRegisterResponse{
			Response: Response{StatusCode: -1, StatusMsg: "生成token错误：" + err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, UserRegisterResponse{
		Response: Response{StatusCode: 0, StatusMsg: "注册成功"},
		UserId:   user.ID,
		Token:    token,
	})
}

// Login
// @Tags         基础接口
// @Summary      用户登录
// @Param        username query string true "登录用户名"
// @Param        password query string true "登录密码"
// @Success      200  {string}  json "{"status_code":"200","status_msg":"", "user_id":"", "token":""}"
// @Router       /douyin/user/login/ [post]
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	if username == "" || password == "" {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "缺少username或password",
			},
		})
		return
	}
	loginUser := store.LoginCheck(username, password)
	if loginUser.ID == 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: -1, StatusMsg: "username 或 password 错误"},
		})
		return
	}
	//生成token
	token, err := jwtHelper.GenToken(loginUser.ID, username)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: -1, StatusMsg: "生成token错误：" + err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0, StatusMsg: "登录成功"},
		UserId:   loginUser.ID,
		Token:    token,
	})
}

// UserInfo
// @Tags         基础接口
// @Summary      用户信息
// @Param        user_id query string true "用户id"
// @Param        token query string true "用户鉴权token"
// @Success      200  {string}  json "{"status_code":"200","status_msg":"", "user_id":"", "token":""}"
// @Router       /douyin/user/ [get]
func UserInfo(c *gin.Context) {
	userId := c.Query("user_id")
	token := c.Query("token")
	if userId == "" || token == "" {
		c.JSON(http.StatusOK, UserInfoResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "缺少user_id或token",
			},
		})
		return
	}

	id, _ := strconv.Atoi(userId)
	user, err := store.GetUser(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, UserInfoResponse{
				Response: Response{StatusCode: -1, StatusMsg: "user id not exists:" + err.Error()},
			})
		} else {
			c.JSON(http.StatusOK, UserInfoResponse{
				Response: Response{StatusCode: -1, StatusMsg: err.Error()},
			})
		}
		return
	}
	userClaims, err := jwtHelper.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, UserInfoResponse{
			Response: Response{StatusCode: -1, StatusMsg: "parse token error :" + err.Error()},
		})
		return
	}
	if userClaims.UserID != uint(id) {
		c.JSON(http.StatusOK, UserInfoResponse{
			Response: Response{StatusCode: -1, StatusMsg: "token error"},
		})
		return
	}
	c.JSON(http.StatusOK, UserInfoResponse{
		Response: Response{StatusCode: 200, StatusMsg: "查询成功"},
		User: UserResponse{
			ID:            user.ID,
			Name:          user.Username,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow:      user.IsFollow,
		},
	})
}
