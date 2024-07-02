package main

import (
	"service-app/service"
	"service-app/service/deamon"
)

func main() {
	service.Run(&deamon.Config{
		ProductName: "AAServiceTest",
		DisplayName: "A Test Service",
		Description: "A Golang Service..",
		AppVersion:  "0.0.1",
	})
}
