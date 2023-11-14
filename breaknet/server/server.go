package server

import (
	"breaknet/types"
	"breaknet/utils"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

var (
	mux         sync.Mutex
	g_config    *types.ServerConfig                // 服务端配置
	resourceMap = make(map[uint16]*types.Resource) // 端口-资源对应
)

// Do 主执行入口
func Do(config *types.ServerConfig) {
	g_config = config
	// 启动服务监听端口
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", g_config.Port))
	if err != nil {
		panic(fmt.Errorf("initialization err:%+v", err))
	}
	defer lis.Close()

	for {
		rconn, err := lis.Accept()
		if err != nil {
			log.Printf("Accept err:%+v\n", err)
			continue
		}
		go doConn(rconn, config)
	}
}

// doConn 处理客户端新连接
func doConn(conn net.Conn, config *types.ServerConfig) {
	defer utils.Recover()

	var cmd = make([]byte, 1)
	if _, err := io.ReadAtLeast(conn, cmd, 1); err != nil {
		conn.Close()
		return
	}

	switch cmd[0] {
	case types.START:
		// 客户端服务注册，并启用外部端口监听
		if err := registerClient(conn); err != nil {
			return
		}

		conn.Write([]byte{types.SUCCESS})
		for {
			n, err := conn.Read(cmd)
			if err != nil {
				return
			}
			if n != 0 {
				switch cmd[0] {
				case types.KILL:
					return
				case types.IDLE:
					continue
				}
			}
		}
	case types.NEWCONN:
		// 客户端建立新连接
		sport := make([]byte, 3)
		io.ReadAtLeast(conn, sport, 3)
		port := (uint16(sport[0]) << 8) + uint16(sport[1])
		id := uint8(sport[2])
		client := resourceMap[port]
		if client != nil {
			if int(id) >= len(client.WaitWorker) {
				conn.Close()
				return
			} else {
				wk := client.WaitWorker[id]
				if wk == nil {
					conn.Close()
					return
				}

				if wk.LastTime < time.Now().Unix() {
					if wk.Conn != nil {
						wk.Conn.Close()
					}
					conn.Close()
					return
				}

				var s utils.NCopy
				key, iv := utils.GetKeyIv(g_config.Key)
				s.Init(conn, key, iv)

				go utils.WCopy(&s, wk.Conn) // 写加密，读不加密
				go utils.RCopy(wk.Conn, &s) // 读解密，写不加密

				client.Mux.Lock()
				client.WaitWorker[id] = nil
				client.Mux.Unlock()
			}
		} else {
			conn.Close()
		}
	default:
		conn.Close()
	}
}

// registerClient 客户端服务注册
// Params: conn - 客户端连接句柄
func registerClient(conn net.Conn) error {
	defer conn.Close()

	// 解析消息体长度
	// START info_len info
	info_len := make([]byte, 8)
	if _, err := io.ReadAtLeast(conn, info_len, 8); err != nil {
		return fmt.Errorf("read body info len fail, err:%v", err)
	}

	var ilen = (uint64(info_len[0]) << 56) | (uint64(info_len[1]) << 48) | (uint64(info_len[2]) << 40) |
		(uint64(info_len[3]) << 32) | (uint64(info_len[4]) << 24) | (uint64(info_len[5]) << 16) | (uint64(info_len[6]) << 8) | (uint64(info_len[7]))
	if ilen > 1024*1024 {
		// 限制消息最大内存使用量 1M
		return fmt.Errorf("body over limit")
	}

	// 读取消息结构体
	var cinfo = make([]byte, ilen)
	if _, err := io.ReadAtLeast(conn, cinfo, int(ilen)); err != nil {
		return fmt.Errorf("read body fail, err:%v", err)
	}

	var cconf types.ClientConfig
	if err := json.Unmarshal(cinfo, &cconf); err != nil {
		return fmt.Errorf("unmarshal body fail, err:%v", err)
	}

	if cconf.Key != g_config.Key {
		conn.Write([]byte{types.ERROR_PWD})
		return fmt.Errorf("passwd invalid")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 打开端口
	for _, cc := range cconf.Map {
		// 端口有效性校验
		if len(g_config.LimitPort) > 1 {
			if cc.Outer < g_config.LimitPort[0] || cc.Outer > g_config.LimitPort[1] {
				// 不满足端口范围
				log.Printf("Does not meet the port range[%v, %v] %v", g_config.LimitPort[0], g_config.LimitPort[1], cc.Outer)
				conn.Write([]byte{types.ERROR_LIMIT_PORT})
				return fmt.Errorf("does not meet the port range[%v, %v] %v", g_config.LimitPort[0], g_config.LimitPort[1], cc.Outer)
			}
		}

		// 启动对外端口服务
		clis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", cc.Outer))
		if err != nil {
			log.Println("Port is occupied", cc.Outer)
			conn.Write([]byte{types.ERROR_BUSY})
			return fmt.Errorf("port is occupied, err:%v", err)
		}

		// 收集对外服务的资源到资源集合中
		mux.Lock()
		resourceMap[cc.Outer] = &types.Resource{
			Listener: clis,
			Running:  true,
		}
		mux.Unlock()

		go doListen(ctx, conn, cc.Outer)
	}
	return nil
}

// doListen 监听处理客户端信息
/**************************************************************************************************
 概述：
 	1. 客户端启动服务，会自动连接Server端，进行客户端服务注册；
 	2. 服务端接收到客户端信息后，会先进行信息的校验，通过后会维护端口和客户端来源；
	3. 启动客户端服务暴露的外部端口，并进行消息的监听；
	4. 端口收到信息时，会先找到对应的客户端来源，然后通知客户端与服务端另建立连接；
	5. 客户端收到服务的新建连接时，会根据透传的接口匹配到对应的目的服务；
	6. 客户端分别与服务端、目的服务建立连接，然后将双方句柄互换，即实现了服务端与目的服务的通信；
	7. 服务端再次接收到客户端的连接时，即为与目的服务的连接；
	8. 服务端基于客户端的端口寻址到用户请求，并将用户请求句柄和目的服务端请求句柄进行互换，最终实现服务打通；
**************************************************************************************************/
func doListen(ctx context.Context, inConn net.Conn, port uint16) {
	var rs = resourceMap[port]
	if rs == nil {
		return
	}

	defer func() {
		defer utils.Recover()
		for _, v := range rs.WaitWorker {
			if v != nil && v.Conn != nil {
				v.Conn.Close()
			}
		}
		rs.Listener.Close()
		mux.Lock()
		delete(resourceMap, port)
		mux.Unlock()
		log.Println("Close port:", port)
	}()

	log.Println("Open port:", port)
	go func() {
		defer utils.Recover()

		// 监听访问信息
		for {
			outConn, err := rs.Listener.Accept()
			if err != nil {
				return
			}

			// 通知客户端建立连接
			ok, id := rs.NewConn(inConn)
			if ok {
				var buffer bytes.Buffer
				buffer.Write([]byte{types.NEWSOCKET})
				buffer.Write([]byte{uint8(port >> 8), uint8(port & 0xff)})
				buffer.Write([]byte{id})
				body := buffer.Bytes()
				buffer.Reset()
				inConn.Write(body)
			} else {
				outConn.Close()
			}
		}
	}()
	<-ctx.Done()
}
