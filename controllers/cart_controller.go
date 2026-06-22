package controllers

import (
	"net/http"

	"online-ordering-system/models"

	"github.com/gin-gonic/gin"
)

// ListCart 获取当前用户的购物车列表
// 从认证中间件获取用户ID，查询该用户的所有购物车项并预加载关联的菜品信息
func ListCart(c *gin.Context) {
	// 从上下文中获取已认证的用户ID（由AuthMiddleware设置）
	userID, _ := c.Get("user_id")
	// 声明购物车项切片用于存储查询结果
	var items []models.CartItem
	// 查询当前用户的购物车项，并预加载关联的菜品详情
	models.DB.Where("user_id = ?", userID).Preload("Dish").Find(&items)
	// 返回购物车列表JSON
	c.JSON(http.StatusOK, items)
}

// AddToCart 添加菜品到购物车
// 若该菜品已在购物车中则累加数量，否则新建购物车项
func AddToCart(c *gin.Context) {
	// 从上下文中获取已认证的用户ID
	userID, _ := c.Get("user_id")
	// 定义请求参数结构体，dish_id为必填，count可选
	var input struct {
		DishID uint `json:"dish_id" binding:"required"`
		Count  int  `json:"count"`
	}
	// 绑定并校验JSON请求参数
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 若未指定数量或数量不合法，默认为1
	if input.Count <= 0 {
		input.Count = 1
	}

	// 查询该用户的购物车中是否已存在相同菜品
	var existing models.CartItem
	if err := models.DB.Where("user_id = ? AND dish_id = ?", userID, input.DishID).First(&existing).Error; err == nil {
		// 已存在则累加数量
		models.DB.Model(&existing).Update("count", existing.Count+input.Count)
		c.JSON(http.StatusOK, gin.H{"message": "已添加到购物车", "item": existing})
		return
	}

	// 不存在则创建新的购物车项
	item := models.CartItem{
		UserID: userID.(uint),
		DishID: input.DishID,
		Count:  input.Count,
	}
	models.DB.Create(&item)
	c.JSON(http.StatusCreated, item)
}

// UpdateCartItem 修改购物车中某项的数量
// 仅允许修改属于当前用户的购物车项，数量必须 ≥ 1
func UpdateCartItem(c *gin.Context) {
	// 从上下文中获取已认证的用户ID
	userID, _ := c.Get("user_id")
	// 从URL路径中获取购物车项ID
	id := c.Param("id")

	// 查询该购物车项，确保属于当前用户
	var item models.CartItem
	if err := models.DB.Where("id = ? AND user_id = ?", id, userID).First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "购物车项不存在"})
		return
	}

	// 绑定请求参数，count必填且最小值为1
	var input struct {
		Count int `json:"count" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 更新购物车项数量
	models.DB.Model(&item).Update("count", input.Count)
	c.JSON(http.StatusOK, item)
}

// RemoveFromCart 移除购物车中的指定项
// 仅允许移除属于当前用户的购物车项
func RemoveFromCart(c *gin.Context) {
	// 从上下文中获取已认证的用户ID
	userID, _ := c.Get("user_id")
	// 从URL路径中获取购物车项ID
	id := c.Param("id")
	// 删除该用户指定的购物车项
	models.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.CartItem{})
	c.JSON(http.StatusOK, gin.H{"message": "已移除"})
}

// ClearCart 清空当前用户的购物车
func ClearCart(c *gin.Context) {
	// 从上下文中获取已认证的用户ID
	userID, _ := c.Get("user_id")
	// 删除该用户的所有购物车项
	models.DB.Where("user_id = ?", userID).Delete(&models.CartItem{})
	c.JSON(http.StatusOK, gin.H{"message": "购物车已清空"})
}
