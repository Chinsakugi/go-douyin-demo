package service

import (
	"github.com/gin-gonic/gin"
	jwtHelper "go-douyin-demo/middleware/jwt"
	"go-douyin-demo/middleware/util"
	"go-douyin-demo/store"
	"net/http"
	"strconv"
)

type VideoListResponse struct {
	Response
	VideoList []store.Video `json:"video_list"`
}

// Publish
// @Tags         基础接口
// @Summary      投稿接口
// @Param        data formData file true "视频数据"
// @Param        token formData string true "用户鉴权token"
// @Param        title formData string true "视频标题"
// @Success      200  {string}  json "{"status_code":"200","status_msg":"", "user_id":"", "token":""}"
// @Router       /douyin/publish/action/ [post]
func Publish(c *gin.Context) {
	file, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: -1,
			StatusMsg:  "file upload error: " + err.Error(),
		})
		return
	}
	token := c.PostForm("token")
	title := c.PostForm("title")
	if token == "" || title == "" {
		c.JSON(http.StatusOK, Response{
			StatusCode: -1,
			StatusMsg:  "缺少token或title",
		})
		return
	}
	//验证token
	userClaims, err := jwtHelper.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: -1, StatusMsg: "parse token error: " + err.Error()})
		return
	}
	userId := userClaims.UserID

	//保存文件
	filePath := "D:/go/go-douyin-demo/public/video_data/" + file.Filename
	err = c.SaveUploadedFile(file, filePath)
	playUrl := "http://192.168.100.19:8080/static/video_data/" + file.Filename
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: -1, StatusMsg: "save file error" + err.Error()})
		return
	}

	//提取封面并保存
	err, coverName := util.GetSnapshot(filePath, "D:/go/go-douyin-demo/public/video_cover/", 3)
	coverUrl := "http://192.168.100.19:8080/static/video_cover/" + coverName
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: -1, StatusMsg: "get video cover error: " + err.Error()})
		return
	}

	//保存视频信息
	video := store.Video{
		UserID:   userId,
		PlayUrl:  playUrl,
		CoverUrl: coverUrl,
		Title:    title,
	}
	err = store.Db.Model(&store.Video{}).Create(&video).Error
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: -1, StatusMsg: "create video info error" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "发布成功",
	})
}

// PublishList
// @Tags         基础接口
// @Summary      发布列表
// @Param        user_id query string true "用户id"
// @Param        token query string true "用户鉴权token"
// @Success      200  {string}  json "{"status_code":"200","status_msg":"", "user_id":"", "token":""}"
// @Router       /douyin/publish/list/ [get]
func PublishList(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Query("user_id"))
	token := c.Query("token")
	if userId == 0 || token == "" {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "缺少user_id或token",
			},
			VideoList: nil,
		})
		return
	}

	//验证token
	userClaims, err := jwtHelper.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response:  Response{StatusCode: -1, StatusMsg: "parse token error :" + err.Error()},
			VideoList: nil,
		})
		return
	}
	if userClaims == nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response:  Response{StatusCode: -1, StatusMsg: "token error"},
			VideoList: nil,
		})
		return
	}

	//获取用户发布视频列表
	err, videoList := store.GetVideoList(uint(userId))
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response:  Response{StatusCode: -1, StatusMsg: "token error"},
			VideoList: nil,
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
		var count int64
		store.Db.Model(&store.FavoriteVideo{}).Where("user_id = ? and video_id = ?", userClaims.UserID, video.ID).Count(&count)
		if count > 0 {
			videoRes["is_favorite"] = true
		} else {
			videoRes["is_favorite"] = false
		}
		videoRes["title"] = video.Title

		author := make(map[string]interface{})
		author["id"] = video.Author.ID
		author["name"] = video.Author.Username
		author["follow_count"] = video.Author.FollowCount
		author["follower_count"] = video.Author.FollowerCount
		var isFollow int64
		store.Db.Model(&store.UserRelation{}).Where("user_id = ? and to_user_id = ?", userClaims.UserID, video.Author.ID).Count(&isFollow)
		if isFollow > 0 {
			author["is_follow"] = true
		} else {
			author["is_follow"] = false
		}
		videoRes["author"] = author
		videoListRes = append(videoListRes, videoRes)
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "查询成功",
		"video_list":  videoListRes,
	})
}
