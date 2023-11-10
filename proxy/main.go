package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	jsoniter "github.com/json-iterator/go"
)

// Jsoniter 别名
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Config 配置
type Config struct {
	Listen  uint16   // 监听端口
	Forward []string // 转发的目标服务器配置
}

func main() {
	// 参数解析
	conf := flag.String("f", "config.json", "Config file")
	flag.Parse()

	// 创建信息通道 Channel
	psignal := make(chan os.Signal, 1)
	// 枚举可接受的信号，通过signal库写入到信号通道
	// ctrl+c -> SIGINT, kill -9 -> SIGKILL
	signal.Notify(psignal, syscall.SIGINT, syscall.SIGKILL)

	// 配置文件读取
	confBytes, err := ioutil.ReadFile(*conf)
	if err != nil {
		panic(fmt.Sprintf("Config file read fail, err:%+v", err))
	}

	// 配置文件解析
	var config []Config
	err = json.Unmarshal(confBytes, &config)
	if err != nil {
		panic(fmt.Sprintf("Conf Unmarshal fail, err:%+v", err))
	}

	// 执行代理
	go DoServer(config)

	// 通过监听信息通道的方式阻塞程序推出
	<-psignal

	log.Println("Byte~")
}

// DoServer 执行多代理服务处理
func DoServer(configs []Config) {
	for _, config := range configs {
		go handle(config)
	}
}

// handle 单代理服务处理
func handle(config Config) {
	// 负载均衡计数器
	var fid = -1

	// 获取转发地址
	var getForward = func() string {
		// 仅一个转发配置时
		if len(config.Forward) == 1 {
			return config.Forward[0]
		}

		// 存在多个转发配置时进行负载均衡
		fid++
		if fid >= len(config.Forward) {
			fid = 0
		}
		return config.Forward[fid]
	}

	// 链接切换处理
	var doConn = func(conn net.Conn) {
		defer conn.Close()

		forward := getForward()
		log.Println("Dest addr:", forward)

		fconn, err := net.Dial("tcp", forward)
		if err != nil {
			log.Printf("Dial fail, addr[%v] err[%v]\n", forward, err)
			return
		}
		defer fconn.Close()

		// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //
		// 问：为什么这里会使用协程的方式执行`io.Copy`?
		// 答：因为协程可以并发执行，能够提高程序的并发性能。
		//
		// 详解：
		// 		`io.Copy` 函数用于复制数据从一个连接（conn）到另一个连接（fconn）。
		// 如果直接在主程序中执行这个操作，那么在数据传输期间，主程序会阻塞等待数据
		// 传输完成。而使用协程可以将这个操作放在一个单独的goroutine中执行，这样主
		// 程序可以继续处理其他任务，而不需要等待数据传输完成。
		// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //
		go io.Copy(conn, fconn)
		io.Copy(fconn, conn)
	}

	// 启动服务链接
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", config.Listen))
	if err != nil {
		panic(fmt.Sprintf("Start server[%v] fail, err:%v", config.Listen, err))
	}
	defer lis.Close()
	log.Println("Listen on", config.Listen)

	// 持续建立请求，并做请求切换
	for {
		conn, err := lis.Accept()
		if err != nil {
			continue
		}
		go doConn(conn)
	}
}
