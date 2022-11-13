package router

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "persion_test/docs"
	"persion_test/middlewares"
	"persion_test/models"
	"persion_test/service"
)

func Router() *gin.Engine {
	models.InitDb()
	models.InitRedis()
	r := gin.Default()
	r.Use(middlewares.Cors(), middlewares.LoggerFile())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.POST("/admin-login", service.AdminLogin)
	r.POST("/user-login", service.UserLogin)
	r.POST("/register", service.Register)
	r.POST("/send-code", service.SendCode)
	r.POST("/upload-img", service.UploadImg)
	r.GET("/image-code", service.GetImageCode)
	r.GET("/user-list", service.GetUserList)
	authAdmin := r.Group("/admin", middlewares.AuthAdminCheck())
	//authAdmin.GET("/user-list", service.GetUserList)
	authAdmin.POST("/create-product", service.CreateProduct)
	authAdmin.PUT("/update-product", service.UpdateProduct)
	authAdmin.DELETE("/delete-product", service.DeleteProduct)
	authAdmin.GET("/product-list", service.GetProductList)
	authAdmin.GET("/product-img-url", service.GetProductImgUrl)
	authAdmin.POST("/create-product-seckill", service.CreateProductSecKill)
	authAdmin.PUT("/update-product-seckill", service.UpdateProductSecKill)
	authAdmin.DELETE("/delete-product-seckill", service.DeleteProductSecKill)
	authAdmin.GET("/product-seckill-list", service.GetProductSecKillList)

	authUser := r.Group("/user", middlewares.AuthUserCheck())
	authUser.POST("/create-order", service.CreateOrder)
	authUser.PUT("/update-order", service.UpdateOrder)
	authUser.DELETE("/delete-order", service.DeleteOrder)
	return r
}
