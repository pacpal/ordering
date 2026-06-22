package routes

import (
	"online-ordering-system/controllers"
	"online-ordering-system/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	r.GET("/admin", func(c *gin.Context) {
		c.HTML(200, "admin.html", nil)
	})

	api := r.Group("/api")
	{
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)
		api.GET("/dishes", controllers.ListDishes)
		api.GET("/dishes/:id", controllers.GetDish)
		api.GET("/categories", controllers.ListCategories)
	}

	auth := api.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/user", controllers.GetCurrentUser)
		auth.PUT("/user", controllers.UpdateProfile)
		auth.GET("/cart", controllers.ListCart)
		auth.POST("/cart", controllers.AddToCart)
		auth.PUT("/cart/:id", controllers.UpdateCartItem)
		auth.DELETE("/cart/:id", controllers.RemoveFromCart)
		auth.DELETE("/cart", controllers.ClearCart)
		auth.POST("/orders", controllers.CreateOrder)
		auth.GET("/orders", controllers.ListMyOrders)
		auth.PUT("/orders/:id/cancel", controllers.CancelOrder)
	}

	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		admin.GET("/users", controllers.ListUsers)
		admin.POST("/categories", controllers.CreateCategory)
		admin.PUT("/categories/:id", controllers.UpdateCategory)
		admin.DELETE("/categories/:id", controllers.DeleteCategory)
		admin.POST("/dishes", controllers.CreateDish)
		admin.PUT("/dishes/:id", controllers.UpdateDish)
		admin.DELETE("/dishes/:id", controllers.DeleteDish)
		admin.POST("/upload", controllers.UploadDishImage)
		admin.GET("/orders", controllers.ListAllOrders)
		admin.PUT("/orders/:id/status", controllers.UpdateOrderStatus)
	}

	return r
}
