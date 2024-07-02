package deamon

import (
	"fmt"
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
	config      *Config
	demoner     *daemon
	installPath string
}

func NewInstall(_config *Config, _installPath string) *Installer {
	return &Installer{
		installPath: _installPath,
		config:      _config,
		demoner:     newDaemon(_config),
	}
}
func (this *Installer) Install() {
	defer glog.Flush()
	defer glog.Println("install end")
	//installPath := this.installPath + string(filepath.Separator) + this.config.ProductName
	defaultBinName := this.config.ProductName
	if util.IsWindows() {
		defaultBinName += ".exe"
	}
	glog.Println("install installPath", this.installPath)
	// auto uninstall
	err := os.MkdirAll(this.installPath, 0775)
	if err != nil {
		glog.Printf("MkdirAll %s error:%s", this.installPath, err)
		return
	}
	err = os.Chdir(this.installPath)
	if err != nil {
		glog.Println("cd error:", err)
		return
	}

	this.Uninstall()

	targetPath := filepath.Join(this.installPath, defaultBinName)

	binPath, err1 := os.Executable()
	if err1 != nil {
		glog.Fatal("os.Executable() error", err1)
		return
	}
	glog.Println("binPath", binPath)
	src, errFiles := os.Open(binPath) // can not use args[0], on Windows call openp2p is ok(=openp2p.exe)
	if errFiles != nil {
		glog.Printf("os.OpenFile %s error:%s", os.Args[0], errFiles)
		return
	}
	dst, errFiles := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0775)
	if errFiles != nil {
		glog.Printf("os.OpenFile %s error:%s", targetPath, errFiles)
		return
	}

	_, errFiles = io.Copy(dst, src)
	if errFiles != nil {
		glog.Printf("io.Copy error:%s", errFiles)
		return
	}
	src.Close()
	dst.Close()
	// install system service
	glog.Println("targetPath:", targetPath)
	err = this.demoner.Control("install", targetPath, []string{"-d"})
	if err == nil {
		glog.Println("install system service ok.")
	} else {
		glog.Println("install system service error:", err)
	}
	time.Sleep(time.Second * 2)
	err = this.demoner.Control("start", targetPath, []string{"-d"})
	if err != nil {
		glog.Println("start service error:", err)
	} else {
		glog.Println("start service ok.")
	}
}

func (this *Installer) Uninstall() {
	defer glog.Flush()
	//defaultInstallPath := config.InstallPath
	defaultBinName := this.config.ProductName
	if util.IsWindows() {
		defaultBinName += ".exe"
	}
	glog.Println("uninstall start")
	defer glog.Println("uninstall end")
	err := this.demoner.Control("stop", "", nil)
	if err != nil { // service maybe not install
		return
	}
	err = this.demoner.Control("uninstall", "", nil)
	if err != nil {
		glog.Println("uninstall system service error:", err)
	} else {
		glog.Println("uninstall system service ok.")
	}
	glog.Println("uninstall installPath", this.installPath)
	binPath := filepath.Join(this.installPath, defaultBinName)
	os.Remove(binPath + "0")
	os.Remove(binPath)
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
	glog.Println("restart")
	defer glog.Println("restart end")
	err := this.demoner.Control("restart", "", nil)
	if err != nil {
		glog.Println("restart system service error:", err)
	} else {
		glog.Println("restart system service ok.")
	}
}

func (this *Installer) Start() {
	defer glog.Flush()
	glog.Println("start")
	defer glog.Println("start end")
	err := this.demoner.Control("start", "", nil)
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
	err := this.demoner.Control("stop", "", nil)
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
