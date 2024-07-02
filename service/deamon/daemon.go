package deamon

import (
	"fmt"
	"github.com/kardianos/service"
	"github.com/xxl6097/go-glog/glog"
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
	running bool
	proc    *os.Process
	config  *Config
	svr     *service.Service
}

func newDaemon(config *Config) *daemon {
	return &daemon{
		config: config,
	}
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
	glog.Println("daemon start")
	glog.Println("Status", status, err)
	glog.Println("Platform", s.Platform())
	glog.Println("String", s.String())
	return nil
}

func (d *daemon) Stop(s service.Service) error {
	d.svr = &s
	defer glog.Flush()
	glog.Println("service stop")
	d.running = false
	if d.proc != nil {
		glog.Println("stop worker")
		d.proc.Kill()
	}
	if service.Interactive() {
		glog.Println("stop daemon")
		os.Exit(0)
	}
	return nil
}

func (d *daemon) Status() {
	if d.svr == nil {
		return
	}
	status, err := (*d.svr).Status()
	if err != nil {
		glog.Println(err)
	}
	if status == service.StatusRunning {
		glog.Println("running")
	} else if status == service.StatusRunning {
		glog.Println("stopped")
	} else {
		glog.Println("StatusUnknown", status)
	}
}

func (d *daemon) Run() {
	defer glog.Flush()
	glog.Println("daemon run start")
	defer glog.Println("daemon run end")
	d.running = true
	binPath, _ := os.Executable()
	mydir, err := os.Getwd()
	if err != nil {
		glog.Println(err)
	}
	glog.Println("mydir", mydir)
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
	go s.Run()
	var args []string
	// rm -d parameter
	for i := 0; i < len(os.Args); i++ {
		if os.Args[i] == "-d" {
			args = append(os.Args[0:i], os.Args[i+1:]...)
			break
		}
	}
	args = append(args, "-nv")
	for {
		// start worker
		tmpDump := filepath.Join(mydir+string(filepath.Separator)+"logs", "dump.log.tmp")
		dumpFile := filepath.Join(mydir+string(filepath.Separator)+"logs", "dump.log")
		f, err := os.Create(filepath.Join(tmpDump))
		if err != nil {
			glog.Printf("start worker error:%s", err)
			return
		}
		glog.Println("start worker process, args:", args)
		execSpec := &os.ProcAttr{Env: append(os.Environ(), "GOTRACEBACK=crash"), Files: []*os.File{os.Stdin, os.Stdout, f}}
		p, err := os.StartProcess(binPath, args, execSpec)
		if err != nil {
			glog.Printf("start worker error:%s", err)
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
		glog.Printf("worker stop, restart it after 10s")
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
		fmt.Println(e)
		return e
	}
	d.svr = &s
	e = service.Control(s, ctrlComm)
	if e != nil {
		fmt.Println(e)
		return e
	}

	return nil
}
