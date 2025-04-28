package seed

import (
	"control-go/model"

	"gorm.io/gorm"
)

func SeedProduct(db *gorm.DB) error {
	var count int64
	db.Model(&model.Product{}).Count(&count)
	if count > 0 {
		return nil // 表不为空，跳过插入
	}

	products := []model.Product{
		{
			Code:  "P001",
			Price: 100,
		},
		{
			Code:  "P002",
			Price: 200,
		},
		{
			Code:  "P003",
			Price: 300,
		},
	}

	for _, product := range products {
		if err := db.Create(&product).Error; err != nil {
			return err
		}
	}
	return nil
}
