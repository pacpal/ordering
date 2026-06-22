package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"online-ordering-system/models"

	"github.com/gin-gonic/gin"
)

func ListDishes(c *gin.Context) {
	var dishes []models.Dish
	query := models.DB.Preload("Category")

	if catID := c.Query("category_id"); catID != "" {
		query = query.Where("category_id = ?", catID)
	}
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	query.Find(&dishes)
	c.JSON(http.StatusOK, dishes)
}

func GetDish(c *gin.Context) {
	id := c.Param("id")
	var dish models.Dish
	if err := models.DB.Preload("Category").First(&dish, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "菜品不存在"})
		return
	}
	c.JSON(http.StatusOK, dish)
}

func CreateDish(c *gin.Context) {
	var input struct {
		Name       string  `json:"name" binding:"required"`
		Price      float64 `json:"price" binding:"required"`
		Image      string  `json:"image"`
		Desc       string  `json:"desc"`
		CategoryID uint    `json:"category_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	dish := models.Dish{
		Name:       input.Name,
		Price:      input.Price,
		Image:      input.Image,
		Desc:       input.Desc,
		Status:     true,
		CategoryID: input.CategoryID,
	}
	if err := models.DB.Create(&dish).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}
	c.JSON(http.StatusCreated, dish)
}

func UpdateDish(c *gin.Context) {
	id := c.Param("id")
	var dish models.Dish
	if err := models.DB.First(&dish, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "菜品不存在"})
		return
	}

	var input struct {
		Name       string  `json:"name"`
		Price      float64 `json:"price"`
		Image      string  `json:"image"`
		Desc       string  `json:"desc"`
		Status     *bool   `json:"status"`
		CategoryID uint    `json:"category_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	updates := map[string]interface{}{
		"name": input.Name, "price": input.Price,
		"image": input.Image, "desc": input.Desc,
		"category_id": input.CategoryID,
	}
	if input.Status != nil {
		updates["status"] = *input.Status
	}
	if input.Image != "" && input.Image != dish.Image {
		oldImgPath := strings.TrimPrefix(dish.Image, "/")
		os.Remove(oldImgPath)
	}
	models.DB.Model(&dish).Updates(updates)
	c.JSON(http.StatusOK, dish)
}

func DeleteDish(c *gin.Context) {
	id := c.Param("id")
	var dish models.Dish
	if err := models.DB.First(&dish, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "菜品不存在"})
		return
	}
	if dish.Image != "" {
		imgPath := strings.TrimPrefix(dish.Image, "/")
		os.Remove(imgPath)
	}
	models.DB.Delete(&dish)
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func UploadDishImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择图片"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "仅支持 JPG/PNG/GIF/WEBP 格式"})
		return
	}

	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "图片大小不能超过 5MB"})
		return
	}

	uploadDir := "static/uploads"
	os.MkdirAll(uploadDir, 0755)

	timestamp := time.Now().Format("20060102150405")
	randomStr := strconv.FormatInt(time.Now().UnixNano()%10000, 10)
	filename := fmt.Sprintf("%s_%s%s", timestamp, randomStr, ext)
	savePath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": "/" + filepath.ToSlash(savePath)})
}
