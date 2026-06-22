package config

import (
	"image"
	"image/color"
	"image/png"
	"os"

	"online-ordering-system/models"

	"golang.org/x/crypto/bcrypt"
)

func SeedData() {
	var count int64
	models.DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		return
	}

	os.MkdirAll("static/uploads", 0755)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := models.User{
		Username: "admin",
		Password: string(hashedPassword),
		Role:     "admin",
		Phone:    "13800000000",
		Address:  "管理员地址",
	}
	models.DB.Create(&admin)

	categories := []models.Category{
		{Name: "中餐", Sort: 1},
		{Name: "西餐", Sort: 2},
		{Name: "饮品", Sort: 3},
		{Name: "甜点", Sort: 4},
	}
	models.DB.Create(&categories)

	dishes := []models.Dish{
		{Name: "宫保鸡丁", Price: 28.00, CategoryID: 1, Image: "/static/uploads/default_zhongcan.png", Desc: "经典川菜，麻辣鲜香", Status: true},
		{Name: "红烧肉", Price: 38.00, CategoryID: 1, Image: "/static/uploads/default_zhongcan.png", Desc: "肥而不腻，入口即化", Status: true},
		{Name: "鱼香肉丝", Price: 26.00, CategoryID: 1, Image: "/static/uploads/default_zhongcan.png", Desc: "酸甜微辣，下饭首选", Status: true},
		{Name: "麻婆豆腐", Price: 22.00, CategoryID: 1, Image: "/static/uploads/default_zhongcan.png", Desc: "麻辣鲜烫，嫩滑可口", Status: true},
		{Name: "牛排", Price: 68.00, CategoryID: 2, Image: "/static/uploads/default_xican.png", Desc: "澳洲进口，七分熟", Status: true},
		{Name: "意大利面", Price: 32.00, CategoryID: 2, Image: "/static/uploads/default_xican.png", Desc: "经典番茄肉酱意面", Status: true},
		{Name: "凯撒沙拉", Price: 28.00, CategoryID: 2, Image: "/static/uploads/default_xican.png", Desc: "新鲜蔬菜配凯撒酱", Status: true},
		{Name: "珍珠奶茶", Price: 15.00, CategoryID: 3, Image: "/static/uploads/default_yinpin.png", Desc: "Q弹珍珠，香浓奶茶", Status: true},
		{Name: "鲜榨橙汁", Price: 18.00, CategoryID: 3, Image: "/static/uploads/default_yinpin.png", Desc: "新鲜橙子现榨", Status: true},
		{Name: "拿铁咖啡", Price: 22.00, CategoryID: 3, Image: "/static/uploads/default_yinpin.png", Desc: "意式浓缩配丝滑牛奶", Status: true},
		{Name: "提拉米苏", Price: 28.00, CategoryID: 4, Image: "/static/uploads/default_tiandian.png", Desc: "经典意式甜点", Status: true},
		{Name: "芒果布丁", Price: 18.00, CategoryID: 4, Image: "/static/uploads/default_tiandian.png", Desc: "新鲜芒果制作", Status: true},
	}
	models.DB.Create(&dishes)

	generateDefaultImages()
}

func generateDefaultImages() {
	type imgDef struct {
		filename string
		bgColor  color.RGBA
	}
	imgs := []imgDef{
		{"default_zhongcan.png", color.RGBA{231, 76, 60, 255}},
		{"default_xican.png", color.RGBA{52, 152, 219, 255}},
		{"default_yinpin.png", color.RGBA{46, 204, 113, 255}},
		{"default_tiandian.png", color.RGBA{243, 156, 18, 255}},
	}
	for _, img := range imgs {
		func() {
			f, err := os.Create("static/uploads/" + img.filename)
			if err != nil {
				return
			}
			defer f.Close()
			rect := image.Rect(0, 0, 400, 300)
			rgba := image.NewRGBA(rect)
			for y := 0; y < 300; y++ {
				for x := 0; x < 400; x++ {
					rgba.Set(x, y, img.bgColor)
				}
			}
			png.Encode(f, rgba)
		}()
	}
}
