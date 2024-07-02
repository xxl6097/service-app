
go get github.com/kardianos/service


CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o AAServiceApp.exe main.go

go get -u github.com/xxl6097/go-glog@v0.0.7