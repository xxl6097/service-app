package service

import (
	"fmt"
	"github.com/xxl6097/go-glog/glog"
	"github.com/xxl6097/service-app/service/deamon"
	"math/rand"
	"os"
	"path/filepath"
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
		fmt.Println("1. 安装")
		fmt.Println("2. 卸载")
		fmt.Println("3. 启动")
		fmt.Println("4. 停止")
		fmt.Println("5. 重启")
		fmt.Println("6. 状态")
		fmt.Println("7. 版本")
		fmt.Println("8. 退出")
		fmt.Println("请选择一个选项:")
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

func Run(config *deamon.Config, run func()) {
	if config == nil {
		glog.Fatal("config is nil")
	}
	installPath := defaultInstallPath + string(filepath.Separator) + config.ProductName
	initLog(installPath)
	rand.Seed(time.Now().UnixNano())
	baseDir := filepath.Dir(os.Args[0])
	os.Chdir(baseDir) // for system service
	//fmt.Println("baseDir:", baseDir)
	//fmt.Println("os.Args:", len(os.Args), os.Args)
	glog.Info("Run...", len(os.Args), os.Args)
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
			glog.Println("创建进程..")
			installer.Run()
			glog.Println("进程结束..")
			return
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
	if run != nil {
		run()
	}
}
