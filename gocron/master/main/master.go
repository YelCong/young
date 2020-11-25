package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"
	"young/gocron/master"
)

var (
	configFile string //配置文件路径
)

//解析命令行参数
func initArgs() {
	//master -config ./master.json -x 123 -y 456
	//master -h
	flag.StringVar(&configFile, "config", "./master.json", "指定配置文件master.json")
	flag.Parse()
}

//初始化线程数量
func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		err error
	)

	//初始化命令行参数
	initArgs()

	//初始化线程数
	initEnv()

	//加载配置
	if err = master.InitConfig(configFile); err != nil {
		goto ERR
	}

	//任务管理器
	if err = master.InitJonMgr(); err != nil {
		goto ERR
	}

	//启动Api Http服务
	if err = master.InitApiServer(); err != nil {
		goto ERR
	}

	for {
		time.Sleep(1 * time.Second)
	}

	return

ERR:
	fmt.Println(err)

}
