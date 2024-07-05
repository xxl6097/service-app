#!/bin/bash
#修改为自己的应用名称
appname=AAServiceTest
version=0.0.1

#function build_windows_amd641() {
#  #goversioninfo -manifest versioninfo.json
#  rm -rf ${appname}_${version}_windows_amd64.exe
#  go generate
#  #CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-linkmode internal" -o ${appname}_${version}_windows_amd64.exe
#  CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags " -s -w -linkmode internal" -o ${appname}_${version}_windows_amd64.exe
#}

function build_windows_amd64() {
  rm -rf /Volumes/Desktop/service/${appname}_${version}_windows_amd64.exe
  go generate ./cmd/app
  CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags " -s -w -linkmode internal" -o /Volumes/Desktop/service/${appname}_${version}_windows_amd64.exe ./cmd/app
}


function build_windows_arm64() {
  #goversioninfo -manifest versioninfo.json
  #go generate
  CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -trimpath -ldflags "-linkmode internal" -o ${appname}_${version}_windows_arm64.exe
}


function menu() {
  echo -e "\r\n0. 编译 Windows amd64"
  echo "1. 编译 Windows arm64"
  echo "请输入编号:"
  read index
  case "$index" in
  [0]) (build_windows_amd64) ;;
  [1]) (build_windows_arm64) ;;
  *) echo "exit" ;;
  esac

  if ((index >= 4 && index <= 6)); then
    # 获取命令的退出状态码
    exit_status=$?
    # 检查退出状态码
    if [ $exit_status -eq 0 ]; then
      echo "成功推送Docker"
      echo $appversion >version.txt
    else
      echo "失败"
      echo "【$docker_push_result】"
    fi
  fi
}
menu

