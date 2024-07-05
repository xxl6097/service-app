package main

import (
	"fmt"
	service2 "github.com/kardianos/service"
	"github.com/xxl6097/go-glog/glog"
	"github.com/xxl6097/service-app/service"
	"time"
)

//go:generate goversioninfo -icon=resource/icon.ico -manifest=resource/goversioninfo.exe.manifest
func main() {
	//url := "http://uuxia.cn:8086/files/2024/07/03/AAServiceTest_0.0.1_windows_amd64.exe"
	//a, b := util.GetFileAndExtensionFromURL(url)
	//fmt.Println("main...", a, b)
	//service.Run(&deamon.Config{
	//	ProductName: "AAServiceTest",
	//	DisplayName: "A Test Service",
	//	Description: "A Golang Service..",
	//	AppVersion:  "0.0.1",
	//}, func() {
	//	// 初始化服务
	//	fmt.Println("初始化服务...")
	//	for {
	//		glog.Error("uuuu----aaa", time.Now().Format("2006-01-02 15:04:05"))
	//		time.Sleep(time.Second * 120)
	//	}
	//})
	service.Run(&service2.Config{
		Name:        "AAServiceTest",
		DisplayName: "A Test Service",
		Description: "A Golang Service..",
	}, func() {
		// 初始化服务
		fmt.Println("初始化服务...")
		for {
			glog.Error("uuuu----aaa", time.Now().Format("2006-01-02 15:04:05"))
			time.Sleep(time.Second * 120)
		}
	})
}
