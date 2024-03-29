// 又拍云文件列表生成集合目录
// 2023-06-15
// https://github.com/upyun/go-sdk
// package main

package main

import (
    "fmt"
    "github.com/upyun/go-sdk/v3/upyun"
    "strings"
    "time"
    "os"
    "bufio"
    "strconv"
)
var up = upyun.NewUpYun(&upyun.UpYunConfig{
    Bucket:   "xxxxxxxxxx",
    Operator: "xxxxxxxxxx",
    Password: "xxxxxxxxxx",
})

func main() {
    p,_ := UpyunFileList("/")
    if p == nil {
        fmt.Println("执行失败")
    }
    fmt.Println("执行完成")
}

func UpyunFileList(path string) (*[]*upyun.FileInfo,error) {
    fmt.Println(path)
    var result = make([]*upyun.FileInfo, 0)
    var resultpush = make([]*upyun.FileInfo, 0)

    objsChan := make(chan *upyun.FileInfo, 10)
    go func() {
        up.List(&upyun.GetObjectsConfig{
            Path        : path,
            ObjectsChan : objsChan,
        })
    }()
    
    for obj := range objsChan {
        resultpush = append(resultpush, obj)
        if(obj.ContentType == "folder"){
            result0,_ := UpyunFileList(path + obj.Name + "/")
            result1 := *result0
            for obj0 := range result1 {
                result = append(result, result1[obj0])
            }
        }else if obj != nil {
            // obj.Name = path + obj.Name
            result = append(result, obj)
        }
    }
    
    MakePath(path,resultpush)

    p := &result
    return p, nil
}

