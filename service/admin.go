package service

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"log"
	"net/http"
	"persion_test/define"
	"persion_test/helper"
	"persion_test/models"
	"strconv"
	"time"
)

// AdminLogin
// @Tags 公共方法
// @Summary 管理员登录
// @Param username formData string true "username"
// @Param password formData string true "password"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin-login [post]
func AdminLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "必填信息为空",
		})
		return
	}
	password = helper.GetMd5(password)
	data := new(models.SysAdmin)
	err := models.DB.Where("username = ? and password = ?", username, password).First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": -1,
				"msg":  "用户名或密码错误",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "Get Admin Error:" + err.Error(),
		})
		return
	}
	token, err := helper.GenerateToken(data.UserName, data.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "GenerateToken error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": token,
	})
}

// GetUserList
// @Tags 公有方法
// @Summary 获取用户列表
// @Param page query int false "page"
// @Param size query int false "size"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /user-list [get]
func GetUserList(c *gin.Context) {
	size, err := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	if err != nil {
		log.Println("GetUserList size strconv Error" + err.Error())
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("GetUserList Page strconv Error" + err.Error())
		return
	}
	page = (page - 1) * size
	var count int64
	list := make([]*models.SysUser, 0)
	err = models.DB.Preload("SysOrders").Model(new(models.SysUser)).Count(&count).Offset(page).Limit(size).Find(&list).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "Get UserList Error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  list,
			"count": count,
		},
	})
}

// UploadImg
// @Tags 公共方法
// @Summary 上传图片
// @Param pic  formData file true "pic"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /upload-img [post]
func UploadImg(c *gin.Context) {
	file, err := c.FormFile("pic")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "获取文件失败",
		})
	}
	url := helper.UploadImg(file)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": url,
	})
}

// CreateProduct
// @Tags 管理员私有方法
// @Summary 创建产品
// @Param authorization header string true "authorization"
// @Param name formData string true "name"
// @Param price formData float64 true "price"
// @Param num formData int true "num"
// @Param unit formData string false "unit"
// @Param pic  formData file true "pic"
// @Param desc  formData string false "desc"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/create-product [post]
func CreateProduct(c *gin.Context) {
	name := c.PostForm("name")
	price := c.PostForm("price")
	num := c.PostForm("num")
	unit := c.PostForm("unit")
	desc := c.PostForm("desc")
	file, _ := c.FormFile("pic")
	pic := helper.UploadImg(file)
	price1, _ := strconv.ParseFloat(price, 64)
	num1, _ := strconv.Atoi(num)
	product := &models.SysProduct{
		Name:  name,
		Price: price1,
		Num:   num1,
		Unit:  unit,
		Desc:  desc,
		Pic:   pic,
	}
	err := models.DB.Create(product).Error
	if err != nil {
		log.Println("Product Create Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "创建产品失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建成功",
	})
}

// DeleteProduct
// @Tags 管理员私有方法
// @Summary 删除产品
// @Param authorization header string true "authorization"
// @Param name query string true "name"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/delete-product [delete]
func DeleteProduct(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "请输入产品名",
		})
		return
	}
	var count int64
	err := models.DB.Model(new(models.SysProduct)).Where("name = ?", name).Count(&count).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "通过Id获取信息失败",
		})
		return
	}
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "该产品不存在",
		})
		return
	}
	err = models.DB.Where("name = ?", name).Delete(new(models.SysProduct)).Error
	if err != nil {
		log.Println("Delete Product Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "删除失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

// GetProductList
// @Tags 管理员私有方法
// @Summary 获取产品列表
// @Param authorization header string true "authorization"
// @Param page query int false "page"
// @Param size query int false "size"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/product-list [get]
func GetProductList(c *gin.Context) {
	size, err := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	if err != nil {
		log.Println("GetProductList size strconv Error" + err.Error())
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("GetProductList Page strconv Error" + err.Error())
		return
	}
	page = (page - 1) * size
	var count int64
	list := make([]*models.SysProduct, 0)
	err = models.DB.Preload("SysProductSeckill").Preload("SysProductSeckill.SysOrders").Model(new(models.SysProduct)).Count(&count).Limit(size).Offset(page).Find(&list).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "Get ProductList Error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  list,
			"count": count,
		},
	})
}

// UpdateProduct
// @Tags 管理员私有方法
// @Summary 更新产品
// @Param authorization header string true "authorization"
// @param id formData int true "id"
// @Param name formData string false "name"
// @Param price formData float64 false "price"
// @Param num formData int false "num"
// @Param unit formData string false "unit"
// @Param pic  formData file false "pic"
// @Param desc  formData string false "desc"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/update-product [put]
func UpdateProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("id"))
	name := c.PostForm("name")
	price := c.PostForm("price")
	num := c.PostForm("num")
	unit := c.PostForm("unit")
	desc := c.PostForm("desc")
	file, _ := c.FormFile("pic")
	pic := helper.UploadImg(file)
	price1, _ := strconv.ParseFloat(price, 64)
	num1, _ := strconv.Atoi(num)
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "请输入产品Id",
		})
		return
	}
	product := &models.SysProduct{
		Name:  name,
		Price: price1,
		Num:   num1,
		Unit:  unit,
		Desc:  desc,
		Pic:   pic,
	}
	var count int64
	err := models.DB.Model(new(models.SysProduct)).Where("id = ?", id).Count(&count).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "通过Id获取信息失败",
		})
		return
	}
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "该产品不存在",
		})
		return
	}
	err = models.DB.Model(new(models.SysProduct)).Where("id = ?", id).Update(product).Error
	if err != nil {
		log.Println("Update Product Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "修改产品信息失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "修改成功",
	})
}

