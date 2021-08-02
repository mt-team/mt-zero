#mt-zero微服务工程框架

##项目环境
```
go1.15
go-zero
gomod
gorm
goimoprts
etcd
mysql
redis
GNU Make 4.3
goreman 0.2.1
```


##项目目录
```
src 项目
    gateway  网关层
    util    依赖包
bin 编译后执行文件，gitignore
go.mod      gomod文件
Makefile    make命令文件
Procfile    goreman命令文件
run.sh      本地执行项目文件
docker      容器化镜像文件
```

##初始化项目
```
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
go mod download
chmod -R 777 run.sh
response模板调整
```

##response模板调整
```
go-zero的响应模板调整
goctl template init
cp src/util/response/handler.tpl ~/.goctl/api/handler.tpl
```

##make命令
```
make gateway 解析api文件并编译gateway
make {project} 解析proto文件并编译项目，例如meeting
make newrpc PROJECT={project} 创建一个新微服务并生成proto，例如app
make rpc PROJECT={project} 解析指定服务微服务proto，例如app
make build PROJECT={project} 编译指定服务微服务，例如app
make model PROJECT={project} 预生成mysql表映射代码，前提是在微服务目录下创建sql/{mysql}.sql文件
```

##goreman命令
```
多项目启动软件
goreman start 启动Procfile中的服务
goreman run status 启动进程信息
```

##gateway编写规则
```
所有response如果有err，必须使用response.ErrXXXX.WithMessage("error msg")
用来规范对外输出的code和错误信息
```