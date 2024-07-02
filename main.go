package main

import (
	"github.com/xxl6097/service-app/service"
	"github.com/xxl6097/service-app/service/deamon"
)

func main() {
	service.Run(&deamon.Config{
		ProductName: "AAServiceTest",
		DisplayName: "A Test Service",
		Description: "A Golang Service..",
		AppVersion:  "0.0.1",
	}, func() {
		// 初始化服务
	})
}
