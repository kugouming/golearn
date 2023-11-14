package main

import (
	"breaknet/client"
	"breaknet/server"
	"breaknet/types"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	jsoniter "github.com/json-iterator/go"
)

// Jsoniter 别名
var json = jsoniter.ConfigCompatibleWithStandardLibrary

func main() {
	confile := flag.String("f", "./config/all.json", "Config file")
	flag.Parse()

	psignal := make(chan os.Signal, 1)
	// ctrl+c -> SIGINT kill -9 -> SIGKILL
	signal.Notify(psignal, syscall.SIGINT, syscall.SIGKILL)

	confBytes, err := ioutil.ReadFile(*confile)
	if err != nil {
		panic(fmt.Sprintf("Config file read fail, err:%+v", err))
	}

	var config types.Config
	err = json.Unmarshal(confBytes, &config)
	if err != nil {
		panic(fmt.Sprintf("Conf Unmarshal fail, err:%+v", err))
	}

	// 根据配置内容启动不同的客户端
	if config.Server != nil {
		go server.Do(config.Server)
	}
	if config.Client != nil {
		go client.Do(config.Client)
	}

	<-psignal
	log.Println("Byte~")
}
