package models

import (
	"github.com/jinzhu/gorm"
)

type SysAdmin struct {
	gorm.Model
	UserName string `gorm:"column:username;type:varchar(64);" json:"username" valid:"Required;"`
	Password string `gorm:"column:password;type:varchar(64);" json:"password" valid:"Required;MinSize(6)"`
	Desc     string `gorm:"column:desc;type:varchar(255);" json:"desc"`
	Status   int    `gorm:"column:status;type:int(2);" json:"status"`
}

func (admin SysAdmin) TableName() string {
	return "sys_admin"
}
