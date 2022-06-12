package service

import (
	"github.com/gin-gonic/gin"
	jwtHelper "go-douyin-demo/middleware/jwt"
	"go-douyin-demo/store"
	"net/http"
	"strconv"
)

// RelationAction
// @Tags         扩展接口-2
// @Summary      关注操作
// @Param        token query string true "用户鉴权token"
// @Param        to_user_id query int true "对方用户id"
// @Param        action_type query int true "1-关注，2-取消关注"
// @Success      200  {string}  json "{"status_code":"200","status_msg":""}"
// @Router       /douyin/relation/action/ [post]
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")
	actionType := c.Query("action_type")
	if token == "" || toUserId == "" || actionType == "" {
		c.JSON(http.StatusOK, gin.H{
			"status_code": -1,
			"status_msg":  "缺少参数",
		})
		return
	}
	//check token
	userClaims, err := jwtHelper.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status_code": -1,
			"status_msg":  "token parse error:" + err.Error(),
		})
		return
	}
	userID := userClaims.UserID

	//transfer toUserId from string to uint
	toUserIdentity, _ := strconv.Atoi(toUserId)

	//check actionType
	if actionType != "1" && actionType != "2" {
		c.JSON(http.StatusOK, gin.H{
			"status_code": -1,
			"status_msg":  "actionType is not valid, actionType:" + actionType,
		})
		return
	}

	var statusMsg string
	if actionType == "1" {
		err := store.AddUserRelation(userID, uint(toUserIdentity))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status_code": -1,
				"status_msg":  "add user relation error, err: " + err.Error(),
			})
			return
		}
		statusMsg = "关注成功"
	} else {
		err := store.DeleteUserRelation(userID, uint(toUserIdentity))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status_code": -1,
				"status_msg":  "delete user relation error, err: " + err.Error(),
			})
			return
		}
		statusMsg = "取消关注成功"
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  statusMsg,
	})
}

// FollowList
// @Tags         扩展接口-2
// @Summary      关注列表
// @Param        token query string true "用户鉴权token"
// @Param        user_id query int true "用户id"
// @Success      200  {string}  json "{"status_code":"200","status_msg":""}"
// @Router       /douyin/relation/follow/list/ [get]
func FollowList(c *gin.Context) {
	token := c.Query("token")
	userId := c.Query("user_id")
	if token == "" || userId == "" {
		c.JSON(http.StatusOK, gin.H{
			"status_code": -1,
			"status_msg":  "缺少参数",
		})
		return
	}
	//check token
	_, err := jwtHelper.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status_code": -1,
			"status_msg":  "token parse error:" + err.Error(),
		})
		return
	}

	userIdentity, _ := strconv.Atoi(userId)
	followList := store.GetFollowList(uint(userIdentity))
	var followListRes []map[string]interface{}
	for _, user := range followList {
		userRes := make(map[string]interface{})
		userRes["id"] = user.ID
		userRes["name"] = user.Username
		userRes["follow_count"] = user.FollowCount
		userRes["follower_count"] = user.FollowerCount
		userRes["is_follow"] = true
		followListRes = append(followListRes, userRes)
	}
	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "查询成功",
		"user_list":   followListRes,
	})
}

// FollowerList
// @Tags         扩展接口-2
// @Summary      粉丝列表
// @Param        token query string true "用户鉴权token"
// @Param        user_id query int true "用户id"
// @Success      200  {string}  json "{"status_code":"200","status_msg":""}"
// @Router       /douyin/relation/follower/list/ [get]
func FollowerList(c *gin.Context) {
	token := c.Query("token")
	userId := c.Query("user_id")
	if token == "" || userId == "" {
		c.JSON(http.StatusOK, gin.H{
			"status_code": -1,
			"status_msg":  "缺少参数",
		})
		return
	}
	//check token
	_, err := jwtHelper.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status_code": -1,
			"status_msg":  "token parse error:" + err.Error(),
		})
		return
	}

	userIdentity, _ := strconv.Atoi(userId)
	followerList := store.GetFollowerList(uint(userIdentity))
	var followerListRes []map[string]interface{}
	for _, user := range followerList {
		userRes := make(map[string]interface{})
		userRes["id"] = user.ID
		userRes["name"] = user.Username
		userRes["follow_count"] = user.FollowCount
		userRes["follower_count"] = user.FollowerCount
		userRes["is_follow"] = false
		followerListRes = append(followerListRes, userRes)
	}
	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "查询成功",
		"user_list":   followerListRes,
	})
}
