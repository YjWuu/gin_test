package models

import (
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB
var Redis *redis.Client

func InitDb() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", "root:admin@tcp(127.0.0.1:3306)/person_test?parseTime=True&loc=Local")

	if err != nil {
		return nil, err
	}
	DB = db
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetConnMaxLifetime(100)
	db.SingularTable(true)
	db.AutoMigrate(new(SysAdmin), new(SysUser), new(SysOrders), new(SysProduct), new(SysProductSeckill))
	return db, nil
}

func InitRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
