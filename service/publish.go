package service

import (
	"github.com/gin-gonic/gin"
	jwtHelper "go-douyin-demo/middleware/jwt"
	"go-douyin-demo/middleware/util"
	"go-douyin-demo/store"
	"net/http"
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
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: -1, StatusMsg: "save file error" + err.Error()})
		return
	}

	//提取封面并保存
	err, coverPath := util.GetSnapshot(filePath, "D:/go/go-douyin-demo/public/video_cover/", 5)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: -1, StatusMsg: "get video cover error: " + err.Error()})
		return
	}

	//保存视频信息
	video := store.Video{
		UserID:   userId,
		PlayUrl:  filePath,
		CoverUrl: coverPath,
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

}
