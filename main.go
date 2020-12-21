package main

import (
	"context"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

/*
	指定一个图片文件作为参数，上传至minio服务器上。
	允许图片的格式： .jpg、 .png .jpeg


 */

//var target = "192.168.1.104"
var target = "pic.llsong.xyz"

func main() {
	// 参数检查
	if len(os.Args) <2{
		log.Println("please input the filepath");pause()
		return
	}

	_, err := os.Stat(os.Args[1])
	if err!=nil {
		log.Println(err);pause()
		return
	}

	filename := filepath.Base(os.Args[1])
	ext := filepath.Ext(filename)
	filename = strings.TrimSuffix(filename,ext)
	contentType := ""
	switch ext {
	case ".jpg":
		fallthrough
	case ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	default:
		log.Fatalln("图片格式不支持");pause()
		return

	}

	// 连接minio
	endpoint := target + ":9000"
	accessKeyId := "minio"
	secretAccessKey := "linyu168@"
	useSSL := false
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyId, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Println("connect to minio server failed",endpoint);pause()
		return
	}

	ctx := context.Background()
	bucketName := "picbed"
	exists, err := client.BucketExists(ctx, bucketName)
	if err!=nil || !exists {
		log.Printf("%s is not exist",bucketName);pause()
		return
	}

	// 上传文件
	now := time.Now()
	objectName := fmt.Sprintf("%d/%d/%d/%s-%d%d%d%s",now.Year(),now.Month(),now.Day(),filename,now.Hour(),now.Minute(),now.Second(),ext)
	info, err := client.FPutObject(ctx,bucketName, objectName, os.Args[1], minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		log.Printf("upload fail with error", err);pause()
		return
	}
	log.Println(info)
	result :=  fmt.Sprintf("http://pic.llsong.xyz:9000/%s/%s",bucketName,objectName)
	fmt.Println(result)

	// 拷贝到剪贴板
	err = clipboard.WriteAll(result)
	if err != nil {
		fmt.Println(err)
	}

}

func pause(){
	var i int =0
	fmt.Scanln(&i)
}
