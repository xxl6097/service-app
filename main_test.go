package main

import (
	"github.com/xxl6097/go-glog/glog"
	"github.com/xxl6097/service-app/service/util"
	"testing"
)

func Test_Encode_Data(t *testing.T) {
	glog.Println("----- Test_Encode_Data ---")
	url := "http://uuxia.cn:8086/files/2024/07/03/AAServiceTest_0.0.1_windows_amd64.exe"
	a, b := util.GetFileAndExtensionFromURL(url)
	glog.Info("main...", a, b)
}
