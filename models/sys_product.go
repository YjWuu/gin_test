package models

import (
	"github.com/jinzhu/gorm"
)

type SysProduct struct {
	gorm.Model
	Name              string              `gorm:"column:name;type:varchar(64);" json:"name" valid:"Required"`
	Price             float64             `gorm:"column:price;type:decimal(11,2);" json:"price" valid:"Required"`
	Num               int                 `gorm:"column:num;type:int(11);" json:"num"`
	Unit              string              `gorm:"column:unit;type:varchar(32);" json:"unit"`
	Pic               string              `gorm:"column:pic;type:varchar(255);" json:"pic"`
	Desc              string              `gorm:"column:desc;type:varchar(255);" json:"desc"`
	SysProductSeckill []SysProductSeckill `gorm:"ForeignKey:PId;AssiciationForeignKey:Id"`
}

func (product SysProduct) TableName() string {
	return "sys_product"
}
