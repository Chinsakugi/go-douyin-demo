package service

import (
	"github.com/gin-gonic/gin"
	jwtHelper "go-douyin-demo/middleware/jwt"
	"go-douyin-demo/store"
	"net/http"
	"strconv"
)

// CommentAction
// @Tags         扩展接口-1
// @Summary      评论操作
// @Param        token query string true "用户鉴权token"
// @Param        video_id query int true "视频id"
// @Param        action_type query int true "1-发布评论，2-删除评论"
// @Param        comment_text query string false "用户填写的评论内容，在action_type=1的时候使用"
// @Param        comment_id query int false "要删除的评论id，在action_type=2的时候使用"
// @Success      200  {string}  json "{"status_code":"200","status_msg":""}"
// @Router       /douyin/comment/action/ [post]
func CommentAction(c *gin.Context) {
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

	//进行评论或取消评论处理
	var statusMsg string
	if actionType == "1" {
		commentText := c.Query("comment_text")
		if commentText == "" {
			c.JSON(http.StatusOK, gin.H{
				"status_code": -1,
				"status_msg":  "缺少comment_text",
			})
			return
		}
		comment, err := store.AddComment(userID, uint(videoIdentity), commentText)
		statusMsg = "评论成功"
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status_code": -1,
				"status_msg":  "add comment error, err: " + err.Error(),
			})
			return
		}
		user, _ := store.GetUser(userID)
		userRes := make(map[string]interface{})
		userRes["id"] = user.ID
		userRes["name"] = user.Username
		userRes["follow_count"] = user.FollowCount
		userRes["follower_count"] = user.FollowerCount
		userRes["is_follow"] = user.IsFollow
		commentRes := make(map[string]interface{})
		commentRes["id"] = comment.ID
		commentRes["user"] = userRes
		commentRes["content"] = comment.Content
		commentRes["create_date"] = comment.CreatedAt.Format("01-02")
		c.JSON(http.StatusOK, gin.H{
			"status_code": 0,
			"status_msg":  statusMsg,
			"comment":     commentRes,
		})
	} else {
		commentId := c.Query("comment_id")
		if commentId == "" {
			c.JSON(http.StatusOK, gin.H{
				"status_code": -1,
				"status_msg":  "缺少comment_id",
			})
			return
		}
		commentIdentity, _ := strconv.Atoi(commentId)
		err := store.DeleteComment(uint(commentIdentity), uint(videoIdentity))
		statusMsg = "删除评论成功"
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status_code": -1,
				"status_msg":  "delete comment error, err: " + err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status_code": 0,
			"status_msg":  statusMsg,
			"comment":     nil,
		})
	}
}

// CommentList
// @Tags         扩展接口-1
// @Summary      评论列表
// @Param        token query string true "用户鉴权token"
// @Param        video_id query int true "视频id"
// @Success      200  {string}  json "{"status_code":"200","status_msg":""}"
// @Router       /douyin/comment/list/ [get]
func CommentList(c *gin.Context) {
	token := c.Query("token")
	videoId := c.Query("video_id")
	if token == "" || videoId == "" {
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

	videoIdentity, _ := strconv.Atoi(videoId)
	commentList := store.GetCommentList(uint(videoIdentity))
	var commentListRes []map[string]interface{}
	for _, comment := range commentList {
		userRes := make(map[string]interface{})
		userRes["id"] = comment.User.ID
		userRes["name"] = comment.User.Username
		userRes["follow_count"] = comment.User.FollowCount
		userRes["follower_count"] = comment.User.FollowerCount
		userRes["is_follow"] = comment.User.IsFollow
		commentRes := make(map[string]interface{})
		commentRes["id"] = comment.ID
		commentRes["user"] = userRes
		commentRes["content"] = comment.Content
		commentRes["create_date"] = comment.CreatedAt.Format("01-02")
		commentListRes = append(commentListRes, commentRes)
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "查询成功",
		"video_list":  commentListRes,
	})
}
