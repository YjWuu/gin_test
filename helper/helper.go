package helper

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"math/rand"
	"mime/multipart"
	"persion_test/define"
	"persion_test/models"
	"strconv"
	"time"
)

type UserClaims struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	jwt.StandardClaims
}

var myKey = []byte("person_test")

func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func GetRand() string {
	rand.Seed(time.Now().UnixNano())
	s := ""
	for i := 0; i < 6; i++ {
		s += strconv.Itoa(rand.Intn(10))
	}
	return s
}

func GenerateToken(name string, password string) (string, error) {
	UserClaim := &UserClaims{
		Name:           name,
		Password:       password,
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func AnalyseToken(tokenString string) (*UserClaims, error) {
	userClaims := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userClaims, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, fmt.Errorf("analyse Token Error:%v", err)
	}
	return userClaims, nil
}

func SendCode(phone, code string) error {
	client, _err := func(accessKeyId *string, accessKeySecret *string) (_result *dysmsapi20170525.Client, _err error) {
		config := &openapi.Config{
			// 您的AccessKey ID
			AccessKeyId: accessKeyId,
			// 您的AccessKey Secret
			AccessKeySecret: accessKeySecret,
		}
		// 访问的域名
		config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
		_result = &dysmsapi20170525.Client{}
		_result, _err = dysmsapi20170525.NewClient(config)
		return _result, _err
	}(tea.String("LTAI5tHaKTjoLo8W4GcVnTpU"), tea.String("uc4evQGErrAev9serO8G2YP14Iqmri"))
	if _err != nil {
		fmt.Println(_err)
		return _err
	}
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String(phone),
		SignName:      tea.String("阿里云短信测试"),
		TemplateCode:  tea.String("SMS_154950909"),
		TemplateParam: tea.String(`{"code":"` + code + `"}`),
	}
	// 复制代码运行请自行打印 API 的返回值
	_, _err = client.SendSms(sendSmsRequest)
	return _err
}

func UploadImg(file *multipart.FileHeader) string {
	fileSize := file.Size
	f, _ := file.Open()
	buf := make([]byte, file.Size)
	f.Read(buf)
	putPolicy := storage.PutPolicy{
		Scope: define.Bucket,
	}

	mac := qbox.NewMac(define.AccessKey, define.SecretKey)

	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan, // 这里选择你的空间所在 华南华北华东等
		UseCdnDomains: false,
		UseHTTPS:      false,
	}
	putExtra := storage.PutExtra{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	reader := bytes.NewReader(buf)
	err := formUploader.PutWithoutKey(context.Background(), &ret, upToken, reader, fileSize, &putExtra)
	if err != nil {
		fmt.Println("formUploader.PutWithoutKey err: ", err)
	}

	return ret.Key
}

func CheckImgCode(c *gin.Context, uuid, imgCode string) bool {
	code, err := models.Redis.Get(c, uuid).Result()
	if err != nil {
		fmt.Println("从redis获取imgCode失败:" + err.Error())
	}
	return code == imgCode
}
