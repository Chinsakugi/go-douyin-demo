package service

import (
	"github.com/gin-gonic/gin"
	"go-douyin-demo/store"
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
}
