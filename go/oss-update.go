package main

//这是一个对比OSS文件并更新的程序
import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

const (
	ossBucket          = "cod***ges"
	ossPoint           = "http://oss-cn-hongkong.aliyuncs.com"
	ossAccessKeyId     = "LTA***toz"
	ossAccessKeySecret = "MIv***JOT"
	originPath         = "/home/wwwroot/l***g/"
)

var ossObject *oss.Client
var ossBucketObj *oss.Bucket

// 异常处理
func HandleError(err error) {
	log.Println("Error:", err)
	os.Exit(-1)
}

// 获取本地文件列表
func localFileList(originpath string) map[string]string {
	var localfile map[string]string
	localfile = make(map[string]string)
	filepath.Walk(originpath,
		func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				//log.Println("dir:", path)
				return nil
			}
			// 忽略.git的目录文件
			if strings.HasPrefix(path, originpath+".git") {
				return nil
			}
			file, inerr := os.Open(path)
			if inerr == nil {
				md5h := md5.New()
				io.Copy(md5h, file)
				md5string := fmt.Sprintf("%x", md5h.Sum([]byte("")))
				localfile[strings.TrimPrefix(path, originpath)] = md5string
			}
			return nil
		})
	return localfile
}

// 获取阿里云OSS文件列表
func ossFileList() map[string]string {
	// 分页列举所有文件。每页列举100个。
	continueToken := ""

	var ossfile map[string]string
	ossfile = make(map[string]string)
	var ossfiletag map[string]int
	ossfiletag = make(map[string]int)
	for {
		lsRes, err := ossBucketObj.ListObjectsV2(oss.MaxKeys(600), oss.ContinuationToken(continueToken))
		if err != nil {
			log.Println(err)
		}
		// 打印列举结果。默认情况下，一次返回100条记录。
		for _, object := range lsRes.Objects {
			ossfile[object.Key] = strings.ToLower(strings.Trim(object.ETag, "\""))
			ossfiletag[object.Key] = 0
		}

		if lsRes.IsTruncated {
			continueToken = lsRes.NextContinuationToken
		} else {
			break
		}
	}
	return ossfile
}

// 更新文件至阿里云OSS
func ossUploadFile(key string, localfile string) {
	// 上传本地文件。
	err := ossBucketObj.PutObjectFromFile(key, localfile)
	if err != nil {
		log.Println("Error:", err)
	}
	log.Println("[更新成功]", key)
}

func main() {
	// 日志记录程序
	logFile, err := os.OpenFile("./oss-update.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}
	log.SetOutput(logFile)
	log.Println("开始执行【" + ossBucket + "】")

	// 判断目录是否存在
	s, err := os.Stat(originPath)
	if err != nil {
		HandleError(err)
	}
	if s.IsDir() != true {
		log.Println("目录不存在")
		return
	}

	// 创建OSSClient实例。
	ossObject, err = oss.New(ossPoint, ossAccessKeyId, ossAccessKeySecret)
	if err != nil {
		HandleError(err)
	}

	ossBucketObj, err = ossObject.Bucket(ossBucket)
	if err != nil {
		HandleError(err)
	}

	lofile := localFileList(originPath)
	osfile := ossFileList()

	// 遍历本地
	for lfile := range lofile {
		// OSS不存在某个文件
		if _, ok := osfile[lfile]; !ok {
			// 上传文件
			log.Println("[上传文件]", lfile)
			ossUploadFile(lfile, originPath+lfile)
		} else if lofile[lfile] != osfile[lfile] {
			// MD5不一致，更新文件
			log.Println("[更新文件]", lfile, lofile[lfile], osfile[lfile])
			ossUploadFile(lfile, originPath+lfile)
		} else {
			// 文件一致，不处理
			// fmt.Println("[文件正常]", lfile)
		}
	}
}
