package service

import (
	"github.com/afocus/captcha"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"persion_test/helper"
	"persion_test/models"
	"strconv"
	"time"
)

// SendCode
// @Tags 公共方法
// @Summary 发送验证码
// @Param uuid formData string true "uuid"
// @Param imgCode formData string true "imgCode"
// @Param phone formData string true "phone"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /send-code [post]
func SendCode(c *gin.Context) {
	uuid := c.PostForm("uuid")
	imgCode := c.PostForm("imgCode")
	phone := c.PostForm("phone")
	result := helper.CheckImgCode(c, uuid, imgCode)
	if !result {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "验证码错误",
		})
		return
	}
	if phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "手机号不能为空",
		})
		return
	}
	code := helper.GetRand()
	models.Redis.Set(c, phone, code, time.Second*300)
	err := helper.SendCode(phone, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "Send Code Error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": "验证码发送成功",
	})
}

// Register
// @Tags 公共方法
// @Summary 用户注册
// @Param phone formData string true "phone"
// @Param code formData string true "code"
// @Param password formData string true "password"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /register [post]
func Register(c *gin.Context) {
	phone := c.PostForm("phone")
	userCode := c.PostForm("code")
	password := c.PostForm("password")
	if phone == "" || userCode == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "请重新输入参数",
		})
		return
	}
	sysCode, err := models.Redis.Get(c, phone).Result()
	if err != nil {
		log.Printf("Get code error: %v \n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "获取验证码失败，请重新获取",
		})
		return
	}
	if sysCode != userCode {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "验证码错误，请重新输入",
		})
		return
	}
	var count int
	err = models.DB.Where("phone = ?", phone).Model(new(models.SysUser)).Count(&count).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "Get User Error:" + err.Error(),
		})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "该手机号已被注册",
		})
		return
	}

	data := models.SysUser{
		Phone:    phone,
		Password: helper.GetMd5(password),
	}
	err = models.DB.Create(&data).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "create user error" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": "注册成功",
	})
}

// UserLogin
// @Tags 公共方法
// @Summary 用户登录
// @Param phone formData string true "phone"
// @Param password formData string true "password"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /user-login [post]
func UserLogin(c *gin.Context) {
	phone := c.PostForm("phone")
	password := c.PostForm("password")
	if phone == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "请输入必填信息",
		})
		return
	}
	password = helper.GetMd5(password)
	data := new(models.SysUser)
	err := models.DB.Where("phone = ? and password = ?", phone, password).First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": -1,
				"msg":  "手机号或密码错误",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "Get UserBasic Error:" + err.Error(),
		})
		return
	}
	token, err := helper.GenerateToken(data.Phone, data.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "GenerateToken Error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": token,
	})
}

// GetImageCode
// @Tags 公共方法
// @Summary 获取图片验证码
// @Param uuid query string true "uuid"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /image-code [get]
func GetImageCode(c *gin.Context) {
	uuid := c.Query("uuid")
	cap := captcha.New()

	// 设置字体
	cap.SetFont("./conf/comic.ttf")

	// 设置验证码大小
	cap.SetSize(128, 64)

	// 设置干扰强度
	cap.SetDisturbance(captcha.MEDIUM)

	// 设置前景色
	cap.SetFrontColor(color.RGBA{0, 0, 0, 255})

	// 设置背景色
	cap.SetBkgColor(color.RGBA{100, 0, 255, 255}, color.RGBA{255, 0, 127, 255}, color.RGBA{255, 255, 10, 255})

	// 生成字体 -- 将图片验证码, 展示到页面中.
	img, str := cap.Create(4, captcha.NUM)

	models.Redis.Set(c, uuid, str, time.Second*300)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": png.Encode(c.Writer, img),
	})
}

// CreateOrder
// @Tags 用户私有方法
// @Summary 创建订单
// @Param authorization header string true "authorization"
// @Param orderNum formData string true "orderNum"
// @Param uid formData int true "uid"
// @Param sid formData int true "sid"
// @Param payStatus formData int false "payStatus"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /user/create-order [post]
func CreateOrder(c *gin.Context) {
	orderNum := c.PostForm("orderNum")
	uid, _ := strconv.Atoi(c.PostForm("uid"))
	sid, _ := strconv.Atoi(c.PostForm("sid"))
	payStatus, _ := strconv.Atoi(c.PostForm("payStatus"))
	order := &models.SysOrders{
		OrderNum:  orderNum,
		UId:       uid,
		SId:       sid,
		PayStatus: payStatus,
	}
	err := models.DB.Create(order).Error
	if err != nil {
		log.Println("Order Create Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "创建订单失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建成功",
	})
}

// UpdateOrder
// @Tags 用户私有方法
// @Summary 更新订单
// @Param authorization header string true "authorization"
// @Param id  formData string true "id"
// @Param orderNum formData string false "orderNum"
// @Param uid formData int false "uid"
// @Param sid formData int false "sid"
// @Param payStatus formData int false "payStatus"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /user/update-order [put]
func UpdateOrder(c *gin.Context) {
	id := c.PostForm("id")
	orderNum := c.PostForm("orderNum")
	uid, _ := strconv.Atoi(c.PostForm("uid"))
	sid, _ := strconv.Atoi(c.PostForm("sid"))
	payStatus, _ := strconv.Atoi(c.PostForm("payStatus"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "请输入产品Id",
		})
		return
	}
	order := &models.SysOrders{
		OrderNum:  orderNum,
		UId:       uid,
		SId:       sid,
		PayStatus: payStatus,
	}
	var count int64
	err := models.DB.Model(new(models.SysOrders)).Where("id = ?", id).Count(&count).Error
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
			"msg":  "该订单不存在",
		})
		return
	}
	err = models.DB.Model(new(models.SysOrders)).Where("id = ?", id).Update(order).Error
	if err != nil {
		log.Println("Update Order Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "修改订单信息失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "修改成功",
	})
}

// DeleteOrder
// @Tags 用户私有方法
// @Summary 删除订单
// @Param authorization header string true "authorization"
// @Param id query string true "id"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /user/delete-order [delete]
func DeleteOrder(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "请输入订单号",
		})
		return
	}
	var count int64
	err := models.DB.Model(new(models.SysOrders)).Where("id = ?", id).Count(&count).Error
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
			"msg":  "该订单不存在",
		})
		return
	}
	err = models.DB.Where("id = ?", id).Delete(new(models.SysOrders)).Error
	if err != nil {
		log.Println("Delete Order Error:", err)
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
