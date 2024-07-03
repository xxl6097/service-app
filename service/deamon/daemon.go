package deamon

import (
	"fmt"
	"github.com/kardianos/service"
	"github.com/xxl6097/go-glog/glog"
	"github.com/xxl6097/service-app/service/util"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	AppVersion  string
	ProductName string
	DisplayName string
	Description string
}

type daemon struct {
	running    bool
	proc       *os.Process
	config     *Config
	svr        *service.Service
	upgradeUrl string
}

func newDaemon(config *Config) *daemon {
	this := &daemon{
		config:     config,
		upgradeUrl: "",
	}
	return this
}

func (d *daemon) GetService() *service.Service {
	return d.svr
}

func (d *daemon) Shutdown(s service.Service) error {
	status, err := s.Status()
	glog.Println("daemon Shutdown")
	glog.Println("Status", status, err)
	glog.Println("Platform", s.Platform())
	glog.Println("String", s.String())
	return nil
}

func (d *daemon) Start(s service.Service) error {
	d.svr = &s
	status, err := s.Status()
	glog.Println("启动服务")
	glog.Println("Status", status, err)
	glog.Println("Platform", s.Platform())
	glog.Println("String", s.String())
	return nil
}

func (d *daemon) Stop(s service.Service) error {
	d.svr = &s
	defer glog.Flush()
	glog.Println("停止服务")
	d.running = false
	if d.proc != nil {
		glog.Println("停止worker进程")
		d.proc.Kill()
	}
	if service.Interactive() {
		glog.Println("停止deamon")
		os.Exit(0)
	}
	return nil
}

func (d *daemon) IsRunning() bool {
	if d.svr == nil {
		return false
	}
	status, err := (*d.svr).Status()
	if err != nil {
		glog.Println(err)
		return false
	}
	//glog.Println("status", status)
	if status == service.StatusRunning {
		glog.Println(d.config.ProductName, "is running")
		return true
	} else if status == service.StatusStopped {
		glog.Println(d.config.ProductName, "is stopped")
	} else {
		glog.Println(d.config.ProductName, "StatusUnknown", status)
	}
	return false
}

func (d *daemon) Status() {
	d.IsRunning()
}

// 定义一个简单的处理函数
func (d *daemon) helloHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	glog.Info(url)
	if url != "" {
		d.upgradeUrl = url
		if d.proc != nil {
			glog.Println("停止worker进程")
			d.proc.Kill()
		}
		fmt.Fprintf(w, "Hello, World!")
	} else {
		fmt.Fprintf(w, "no url")
	}

}

func (d *daemon) upgrade() {
	http.HandleFunc("/api", util.BasicAuth(d.helloHandler, "admin", "het002402"))
	// 启动HTTP服务器，监听8080端口
	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}

func (d *daemon) download(exeDir string) string {
	defer func() {
		d.upgradeUrl = ""
	}()
	if d.upgradeUrl != "" {
		filePath, err := util.DownloadFile(exeDir, d.upgradeUrl)
		glog.Info(filePath)
		if err != nil {
			return ""
		} else {
			return filePath
		}
	}
	return ""
}

func (d *daemon) Run() {
	defer glog.Flush()
	glog.Println("新建进程启动程序")
	defer glog.Println("新建进程结束")
	d.running = true
	binPath, _ := os.Executable()
	mydir, err := os.Getwd()
	if err != nil {
		glog.Println(err)
	}
	//glog.Println("mydir", mydir)
	glog.Println("binPath", binPath)
	conf := &service.Config{
		Name:        d.config.ProductName,
		DisplayName: d.config.DisplayName,
		Description: d.config.Description,
		Executable:  binPath,
	}

	s, err1 := service.New(d, conf)
	if err1 != nil {
		glog.Println(err1)
		return
	}
	d.svr = &s
	go d.upgrade()
	go s.Run()
	var args []string
	//删除-d参数
	for i := 0; i < len(os.Args); i++ {
		if os.Args[i] == "-d" {
			args = append(os.Args[0:i], os.Args[i+1:]...)
			break
		}
	}
	args = append(args, "-nv")
	for {
		newfile := d.download(mydir)
		if newfile != "" {
			binPath = newfile
		}
		// start worker
		tmpDump := filepath.Join(mydir+string(filepath.Separator)+"logs", "dump.log.tmp")
		dumpFile := filepath.Join(mydir+string(filepath.Separator)+"logs", "dump.log")
		f, err2 := os.Create(filepath.Join(tmpDump))
		if err2 != nil {
			glog.Printf("start worker error:%s", err2)
			return
		}
		glog.Println("启动worker进程，参数：", args)
		execSpec := &os.ProcAttr{Env: append(os.Environ(), "GOTRACEBACK=crash"), Files: []*os.File{os.Stdin, os.Stdout, f}}
		p, err3 := os.StartProcess(binPath, args, execSpec)
		if err3 != nil {
			glog.Printf("启动worker进程失败，错误信息：%s", err3)
			return
		}
		d.proc = p
		_, _ = p.Wait()
		f.Close()
		time.Sleep(time.Second)
		err = os.Rename(tmpDump, dumpFile)
		if err != nil {
			glog.Printf("rename dump error:%s", err)
		}
		if !d.running {
			return
		}
		glog.Printf("worker进程停止,10秒后重新启动")
		time.Sleep(time.Second * 10)
	}
}

func (d *daemon) Control(ctrlComm string, exeAbsPath string, args []string) error {
	svcConfig := &service.Config{
		Name:        d.config.ProductName,
		DisplayName: d.config.DisplayName,
		Description: d.config.Description,
		Executable:  exeAbsPath,
		Arguments:   args,
	}

	//glog.Debugf("%s %+v", ctrlComm, svcConfig)
	s, e := service.New(d, svcConfig)
	if e != nil {
		glog.Println("New", e)
		return e
	}
	//status := s.Status()
	d.svr = &s
	e = service.Control(s, ctrlComm)
	if e != nil {
		glog.Println("Control", e)
		return e
	}

	return nil
}
