package main

import (
	"fmt"
	"github.com/xxl6097/go-glog/glog"
	"github.com/xxl6097/service-app/service"
	"github.com/xxl6097/service-app/service/deamon"
	"time"
)

//go:generate goversioninfo -icon=resource/icon.ico -manifest=resource/goversioninfo.exe.manifest
func main() {
	service.Run(&deamon.Config{
		ProductName: "AAServiceTest",
		DisplayName: "A Test Service",
		Description: "A Golang Service..",
		AppVersion:  "0.0.1",
	}, func() {
		// 初始化服务
		fmt.Println("初始化服务...")
		for {
			glog.Error("aaa", time.Now().Format("2006-01-02 15:04:05"))
			time.Sleep(time.Second * 10)
		}
	})
}
