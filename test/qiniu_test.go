package test

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"os"
	"testing"
	"time"
)

const (
	accessKey = "20xexftuPIXXLuX0K5bwVX9cKtz5DqmBDAbV5Yt8"
	secretKey = "wEQluEzyfODMMEchKrBjD73gwMpFyc2b2tQb61tZ"
	bucket    = "xiaocaicai"
)

func TestUpload(t *testing.T) {
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}

	mac := qbox.NewMac(accessKey, secretKey)

	localFile := "C:/Users/admin/Pictures/Saved Pictures/1.jpg"
	key := "2-x.png"

	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	resumeUploader := storage.NewResumeUploaderV2(&cfg)
	ret := storage.PutRet{}
	recorder, err := storage.NewFileRecorder(os.TempDir())
	if err != nil {
		fmt.Println(err)
		return
	}
	putExtra := storage.RputV2Extra{
		Recorder: recorder,
	}
	err = resumeUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ret.Key, ret.Hash)
}

func TestDownLoad(t *testing.T) {
	mac := qbox.NewMac(accessKey, secretKey)
	domain := "http://rao4os4ip.hn-bkt.clouddn.com"
	key := "1-x.png"

	deadline := time.Now().Add(time.Second * 3600).Unix() //1小时有效期
	privateAccessURL := storage.MakePrivateURL(mac, domain, key, deadline)
	fmt.Println(privateAccessURL)
}
