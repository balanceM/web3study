package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title   string
	Content string
	UserID  uint
	User    User
}

type CreatePostReq struct {
	Title   string `form:"title" binding:"required,min=1,max=100"`
	Content string `form:"content" binding:"required,min=1"`
}

func getCurrentUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "can't get user"})
		return 0, false
	}
	uid, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID error"})
		return 0, false
	}
	return uid, true
}

func getPostAndCheckOwner(c *gin.Context, postID string, userID uint) (*Post, bool) {
	var post Post
	if err := db.Where("id = ?", postID).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "can't get post"})
		return nil, false
	}
	if post.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "post is not belongs to the user"})
		return nil, false
	}
	return &post, true
}

func validatePostID(c *gin.Context) (string, bool) {
	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "post id is null"})
		return "", false
	}
	return postID, true
}

func CreatePostHandler(c *gin.Context) {
	var req CreatePostReq
	if err := c.ShouldBind(&req); err != nil {
		zap.L().Error("CreatePost failed", zap.String("error", err.Error()), zap.String("path", c.Request.URL.Path), zap.String("method", c.Request.Method))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid, ok := getCurrentUserID(c)
	if !ok {
		zap.L().Error("CreatePost failed", zap.String("error", "can't get user id"), zap.String("path", c.Request.URL.Path), zap.String("method", c.Request.Method))
		return
	}

	post := Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  uid,
	}

	if err := db.Create(&post).Error; err != nil {
		zap.L().Error("CreatePost failed", zap.String("error", err.Error()), zap.String("path", c.Request.URL.Path), zap.String("method", c.Request.Method))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create post failed"})
		return
	}

	zap.L().Info("CreatePost successfully", zap.Uint("post_id", post.ID), zap.Uint("user_id", uid))
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"post": gin.H{
			"id":      post.ID,
			"title":   post.Title,
			"content": post.Content,
			"user_id": post.UserID,
			"created": post.CreatedAt,
		},
	})
}

type UpdatePostReq struct {
	Title   string `form:"title"`
	Content string `form:"content"`
}

func UpdatePostHandler(c *gin.Context) {
	postID, ok := validatePostID(c)
	if !ok {
		return
	}

	var req UpdatePostReq
	if err := c.ShouldBind(&req); err != nil {
		zap.L().Error("UpdatePost failed", zap.String("error", err.Error()), zap.String("path", c.Request.URL.Path), zap.String("method", c.Request.Method))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid, ok := getCurrentUserID(c)
	if !ok {
		zap.L().Error("UpdatePost failed", zap.String("error", "can't get user id"), zap.String("path", c.Request.URL.Path), zap.String("method", c.Request.Method))
		return
	}
	post, ok := getPostAndCheckOwner(c, postID, uid)
	if !ok {
		return
	}

	updateData := make(map[string]interface{})
	if req.Title != "" {
		updateData["Title"] = req.Title
	}
	if req.Content != "" {
		updateData["Content"] = req.Title
	}
	if err := db.Model(&post).Updates(updateData).Error; err != nil {
		zap.L().Error("UpdatePost failed", zap.String("error", err.Error()), zap.String("path", c.Request.URL.Path), zap.String("method", c.Request.Method))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	zap.L().Info("UpdatePost successfully", zap.Uint("post_id", post.ID), zap.Uint("user_id", uid))
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"post": gin.H{
			"id":      post.ID,
			"title":   post.Title,
			"content": post.Content,
			"updated": post.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func GetPostHandler(c *gin.Context) {
	postID, ok := validatePostID(c)
	if !ok {
		zap.L().Error("GetPost failed", zap.String("error", "valid postID failed"), zap.String("path", c.Request.URL.Path), zap.String("method", c.Request.Method))
		return
	}

	var post Post
	if err := db.Where("id = ?", postID).First(&post).Error; err != nil {
		zap.L().Error("GetPost failed", zap.String("error", "can't get post"), zap.String("path", c.Request.URL.Path), zap.String("method", c.Request.Method))
		c.JSON(http.StatusNotFound, gin.H{"error": "can't get post"})
		return
	}

	zap.L().Info("GetPost successfully", zap.Uint("post_id", post.ID))
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"post": gin.H{
			"id":      post.ID,
			"title":   post.Title,
			"content": post.Content,
			"created": post.CreatedAt.Format("2006-01-02 15:04:05"),
			"updated": post.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func DeletePostHandler(c *gin.Context) {
	postID, ok := validatePostID(c)
	if !ok {
		zap.L().Error("DelPost failed", zap.String("error", "validatePostID failed"), zap.String("path", c.Request.URL.Path), zap.String("method", c.Request.Method))
		return
	}

	uid, ok := getCurrentUserID(c)
	if !ok {
		zap.L().Error("DelPost failed", zap.String("error", "getCurrentUserID failed"), zap.String("path", c.Request.URL.Path), zap.String("method", c.Request.Method))
		return
	}

	post, ok := getPostAndCheckOwner(c, postID, uid)
	if !ok {
		zap.L().Error("DelPost failed", zap.String("error", "getPostAndCheckOwner failed"), zap.String("path", c.Request.URL.Path), zap.String("method", c.Request.Method))
		return
	}

	if err := db.Delete(&post).Error; err != nil {
		zap.L().Error("DelPost failed", zap.String("error", err.Error()), zap.String("path", c.Request.URL.Path), zap.String("method", c.Request.Method))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	zap.L().Info("DelPost successfully", zap.Uint("post_id", post.ID))
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"post_id": post.ID,
	})
}
