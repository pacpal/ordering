package controllers

import (
	"net/http"

	"online-ordering-system/models"

	"github.com/gin-gonic/gin"
)

func ListCategories(c *gin.Context) {
	var categories []models.Category
	models.DB.Order("sort ASC, id ASC").Find(&categories)
	c.JSON(http.StatusOK, categories)
}

func CreateCategory(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
		Sort int    `json:"sort"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	category := models.Category{Name: input.Name, Sort: input.Sort}
	if err := models.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}
	c.JSON(http.StatusCreated, category)
}

func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	if err := models.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
		return
	}

	var input struct {
		Name string `json:"name"`
		Sort int    `json:"sort"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	models.DB.Model(&category).Updates(map[string]interface{}{
		"name": input.Name,
		"sort": input.Sort,
	})
	c.JSON(http.StatusOK, category)
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	models.DB.Delete(&models.Category{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
