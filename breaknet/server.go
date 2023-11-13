package main

import (
	"breaknet/encrypto"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

const (
	_ uint8 = iota
	// START 第一次连接服务器
	START
	// NEWSOCKET 新连接
	NEWSOCKET
	// NEWCONN 新连接发送到服务端命令
	NEWCONN
	// ERROR 处理失败
	ERROR
	// SUCCESS 处理成功
	SUCCESS
	// IDLE 空闲命令 什么也不做
	IDLE
	// KILL 退出命令
	KILL
	// ERROR_PWD 密码错误
	ERROR_PWD
	// ERROR_BUSY 端口被占用
	ERROR_BUSY
	// ERROR_LIMIT_PORT 不满足端口范围
	ERROR_LIMIT_PORT
)

const (
	// RetryTime 断线重连时间
	RetryTime          = time.Second
	TcpKeepAlivePeriod = 30 * time.Second
	WaitTimeOut        = 30 * time.Second // 连接等待超时时间
	WaitMax            = 10
)

type Worker struct {
	Conn     net.Conn // 客户端连接
	LastTime int64    // 客户端连接超时时间
}

type Resource struct {
	Listener   net.Listener
	WaitWorker [WaitMax]*Worker // 工作负载
	Running    bool
	mu         sync.Mutex // 工作负载锁
}

// NewConn 建立新连接
func (r *Resource) NewConn(conn net.Conn) (bool, uint8) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, v := range r.WaitWorker {
		if v == nil {
			r.WaitWorker[i] = &Worker{
				Conn:     conn,
				LastTime: time.Now().Add(WaitTimeOut).Unix(),
			}

			return true, uint8(i)
		}

		// 超时
		if time.Now().Unix() > v.LastTime {
			v.Conn.Close()
			r.WaitWorker[i] = &Worker{
				Conn:     conn,
				LastTime: time.Now().Add(WaitTimeOut).Unix(),
			}
			return true, uint8(i)
		}
	}
	return false, 0
}

func Recover() {
	if err := recover(); err != nil {
		log.Printf("Recover err:%+v", err)
	}
}

// 端口-资源对应
var resourceMap = make(map[uint16]*Resource)
var resourceMux sync.Mutex

func DoServer(config *ServerConfig) {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", config.Port))
	if err != nil {
		log.Printf("Initialization err:%+v\n", err)
		return
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

// doListen 处理客户端监听
func doListen(ctx context.Context, inConn net.Conn, port uint16) {
	var rs = resourceMap[port]
	if rs == nil {
		return
	}

	defer func() {
		Recover()
		defer Recover()
		for _, v := range rs.WaitWorker {
			if v != nil && v.Conn != nil {
				v.Conn.Close()
			}
		}
		rs.Listener.Close()
		resourceMux.Lock()
		delete(resourceMap, port)
		resourceMux.Unlock()
		log.Println("Close port:", port)
	}()

	log.Println("Open port:", port)
	go func() {
		defer Recover()
		for {
			outConn, err := rs.Listener.Accept()
			if err != nil {
				return
			}

			// 通知客户端建立连接
			ok, id := rs.NewConn(inConn)
			if ok {
				var buffer bytes.Buffer
				buffer.Write([]byte{NEWSOCKET})
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

// doConn 处理客户端新连接
func doConn(conn net.Conn, config *ServerConfig) {
	defer Recover()

	var cmd = make([]byte, 1)
	if _, err := io.ReadAtLeast(conn, cmd, 1); err != nil {
		conn.Close()
		return
	}

	switch cmd[0] {
	case START:
		defer conn.Close()
		// 初始化
		// START info_len info
		info_len := make([]byte, 8)
		if _, err := io.ReadAtLeast(conn, info_len, 8); err != nil {
			return
		}

		var ilen = (uint64(info_len[0]) << 56) | (uint64(info_len[1]) << 48) | (uint64(info_len[2]) << 40) | (uint64(info_len[3]) << 32) | (uint64(info_len[4]) << 24) | (uint64(info_len[5]) << 16) | (uint64(info_len[6]) << 8) | (uint64(info_len[7]))
		if ilen > 1024*1024 {
			// 限制消息最大内存使用量 1M
			return
		}

		var cinfo = make([]byte, ilen)
		if _, err := io.ReadAtLeast(conn, cinfo, int(ilen)); err != nil {
			return
		}

		var cconf ClientConfig
		if nil != json.Unmarshal(cinfo, &cconf) {
			return
		}

		if cconf.Key != config.Key {
			conn.Write([]byte{ERROR_PWD})
			return
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// 打开端口
		for _, cc := range cconf.Map {
			if len(config.LimitPort) > 1 {
				if cc.Outer < config.LimitPort[0] || cc.Outer > config.LimitPort[1] {
					// 不满足端口范围
					log.Printf("Does not meet the port range[%v, %v] %v", config.LimitPort[0], config.LimitPort[1], cc.Outer)
					conn.Write([]byte{ERROR_LIMIT_PORT})
					return
				}
			}

			clis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", cc.Outer))
			if err != nil {
				log.Println("Port is occupied", cc.Outer)
				conn.Write([]byte{ERROR_BUSY})
				return
			}
			resourceMux.Lock()
			resourceMap[cc.Outer] = &Resource{
				Listener: clis,
				Running:  true,
			}
			resourceMux.Unlock()

			go doListen(ctx, conn, cc.Outer)
		}

		conn.Write([]byte{SUCCESS})
		for {
			n, err := conn.Read(cmd)
			if err != nil {
				return
			}
			if n != 0 {
				switch cmd[0] {
				case KILL:
					return
				case IDLE:
					continue
				}
			}
		}
	case NEWCONN:
		// 客户端建立新连接
		sport := make([]byte, 3)
		io.ReadAtLeast(conn, sport, 3)
		port := (uint16(sport[0]<<8) + uint16(sport[1]))
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

				var s encrypto.NCopy
				key, iv := encrypto.GetKeyIv(config.Key)
				s.Init(conn, key, iv)

				go encrypto.WCopy(&s, wk.Conn)
				go encrypto.RCopy(wk.Conn, &s)

				client.mu.Lock()
				defer client.mu.Unlock()
				client.WaitWorker[id] = nil
			}
		} else {
			conn.Close()
		}
	default:
		conn.Close()
	}
}