func MakePath(path string,result []*upyun.FileInfo){
    vpath := strings.Replace(strings.Trim(path,"/"),"/","_",-1)
    rpath := vpath
    lpath := vpath
    if vpath == "" {
        rpath = "default"
    }else{
        lpath = lpath + "_"
    }
    liststr := "<html lang=\"zh-CN\"><head><title>又拍云文件浏览器</title><meta charset=\"utf-8\"><meta name=\"viewport\" content=\"width=device-width,initial-scale=1\"><meta name=\"description\" content=\"Index of /upyun/\">\n<style>body {margin: 0;font-family: \"ubuntu\", \"Tahoma\", \"Microsoft YaHei\", Arial, Serif;}.container {padding-right: 15px;padding-left: 15px;margin-right: auto;margin-left: auto;}@media (min-width: 768px) {.container {max-width: 750px;}}@media (min-width: 992px) {.container {max-width: 970px;}}@media (min-width: 1200px) {.container {max-width: 1170px;}}table {width: 100%;max-width: 100%;margin-bottom: 20px;border: 1px solid #ddd;padding: 0;border-collapse: collapse;}table th {font-size: 14px;}table tr {border: 1px solid #ddd;padding: 5px;}table tr:nth-child(odd) {background: #f9f9f9 }table th, table td {border: 1px solid #ddd;font-size: 14px;line-height: 20px;padding: 3px;text-align: left;}a {color: #337ab7;text-decoration: none;}a:hover, a:focus {color: #2a6496;text-decoration: underline;}table.table-hover>tbody>tr:hover>td, table.table-hover>tbody>tr:hover>th {background-color: #f5f5f5;}.markdown-body {float: left;font-family: \"ubuntu\", \"Tahoma\", \"Microsoft YaHei\", Arial, Serif;}.octicon {background-position: center left;background-repeat: no-repeat;padding-left: 16px;}.file {background-image: url(\"data:image/svg+xml;charset=utf8,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='16' viewBox='0 0 12 16'%3E%3Cpath d='M6 5L2 5 2 4 6 4 6 5 6 5ZM2 8L9 8 9 7 2 7 2 8 2 8ZM2 10L9 10 9 9 2 9 2 10 2 10ZM2 12L9 12 9 11 2 11 2 12 2 12ZM12 4.5L12 14C12 14.6 11.6 15 11 15L1 15C0.5 15 0 14.6 0 14L0 2C0 1.5 0.5 1 1 1L8.5 1 12 4.5 12 4.5ZM11 5L8 2 1 2 1 14 11 14 11 5 11 5Z' fill='%237D94AE'/%3E%3C/svg%3E\");}.file-directory {background-image: url(\"data:image/svg+xml;charset=utf8,%3Csvg xmlns='http://www.w3.org/2000/svg' width='14' height='16' viewBox='0 0 14 16'%3E%3Cpath d='M13 4L7 4 7 3C7 2.3 6.7 2 6 2L1 2C0.5 2 0 2.5 0 3L0 13C0 13.6 0.5 14 1 14L13 14C13.6 14 14 13.6 14 13L14 5C14 4.5 13.6 4 13 4L13 4ZM6 4L1 4 1 3 6 3 6 4 6 4Z' fill='%237D94AE'/%3E%3C/svg%3E\");}.file-zip {background-image: url(\"data:image/svg+xml;charset=utf8,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='16' viewBox='0 0 12 16'%3E%3Cpath d='M8.5 1L1 1C0.4 1 0 1.4 0 2L0 14C0 14.6 0.4 15 1 15L11 15C11.6 15 12 14.6 12 14L12 4.5 8.5 1ZM11 14L1 14 1 2 4 2 4 3 5 3 5 2 8 2 11 5 11 14 11 14ZM5 4L5 3 6 3 6 4 5 4 5 4ZM4 4L5 4 5 5 4 5 4 4 4 4ZM5 6L5 5 6 5 6 6 5 6 5 6ZM4 6L5 6 5 7 4 7 4 6 4 6ZM5 8L5 7 6 7 6 8 5 8 5 8ZM4 9.3C3.4 9.6 3 10.3 3 11L3 12 7 12 7 11C7 9.9 6.1 9 5 9L5 8 4 8 4 9.3 4 9.3ZM6 10L6 11 4 11 4 10 6 10 6 10Z' fill='%237D94AE'/%3E%3C/svg%3E\");}.file-code {background-image: url(\"data:image/svg+xml;charset=utf8,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='16' viewBox='0 0 12 16'%3E%3Cpath d='M8.5,1 L1,1 C0.45,1 0,1.45 0,2 L0,14 C0,14.55 0.45,15 1,15 L11,15 C11.55,15 12,14.55 12,14 L12,4.5 L8.5,1 L8.5,1 Z M11,14 L1,14 L1,2 L8,2 L11,5 L11,14 L11,14 Z M5,6.98 L3.5,8.5 L5,10 L4.5,11 L2,8.5 L4.5,6 L5,6.98 L5,6.98 Z M7.5,6 L10,8.5 L7.5,11 L7,10.02 L8.5,8.5 L7,7 L7.5,6 L7.5,6 Z' fill='%237D94AE' /%3E%3C/svg%3E\");}.file-media {background-image: url(\"data:image/svg+xml;charset=utf8,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='16' viewBox='0 0 12 16'%3E%3Cpath d='M6 5L8 5 8 7 6 7 6 5 6 5ZM12 4.5L12 14C12 14.6 11.6 15 11 15L1 15C0.5 15 0 14.6 0 14L0 2C0 1.5 0.5 1 1 1L8.5 1 12 4.5 12 4.5ZM11 5L8 2 1 2 1 13 4 8 6 12 8 10 11 13 11 5 11 5Z' fill='%237D94AE'/%3E%3C/svg%3E\");}.device-camera-video {background-image: url(\"data:image/svg+xml;charset=utf8,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 16 16'%3E%3Cpath d='M15.2,2.09 L10,5.72 L10,3 C10,2.45 9.55,2 9,2 L1,2 C0.45,2 0,2.45 0,3 L0,12 C0,12.55 0.45,13 1,13 L9,13 C9.55,13 10,12.55 10,12 L10,9.28 L15.2,12.91 C15.53,13.14 16,12.91 16,12.5 L16,2.5 C16,2.09 15.53,1.86 15.2,2.09 L15.2,2.09 Z' fill='%237D94AE' /%3E%3C/svg%3E\");}.octicon-book {padding-left: 20px;background-image: url(\"data:image/svg+xml;charset=utf8,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 16 16'%3E%3Cpath d='M3,5 L7,5 L7,6 L3,6 L3,5 L3,5 Z M3,8 L7,8 L7,7 L3,7 L3,8 L3,8 Z M3,10 L7,10 L7,9 L3,9 L3,10 L3,10 Z M14,5 L10,5 L10,6 L14,6 L14,5 L14,5 Z M14,7 L10,7 L10,8 L14,8 L14,7 L14,7 Z M14,9 L10,9 L10,10 L14,10 L14,9 L14,9 Z M16,3 L16,12 C16,12.55 15.55,13 15,13 L9.5,13 L8.5,14 L7.5,13 L2,13 C1.45,13 1,12.55 1,12 L1,3 C1,2.45 1.45,2 2,2 L7.5,2 L8.5,3 L9.5,2 L15,2 C15.55,2 16,2.45 16,3 L16,3 Z M8,3.5 L7.5,3 L2,3 L2,12 L8,12 L8,3.5 L8,3.5 Z M15,3 L9.5,3 L9,3.5 L9,12 L15,12 L15,3 L15,3 Z' /%3E%3C/svg%3E\");}.arrow-down {font-weight: bold;text-decoration: none !important;background-image: url(\"data:image/svg+xml;charset=utf8,%3Csvg xmlns='http://www.w3.org/2000/svg' width='10' height='16' viewBox='0 0 10 16'%3E%3Cpolygon id='Shape' points='7 7 7 3 3 3 3 7 0 7 5 13 10 7'%3E%3C/polygon%3E%3C/svg%3E\");}.arrow-up {font-weight: bold;text-decoration: none !important;background-image: url(\"data:image/svg+xml;charset=utf8,%3Csvg xmlns='http://www.w3.org/2000/svg' width='10' height='16' viewBox='0 0 10 16'%3E%3Cpolygon id='Shape' points='5 3 0 9 3 9 3 13 7 13 7 9 10 9'%3E%3C/polygon%3E%3C/svg%3E\");}</style></head>\n<body><div class=\"container\"><table><tbody><tr><th><a href=\"default.html\">Index</a> of /<a href=\""
    liststr = liststr + rpath + ".html\">"+strings.Trim(path,"/")+"</a><span style=\"float:right;\">"+time.Now().Format("2006-01-02 15:04:05")+"</span></th></tr></tbody></table><table class=\"table-hover\"><tbody><tr><td><a href=\"javascript:sortby(0)\" class=\"octicon arrow-up\">Name</a></td><td><a href=\"javascript:sortby(1)\" class=\"octicon arrow-up\">Date</a></td><td><a href=\"javascript:sortby(2)\" class=\"octicon arrow-up\">Size</a></td></tr>"

    if vpath != "" {
        // 返回上一级目录(需要计算上一级目录的情况)
        dpath := ""
        arr := strings.Split(strings.Trim(path,"/"),"/")
        if len(arr) >= 2 {
            arr = arr[:len(arr)-1]
            dpath = strings.Join(arr,"_")
        }else{
            dpath = "default"
        }
        liststr = liststr+"<tr><td><a class=\"octicon file-directory\" href=\""+dpath+".html\">..</a></td><td></td><td>"+"0"+"</td></tr>\n"
    }
    dhtml := ""
    fhtml := ""
    for obj := range result {
        if(result[obj].ContentType == "folder"){
            dhtml = dhtml+"<tr><td><a class=\"octicon file-directory\" href=\""+lpath+result[obj].Name+".html\">"+result[obj].Name+"</a></td><td>"+result[obj].Time.Format("2006-01-02 15:04:05")+"</td><td>"+strconv.FormatInt(result[obj].Size,10)+"</td></tr>\n"
        }else{
            fhtml = fhtml+"<tr><td><a class=\"octicon file\" href=\""+path + result[obj].Name+"\">"+result[obj].Name+"</a></td><td>"+result[obj].Time.Format("2006-01-02 15:04:05")+"</td><td>"+strconv.FormatInt(result[obj].Size,10)+"</td></tr>\n"
        }
    }
    liststr = liststr+dhtml+fhtml
    liststr = liststr+"</tbody></table></div></body></html>"
    file, err := os.Create("a/"+rpath+".html")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer file.Close()
    writer := bufio.NewWriter(file)
    _, err = writer.WriteString(liststr)
    if err != nil {
        fmt.Println(err)
        return
    }
    writer.Flush()
}