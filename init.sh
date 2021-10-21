#goctl安装
go get -u github.com/tal-tech/go-zero
go get github.com/tal-tech/go-zero/tools/goctl

#项目go环境
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
go mod download
chmod -R 777 run.sh

#template安装
goctl template init
cp src/util/tpl/handler.tpl ~/.goctl/api/handler.tpl

#编译项目
make gateway