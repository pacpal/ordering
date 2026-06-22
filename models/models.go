package models

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(dbname string) {
	db, err := gorm.Open(sqlite.Open(dbname), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB = db
	DB.AutoMigrate(&User{}, &Category{}, &Dish{}, &CartItem{}, &Order{}, &OrderItem{})
}

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password string `gorm:"size:255;not null" json:"-"`
	Role     string `gorm:"size:20;default:customer" json:"role"`
	Phone    string `gorm:"size:20" json:"phone"`
	Address  string `gorm:"size:255" json:"address"`
}

type Category struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Sort int    `gorm:"default:0" json:"sort"`
}

type Dish struct {
	ID         uint     `gorm:"primaryKey" json:"id"`
	Name       string   `gorm:"size:100;not null" json:"name"`
	Price      float64  `gorm:"not null" json:"price"`
	Image      string   `gorm:"size:255" json:"image"`
	Desc       string   `gorm:"size:500" json:"desc"`
	Status     bool     `gorm:"default:true" json:"status"`
	CategoryID uint     `json:"category_id"`
	Category   Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

type CartItem struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"not null;index" json:"user_id"`
	DishID uint `gorm:"not null" json:"dish_id"`
	Count  int  `gorm:"default:1" json:"count"`
	Dish   Dish `gorm:"foreignKey:DishID" json:"dish,omitempty"`
}

type Order struct {
	ID      uint        `gorm:"primaryKey" json:"id"`
	UserID  uint        `gorm:"not null;index" json:"user_id"`
	Total   float64     `gorm:"not null" json:"total"`
	Status  string      `gorm:"size:20;default:pending" json:"status"`
	Address string      `gorm:"size:255" json:"address"`
	Phone   string      `gorm:"size:20" json:"phone"`
	Remark  string      `gorm:"size:500" json:"remark"`
	Items   []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	User    User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type OrderItem struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	OrderID  uint    `gorm:"not null;index" json:"order_id"`
	DishID   uint    `gorm:"not null" json:"dish_id"`
	DishName string  `gorm:"size:100" json:"dish_name"`
	Price    float64 `gorm:"not null" json:"price"`
	Count    int     `gorm:"not null" json:"count"`
}
