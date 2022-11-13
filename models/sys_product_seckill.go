package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type SysProductSeckill struct {
	gorm.Model
	Name      string      `gorm:"column:name;type:varchar(64);" json:"name" valid:"Required"`
	Price     float64     `gorm:"column:price;type:decimal(11,2);" json:"price" valid:"Required"`
	Num       int         `gorm:"column:num;type:int(11);" json:"num"`
	PId       int         `gorm:"column:pid;type:int(11)" json:"pid"`
	StartTime time.Time   `gorm:"column:start_time;" json:"start_time"`
	EndTime   time.Time   `gorm:"column:end_time;" json:"end_time"`
	SysOrders []SysOrders `gorm:"ForeignKey:SId;AssiciationForeignKey:Id"`
}

func (productkill SysProductSeckill) TableName() string {
	return "sys_product_seckill"
}