// GetProductImgUrl
// @Tags 管理员私有方法
// @Summary 获取产品图片链接
// @Param authorization header string true "authorization"
// @Param id query int true "id"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/product-img-url [get]
func GetProductImgUrl(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	data := new(models.SysProduct)
	err := models.DB.Model(new(models.SysProduct)).Where("id = ?", id).Find(&data).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "该产品不存在",
		})
		return
	}
	mac := qbox.NewMac(define.AccessKey, define.SecretKey)
	key := data.Pic
	deadline := time.Now().Add(time.Second * 3600).Unix() //1小时有效期
	privateAccessURL := storage.MakePrivateURL(mac, define.ImgUrl, key, deadline)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": privateAccessURL,
	})
}

// CreateProductSecKill
// @Tags 管理员私有方法
// @Summary 创建秒杀
// @Param authorization header string true "authorization"
// @Param name formData string true "name"
// @Param price formData float64 true "price"
// @Param num formData int true "num"
// @Param pid formData int true "pid"
// @Param startTime  formData string true "startTime"
// @Param endTime  formData string false "endTime"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/create-product-seckill [post]
func CreateProductSecKill(c *gin.Context) {
	name := c.PostForm("name")
	price := c.PostForm("price")
	num := c.PostForm("num")
	pid := c.PostForm("pid")
	startTime := c.PostForm("startTime")
	endTime := c.PostForm("endTime")
	price1, _ := strconv.ParseFloat(price, 64)
	num1, _ := strconv.Atoi(num)
	pid1, _ := strconv.Atoi(pid)
	timeLayout := "2006-01-02 15:04:05"  //转化所需模板
	loc, _ := time.LoadLocation("Local") //获取时区
	tmp, _ := time.ParseInLocation(timeLayout, startTime, loc)
	tmp2, _ := time.ParseInLocation(timeLayout, endTime, loc)
	if tmp.After(tmp2) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "开始时间应该早于结束时间",
		})
		return
	}
	productSeckill := &models.SysProductSeckill{
		Name:      name,
		Price:     price1,
		Num:       num1,
		PId:       pid1,
		StartTime: tmp,
		EndTime:   tmp2,
	}
	err := models.DB.Create(productSeckill).Error
	if err != nil {
		log.Println("ProductSeckill Create Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "创建秒杀失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建成功",
	})
}

// DeleteProductSecKill
// @Tags 管理员私有方法
// @Summary 删除秒杀
// @Param authorization header string true "authorization"
// @Param id query string true "id"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/delete-product-seckill [delete]
func DeleteProductSecKill(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "请输入产品名",
		})
		return
	}
	var count int64
	err := models.DB.Model(new(models.SysProductSeckill)).Where("id = ?", id).Count(&count).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "通过Id获取信息失败",
		})
		return
	}
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "该产品不存在",
		})
		return
	}
	err = models.DB.Where("id = ?", id).Delete(new(models.SysProductSeckill)).Error
	if err != nil {
		log.Println("Delete ProductSeckill Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "删除失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

// GetProductSecKillList
// @Tags 管理员私有方法
// @Summary 获取秒杀列表
// @Param authorization header string true "authorization"
// @Param page query int false "page"
// @Param size query int false "size"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/product-seckill-list [get]
func GetProductSecKillList(c *gin.Context) {
	size, err := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	if err != nil {
		log.Println("GetProductSeckillList size strconv Error" + err.Error())
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("GetProductSeckillList Page strconv Error" + err.Error())
		return
	}
	page = (page - 1) * size
	var count int64
	list := make([]*models.SysProductSeckill, 0)
	err = models.DB.Preload("SysOrders").Model(new(models.SysProductSeckill)).Count(&count).Limit(size).Offset(page).Find(&list).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "Get ProductSeckillList Error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  list,
			"count": count,
		},
	})
}

// UpdateProductSecKill
// @Tags 管理员私有方法
// @Summary 更新秒杀
// @Param authorization header string true "authorization"
// @param id formData int true "id"
// @Param name formData string true "name"
// @Param price formData float64 true "price"
// @Param num formData int true "num"
// @Param pid formData int true "pid"
// @Param startTime  formData string true "startTime"
// @Param endTime  formData string false "endTime"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/update-product-seckill [put]
func UpdateProductSecKill(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("id"))
	name := c.PostForm("name")
	price := c.PostForm("price")
	num := c.PostForm("num")
	pid := c.PostForm("pid")
	startTime := c.PostForm("startTime")
	endTime := c.PostForm("endTime")
	price1, _ := strconv.ParseFloat(price, 64)
	num1, _ := strconv.Atoi(num)
	pid1, _ := strconv.Atoi(pid)
	timeLayout := "2006-01-02 15:04:05"  //转化所需模板
	loc, _ := time.LoadLocation("Local") //获取时区
	tmp, _ := time.ParseInLocation(timeLayout, startTime, loc)
	tmp2, _ := time.ParseInLocation(timeLayout, endTime, loc)
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "请输入产品Id",
		})
		return
	}
	if tmp.After(tmp2) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "开始时间应该早于结束时间",
		})
		return
	}
	productSeckill := &models.SysProductSeckill{
		Name:      name,
		Price:     price1,
		Num:       num1,
		PId:       pid1,
		StartTime: tmp,
		EndTime:   tmp2,
	}
	var count int64
	err := models.DB.Model(new(models.SysProductSeckill)).Where("id = ?", id).Count(&count).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "通过Id获取信息失败",
		})
		return
	}
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "该秒杀不存在",
		})
		return
	}
	err = models.DB.Model(new(models.SysProductSeckill)).Where("id = ?", id).Update(productSeckill).Error
	if err != nil {
		log.Println("Update ProductSeckill Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "修改秒杀信息失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "修改成功",
	})
}
