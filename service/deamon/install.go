package deamon

import (
	"fmt"
	"github.com/kardianos/service"
	"github.com/xxl6097/go-glog/glog"
	"github.com/xxl6097/service-app/service/util"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Installer struct {
	demoner *daemon
	binDir  string
	binName string
	binPath string
}

func NewInstall(_config *service.Config, _installPath string) *Installer {
	this := &Installer{
		binDir: _installPath,
	}
	this.binName = _config.Name
	if util.IsWindows() {
		this.binName += ".exe"
	}
	this.binPath = filepath.Join(this.binDir, this.binName)
	_config.Executable = this.binPath
	_config.Arguments = []string{"-d"}
	this.demoner = newDaemon(_config)
	return this
}
func (this *Installer) Install() {
	defer glog.Flush()
	defer glog.Println("安装结束")
	glog.Println("安装路径：", this.binDir)
	err := os.MkdirAll(this.binDir, 0775)
	if err != nil {
		glog.Printf("MkdirAll %s error:%s", this.binDir, err)
		return
	}
	err = os.Chdir(this.binDir)
	if err != nil {
		glog.Println("cd error:", err)
		return
	}

	this.Uninstall()

	binPath, err1 := os.Executable()
	if err1 != nil {
		glog.Fatal("os.Executable() error", err1)
		return
	}
	glog.Println("可执行程序位置：", binPath)
	src, errFiles := os.Open(binPath) // can not use args[0], on Windows call openp2p is ok(=openp2p.exe)
	if errFiles != nil {
		glog.Printf("os.OpenFile %s error:%s", os.Args[0], errFiles)
		return
	}
	dst, errFiles := os.OpenFile(this.binPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0775)
	if errFiles != nil {
		glog.Printf("os.OpenFile %s error:%s", this.binPath, errFiles)
		return
	}

	_, errFiles = io.Copy(dst, src)
	if errFiles != nil {
		glog.Printf("文件拷贝失败，错误信息：%s", errFiles)
		return
	}
	src.Close()
	dst.Close()
	// install system service
	glog.Println("程序位置:", this.binPath)
	err = this.demoner.Controls("install") //.Control("install", this.binPath, []string{"-d"})
	if err == nil {
		glog.Println("服务安装成功!")
	} else {
		glog.Println("服务安装失败，错误信息:", err)
	}
	time.Sleep(time.Second * 2)
	err = this.demoner.Controls("start") //Control("start", this.binPath, []string{"-d"})
	if err != nil {
		glog.Println("服务启动失败，错误信息:", err)
	} else {
		glog.Println("服务启动成功！")
	}
}

func (this *Installer) Uninstall() {
	defer glog.Flush()
	defer glog.Println("卸载结束")
	glog.Println("开始卸载程序")
	if this.demoner.IsRunning() {
		err := this.demoner.Controls("stop") //.Control("stop", "", nil)
		if err != nil {                      // service maybe not install
			glog.Println("卸载失败，错误信息：", err)
			return
		}
	} else {
		glog.Println("服务未运行")
	}

	err := this.demoner.Controls("uninstall") //Control("uninstall", "", nil)
	if err != nil {
		glog.Println("服务卸载失败，错误信息：", err)
	} else {
		glog.Println("服务成功卸载！")
	}
	glog.Println("卸载程序路径", this.binDir)
	os.Remove(this.binPath + "0")
	os.Remove(this.binPath)
}

func (this *Installer) InstallByFilename() {
	defer glog.Flush()
	glog.Println("installByFilename", os.Args[0])
	params := strings.Split(filepath.Base(os.Args[0]), "-")
	if len(params) < 4 {
		return
	}
	glog.Println("params", params)
	serverHost := params[1]
	token := params[2]
	glog.Println("install start")
	targetPath := os.Args[0]
	args := []string{"install"}
	args = append(args, "-serverhost")
	args = append(args, serverHost)
	args = append(args, "-token")
	args = append(args, token)
	env := os.Environ()
	cmd := exec.Command(targetPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = env
	err := cmd.Run()
	if err != nil {
		glog.Println("install by filename, start process error:", err)
		return
	}
	glog.Println("install end")
	glog.Println("Press the Any Key to exit")
	fmt.Scanln()
	os.Exit(0)
}

func (this *Installer) Restart() {
	defer glog.Flush()
	defer glog.Println("restart end")
	glog.Println("重启...")
	err := this.demoner.Controls("restart") //Control("restart", "", nil)
	if err != nil {
		glog.Println("服务重启失败，错误信息：", err)
	} else {
		glog.Println("服务重启成功!")
	}
}

func (this *Installer) Start() {
	defer glog.Flush()
	glog.Println("start")
	defer glog.Println("start end")
	err := this.demoner.Controls("start") //Control("start", "", nil)
	if err != nil {
		glog.Println("start system service error:", err)
	} else {
		glog.Println("start system service ok.")
	}
}
func (this *Installer) Stop() {
	defer glog.Flush()
	glog.Println("stop")
	defer glog.Println("stop end")
	err := this.demoner.Controls("stop") //.Control("stop", "", nil)
	if err != nil {
		glog.Println("stop system service error:", err)
	} else {
		glog.Println("stop system service ok.")
	}
}

func (this *Installer) Run() {
	this.demoner.Run()
}
func (this *Installer) Status() {
	this.demoner.Status()
}
