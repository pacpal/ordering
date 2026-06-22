package controllers

import (
	"net/http"

	"online-ordering-system/models"

	"github.com/gin-gonic/gin"
)

func CreateOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var cartItems []models.CartItem
	models.DB.Where("user_id = ?", userID).Preload("Dish").Find(&cartItems)
	if len(cartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "购物车为空"})
		return
	}

	var input struct {
		Address string `json:"address"`
		Phone   string `json:"phone"`
		Remark  string `json:"remark"`
	}
	c.ShouldBindJSON(&input)

	var user models.User
	models.DB.First(&user, userID)

	address := input.Address
	if address == "" {
		address = user.Address
	}
	phone := input.Phone
	if phone == "" {
		phone = user.Phone
	}

	var total float64
	var orderItems []models.OrderItem
	for _, item := range cartItems {
		total += item.Dish.Price * float64(item.Count)
		orderItems = append(orderItems, models.OrderItem{
			DishID:   item.DishID,
			DishName: item.Dish.Name,
			Price:    item.Dish.Price,
			Count:    item.Count,
		})
	}

	order := models.Order{
		UserID:  userID.(uint),
		Total:   total,
		Status:  "pending",
		Address: address,
		Phone:   phone,
		Remark:  input.Remark,
		Items:   orderItems,
	}

	if err := models.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "下单失败"})
		return
	}

	models.DB.Where("user_id = ?", userID).Delete(&models.CartItem{})
	c.JSON(http.StatusCreated, gin.H{"message": "下单成功", "order": order})
}

func ListMyOrders(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var orders []models.Order
	models.DB.Where("user_id = ?", userID).Preload("Items").Order("id DESC").Find(&orders)
	c.JSON(http.StatusOK, orders)
}

func ListAllOrders(c *gin.Context) {
	var orders []models.Order
	models.DB.Preload("Items").Preload("User").Order("id DESC").Find(&orders)
	c.JSON(http.StatusOK, orders)
}

func UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	var order models.Order
	if err := models.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	var input struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	validStatuses := map[string]bool{"pending": true, "confirmed": true, "preparing": true, "delivering": true, "completed": true, "cancelled": true}
	if !validStatuses[input.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的订单状态"})
		return
	}

	models.DB.Model(&order).Update("status", input.Status)
	c.JSON(http.StatusOK, order)
}

func CancelOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")

	var order models.Order
	if err := models.DB.Where("id = ? AND user_id = ?", id, userID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	if order.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只能取消待处理的订单"})
		return
	}

	models.DB.Model(&order).Update("status", "cancelled")
	c.JSON(http.StatusOK, gin.H{"message": "订单已取消"})
}
