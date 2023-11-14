package client

import (
	"breaknet/types"
	"breaknet/utils"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

var (
	localPortMap map[uint16]string   // 本地服务与服务端口映射
	config       *types.ClientConfig // 客户端服务配置
)

// Do 主执行入口
func Do(conf *types.ClientConfig) {
	parserConfig(conf)

	for {
		ok := func() bool {
			defer utils.Recover()
			defer time.Sleep(types.RetryTime)
			log.Println("Connecting to server ...")

			// 连接服务端（客户端注册到服务端）
			serConn, err := net.Dial("tcp", config.Server)
			if err != nil {
				log.Println("Cann't connect to server")
				return false
			}
			defer serConn.Close()

			// 服务连接状态判断
			if err := doServerState(serConn, config); err != nil {
				log.Printf("Server connect fail, err:%v\n", err)
				return false
			}

			// 将客户端待映射的服务传递给服务端
			doSendMsg(serConn, config)

			// 读取返回信息，并判断连接状态
			if err := doResCode(serConn, config); err != nil {
				log.Printf("Server connect fail, err:%v\n", err)
				return false
			}

			// 执行NAT操作主逻辑
			if err := doResBody(serConn, config); err != nil {
				log.Printf("Server process fail, err:%v\n", err)
				return false
			}

			return true
		}()
		if !ok {
			break
		}
	}
}

// parserConfig 服务配置解析
func parserConfig(conf *types.ClientConfig) {
	config = conf
	// 客户端注册的服务以端口的形式注册到客户端全局变量中
	localPortMap = make(map[uint16]string, len(config.Map))
	for _, m := range config.Map {
		localPortMap[m.Outer] = m.Inner
	}
}

// doServerState 连接服务端状态判断
func doServerState(serConn net.Conn, config *types.ClientConfig) error {
	defer serConn.Close()

	// 将请求设置为长连接，避免请求链接的频繁创建
	serConn.(*net.TCPConn).SetKeepAlive(true)
	serConn.(*net.TCPConn).SetKeepAlivePeriod(types.TcpKeepAlivePeriod)

	// 将客户端待映射的服务传递给服务端
	var buffer bytes.Buffer // 添加字节缓冲
	// 发送客户端信息
	// START info_len info
	buffer.Write([]byte{types.START})
	cinfo, _ := json.Marshal(config)
	binary.Write(&buffer, binary.BigEndian, uint64(len(cinfo)))
	buffer.Write(cinfo)
	serConn.Write(buffer.Bytes())
	buffer.Reset()

	// 读取返回信息，并判断连接状态
	var cmd = make([]byte, 1)
	io.ReadAtLeast(serConn, cmd, 1)

	// 枚举状态码判断
	switch cmd[0] {
	case types.ERROR_PWD:
		return errors.New("wrong password")
	case types.ERROR_BUSY:
		return errors.New("port is occupied")
	case types.ERROR_LIMIT_PORT:
		return errors.New("does not meet the port range")
	}

	if cmd[0] != types.SUCCESS {
		return errors.New("unkown error")
	}

	// 将请求设置为长连接，避免请求链接的频繁创建
	serConn.(*net.TCPConn).SetKeepAlive(true)
	serConn.(*net.TCPConn).SetKeepAlivePeriod(types.TcpKeepAlivePeriod)

	return nil
}

// doSendMsg 给服务端发送数据
func doSendMsg(serConn net.Conn, config *types.ClientConfig) {
	// 将客户端待映射的服务传递给服务端
	var buffer bytes.Buffer // 添加字节缓冲
	// 发送客户端信息
	// START info_len info
	buffer.Write([]byte{types.START})
	cinfo, _ := json.Marshal(config)
	binary.Write(&buffer, binary.BigEndian, uint64(len(cinfo)))
	buffer.Write(cinfo)
	serConn.Write(buffer.Bytes())
	buffer.Reset()
}

// doResCode 服务端返回Code处理
func doResCode(serConn net.Conn, config *types.ClientConfig) error {
	var cmd = make([]byte, 1)
	io.ReadAtLeast(serConn, cmd, 1)
	switch cmd[0] {
	case types.ERROR_PWD:
		return errors.New("wrong password")
	case types.ERROR_BUSY:
		return errors.New("port is occupied")
	case types.ERROR_LIMIT_PORT:
		return errors.New("does not meet the port range")
	}

	if cmd[0] != types.SUCCESS {
		return errors.New("unkown error")
	}

	log.Println("Certification sucessful")
	for _, cc := range config.Map {
		log.Printf("local server port map:[%v->:%v]\n", cc.Inner, cc.Outer)
	}

	return nil
}

// doResBody 服务端返回内容处理（主逻辑）
func doResBody(serConn net.Conn, config *types.ClientConfig) error {
	var cmd = make([]byte, 1)
	// 进入指令读取循环
	for {
		_, err := serConn.Read(cmd)
		if err != nil {
			return fmt.Errorf("read server body fail, err:%v", err)
		}

		switch cmd[0] {
		// 新建连接
		case types.NEWSOCKET:
			// 读取远端端口与ID
			sp := make([]byte, 3)
			io.ReadAtLeast(serConn, sp, 3)
			sport := uint16(sp[0])<<8 + uint16(sp[1])

			// 重新与Server建立连接（用于NAT桥接）
			conn, err := net.Dial("tcp", config.Server)
			if err != nil {
				return fmt.Errorf("connect server fail, err:%v", err)
			}
			go doConn(conn, sport, sp)
		case types.IDLE:
			fallthrough
		default:
			_, err := serConn.Write([]byte{types.SUCCESS})
			if err != nil {
				return fmt.Errorf("write res succ fail, err:%v", err)
			}
		}
	}

	return nil
}

// doConn 执行来访服务与目的服务连接桥接
func doConn(fromConn net.Conn, port uint16, pbytes []byte) {
	defer utils.Recover()

	// 连接至目的服务
	destConn, err := net.Dial("tcp", localPortMap[port])
	if err != nil {
		fromConn.Close()
		log.Println(err)
		return
	}

	// 将服务连接状态以及服务端口写入返回数据包
	fromConn.Write([]byte{types.NEWCONN})
	fromConn.Write(pbytes)

	// 计算出通信密钥
	key, iv := utils.GetKeyIv(config.Key)

	// 将密钥进行加密并写入返回数据包
	var s utils.NCopy
	s.Init(fromConn, key, iv)

	// 通信句柄转换（NAT）
	go utils.WCopy(&s, destConn) // 写加密，读不加密
	go utils.RCopy(destConn, &s) // 读解密，写不加密
}
