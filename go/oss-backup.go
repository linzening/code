package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

const (
	ossBucket          = "cod***ges"
	ossPoint           = "http://oss-cn-hongkong.aliyuncs.com"
	ossAccessKeyId     = "LTA***toz"
	ossAccessKeySecret = "MIv***JOT"
	originPath         = "/data/database/"
)

func main() {
	path := originPath
	file := "docker-go-" + time.Now().Format("20060102150405") + ".sql.gz"
	cmd := exec.Command("/bin/sh", "-c", "/usr/bin/mysqldump --host 127.0.0.1 --port 3306 -uroot -ppass database | gzip > "+path+file)
	err1 := cmd.Run()
	if err1 != nil {
		fmt.Println("error")
		os.Exit(-1)
	}
	fmt.Println("databasefile:" + path + file)

	// 创建OSSClient实例。
	client, err := oss.New(ossPoint, ossAccessKeyId, ossAccessKeySecret)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 获取存储空间。
	bucket, err := client.Bucket(ossBucket)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 上传本地文件。
	err = bucket.PutObjectFromFile("databases2021c/"+file, path+file)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	fmt.Println("Upload Success.")
}
