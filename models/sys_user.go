package models

import (
	"github.com/jinzhu/gorm"
)

type SysUser struct {
	gorm.Model
	Phone     string      `gorm:"column:phone;type:varchar(64);" json:"phone" valid:"Required;Phone"`
	Password  string      `gorm:"column:password;type:varchar(64);" json:"password" valid:"Required;MinSize(6)"`
	Desc      string      `gorm:"column:desc;type:varchar(255);" json:"desc"`
	Status    int         `gorm:"column:status;type:int(2);" json:"status"`
	SysOrders []SysOrders `gorm:"ForeignKey:UId;AssiciationForeignKey:Id"`
}

func (user SysUser) TableName() string {
	return "sys_user"
}
