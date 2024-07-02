package service

import (
	"fmt"
	"github.com/xxl6097/go-glog/glog"
	"math/rand"
	"os"
	"path/filepath"
	"service-app/service/deamon"
	"time"
)

func initLog(installPath string) {
	glog.SetLogFile(installPath+string(filepath.Separator)+"logs", "app.log")
	glog.SetMaxSize(1 * 1024 * 1024)
	glog.SetMaxAge(15)
	glog.SetCons(true)
	glog.SetNoHeader(true)
	glog.SetNoColor(true)
}

func Menu(config *deamon.Config, installer *deamon.Installer) {
	var choice int
	for {
		glog.Println("\r\n1. 安装")
		glog.Println("2. 卸载")
		glog.Println("3. 启动")
		glog.Println("4. 停止")
		glog.Println("5. 重启")
		glog.Println("6. 状态")
		glog.Println("7. 版本")
		glog.Println("8. 退出")
		glog.Println("请选择一个选项:")
		fmt.Scan(&choice)
		switch choice {
		case 1:
			installer.Install()
			break
		case 2:
			installer.Uninstall()
			break
		case 3:
			installer.Start()
			break
		case 4:
			installer.Stop()
			break
		case 5:
			installer.Restart()
			break
		case 6:
			installer.Status()
			break
		case 7:
			glog.Println(config.AppVersion)
			break
		case 8:
			glog.Println("退出程序")
			os.Exit(0)
			return // 退出 main 函数，结束程序
		default:
			glog.Println("无效的选项，请重新输入")
		}
	}
}

func Run(config *deamon.Config) {
	if config == nil {
		glog.Fatal("config is nil")
	}
	installPath := defaultInstallPath + string(filepath.Separator) + config.ProductName
	initLog(installPath)
	rand.Seed(time.Now().UnixNano())
	baseDir := filepath.Dir(os.Args[0])
	os.Chdir(baseDir) // for system service
	//glog.Println("baseDir:", baseDir)
	//glog.Println("os.Args:", len(os.Args), os.Args)
	installer := deamon.NewInstall(config, installPath)
	if len(os.Args) > 1 {
		glog.SetNoHeader(false)
		glog.SetNoColor(false)
		switch os.Args[1] {
		case "version", "-v", "--version":
			glog.Println(config.AppVersion)
			return
		case "install":
			installer.Install()
			return
		case "uninstall":
			installer.Uninstall()
			return
		case "start":
			installer.Start()
		case "stop":
			installer.Stop()
		case "restart":
			installer.Restart()
			return
		case "-d":
			glog.Println("-d daemon run...")
			installer.Run()
			return
			//case "-nv":
			//	glog.Println("-nv daemon run...")
			//	installer.Run()
			//	return
		}
	} else {
		//installer.InstallByFilename()
		Menu(config, installer)
	}
	SetFirewall(config.ProductName)
	err := SetRLimit()
	if err != nil {
		glog.Println("setRLimit error:", err)
	}
	glog.Println("exe start...")
	installer.Status()
	//glog.Println(time.Now().Format("exe start....2016-01-02 15:04:05"))
	//glog.Flush()
	//forever := make(chan bool)
	//<-forever
}