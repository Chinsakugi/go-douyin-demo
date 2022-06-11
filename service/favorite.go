package service

import (
	"github.com/gin-gonic/gin"
	jwtHelper "go-douyin-demo/middleware/jwt"
	"go-douyin-demo/store"
	"net/http"
	"strconv"
)

// FavoriteAction
// @Tags         扩展接口-1
// @Summary      赞操作
// @Param        token query string true "用户鉴权token"
// @Param        video_id query int true "视频id"
// @Param        action_type query int true "1-点赞，2-取消点赞"
// @Success      200  {string}  json "{"status_code":"200","status_msg":""}"
// @Router       /douyin/favorite/action/ [post]
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	videoId := c.Query("video_id")
	actionType := c.Query("action_type")
	if token == "" || videoId == "" || actionType == "" {
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

	//transfer videoId from string to uint
	videoIdentity, _ := strconv.Atoi(videoId)

	//check actionType
	if actionType != "1" && actionType != "2" {
		c.JSON(http.StatusOK, gin.H{
			"status_code": -1,
			"status_msg":  "actionType is not valid, actionType:" + actionType,
		})
		return
	}

	//进行点赞或取消点赞处理
	var statusMsg string
	if actionType == "1" {
		err := store.AddFavoriteVideo(userID, uint(videoIdentity))
		statusMsg = "点赞成功"
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status_code": -1,
				"status_msg":  "add favorite video error, err: " + err.Error(),
			})
			return
		}
	} else {
		err := store.DeleteFavoriteVideo(userID, uint(videoIdentity))
		statusMsg = "取消点赞成功"
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status_code": -1,
				"status_msg":  "delete favorite video error, err: " + err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  statusMsg,
	})

}

// FavoriteList
// @Tags         扩展接口-1
// @Summary      点赞列表
// @Param        token query string true "用户鉴权token"
// @Param        user_id query int true "用户id"
// @Success      200  {string}  json "{"status_code":"200","status_msg":""}"
// @Router       /douyin/favorite/list/ [get]
func FavoriteList(c *gin.Context) {
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
	userClaims, err := jwtHelper.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status_code": -1,
			"status_msg":  "token parse error:" + err.Error(),
		})
		return
	}
	userID := userClaims.UserID
	userIdentity, _ := strconv.Atoi(userId)
	if userID != uint(userIdentity) {
		c.JSON(http.StatusOK, gin.H{
			"status_code": -1,
			"status_msg":  "token error, please login again",
		})
		return
	}

	videoList := store.GetFavoriteVideoList(userID)

	var videoListRes []map[string]interface{}
	for _, video := range videoList {
		videoRes := make(map[string]interface{})
		videoRes["id"] = video.ID
		videoRes["play_url"] = video.PlayUrl
		videoRes["cover_url"] = video.CoverUrl
		videoRes["favorite_count"] = video.FavoriteCount
		videoRes["comment_count"] = video.CommentCount
		videoRes["is_favorite"] = true
		videoRes["title"] = video.Title

		author := make(map[string]interface{})
		author["id"] = video.Author.ID
		author["name"] = video.Author.Username
		author["follow_count"] = video.Author.FollowCount
		author["follower_count"] = video.Author.FollowerCount
		author["is_follow"] = video.Author.IsFollow

		videoRes["author"] = author
		videoListRes = append(videoListRes, videoRes)
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "查询成功",
		"video_list":  videoListRes,
	})

}
