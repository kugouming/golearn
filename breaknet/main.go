package main

import (
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

// ServerConfig 服务端配置
type ServerConfig struct {
	Key       string   `json:"key"`        // 请求密钥，校验合法性
	Port      uint16   `json:"port"`       // 监听端口
	LimitPort []uint16 `json:"limit_port"` // 开发端口列表
}

// ClientMapConfig 客户端Map配置
type ClientMapConfig struct {
	Inner string `json:"inner"`
	Outer uint16 `json:"outer"`
}

// ClientConfig 客户端配置
type ClientConfig struct {
	Key    string            `json:"key"`    // 请求密钥，校验合法性
	Server string            `json:"server"` // Server链接配置
	Map    []ClientMapConfig `json:"map"`    // 本地需要映射的服务配置
}

// Config 服务配置
type Config struct {
	Server *ServerConfig `json:"server"`
	Client *ClientConfig `json:"client"`
}

func main() {
	confile := flag.String("f", "config.json", "Config file")
	flag.Parse()

	psignal := make(chan os.Signal, 1)
	// ctrl+c -> SIGINT kill -9 -> SIGKILL
	signal.Notify(psignal, syscall.SIGINT, syscall.SIGKILL)

	confBytes, err := ioutil.ReadFile(*confile)
	if err != nil {
		panic(fmt.Sprintf("Config file read fail, err:%+v", err))
	}

	var config Config
	err = json.Unmarshal(confBytes, &config)
	if err != nil {
		panic(fmt.Sprintf("Conf Unmarshal fail, err:%+v", err))
	}

	// 根据配置内容启动不同的客户端
	if config.Server != nil {
		go DoServer(config.Server)
	}
	if config.Client != nil {
		go DoClient(config.Client)
	}

	<-psignal
	log.Println("Byte~")
}
