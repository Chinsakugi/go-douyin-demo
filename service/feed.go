package service

import (
	"github.com/gin-gonic/gin"
	jwtHelper "go-douyin-demo/middleware/jwt"
	"go-douyin-demo/store"
	"net/http"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []store.Video `json:"video_list,omitempty"`
	NextTime  int64         `json:"next_time,omitempty"`
}

// Feed
// @Tags         基础接口
// @Summary      视频流接口
// @Param        latest_time query int false "可选参数，限制返回视频的最新投稿时间戳，精确到秒，不填表示当前时间"
// @Param        token query string false "用户登录状态下设置"
// @Success      200  {string}  json "{"status_code":"200","status_msg":"", "user_id":"", "token":""}"
// @Router       /douyin/feed [get]
func Feed(c *gin.Context) {
	//latest_time := c.Query("latest_time")
	token := c.Query("token")

	if token != "" {
		_, err := jwtHelper.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status_code": -1,
				"status_msg":  "token error: " + err.Error(),
				"next_time":   nil,
				"video_list":  nil,
			})
			return
		}
	}

	var timeCondition int = int(time.Now().UnixMilli())

	//if latest_time == "" {
	//	timeCondition = int(time.Now().UnixMilli())
	//} else {
	//	timeCondition, _ = strconv.Atoi(latest_time)
	//}
	err, videoList := store.GetVideoFeed(timeCondition)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status_code": -1,
			"status_msg":  "get video feed error: " + err.Error(),
			"next_time":   nil,
			"video_list":  nil,
		})
		return
	}
	if len(videoList) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status_code": -1,
			"status_msg":  "video number is empty" + err.Error(),
			"next_time":   nil,
			"video_list":  nil,
		})
		return
	}

	var videoListRes []map[string]interface{}

	for _, video := range videoList {
		videoRes := make(map[string]interface{})
		videoRes["id"] = video.ID
		videoRes["play_url"] = video.PlayUrl
		videoRes["cover_url"] = video.CoverUrl
		videoRes["favorite_count"] = video.FavoriteCount
		videoRes["comment_count"] = video.CommentCount
		videoRes["is_favorite"] = video.IsFavorite
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
		"next_time":   videoList[0].CreatedAt.UnixMilli(),
		"video_list":  videoListRes,
	})
}
