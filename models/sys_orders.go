package models

import (
	"github.com/jinzhu/gorm"
)

type SysOrders struct {
	gorm.Model
	OrderNum  string `gorm:"column:order_num;type:varchar(64);" json:"order_num" valid:"Required"`
	UId       int    `gorm:"column:uid;type:int(11);" json:"uid"`
	SId       int    `gorm:"column:sid;type:int(11);" json:"sid"`
	PayStatus int    `gorm:"column:pay_status;type:int(2);" json:"pay_status"`
}

func (orders SysOrders) TableName() string {
	return "sys_orders"
}
