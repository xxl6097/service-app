
go get github.com/kardianos/service


CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build-trimpath -ldflags "-linkmode internal" -o AAServiceApp.exe main.go

go get -u github.com/xxl6097/go-glog@v0.0.7

goversioninfo -manifest versioninfo.json

## windows打包

1、 main.go 文件中添加标签，如下

```
//go:generate goversioninfo -icon=resource/icon.ico -manifest=resource/goversioninfo.exe.manifest
func main() {
}
```

2、编译打包

```
go generate
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-linkmode internal $ldflags" -o ${appname}_${version}_windows_amd64.exe
```

3、生成版本信息
resource文件夹