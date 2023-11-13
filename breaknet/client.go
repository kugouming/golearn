package main

import (
	"breaknet/encrypto"
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
	"time"
)

func DoClient(config *ClientConfig) {
	var pmap = make(map[uint16]string, len(config.Map))
	for _, m := range config.Map {
		pmap[m.Outer] = m.Inner
	}

	var isContinue = true

	// 新连接处理
	var doConn = func(conn net.Conn, port uint16, sp []byte) {
		defer Recover()
		lconn, err := net.Dial("tcp", pmap[port])
		if err != nil {
			conn.Close()
			log.Println(err)
			return
		}

		conn.Write([]byte{NEWCONN})
		conn.Write(sp)
		key, iv := encrypto.GetKeyIv(config.Key)
		var s encrypto.NCopy
		s.Init(conn, key, iv)
		go encrypto.WCopy(&s, lconn)
		go encrypto.RCopy(lconn, &s)
	}

	for isContinue {
		func() {
			defer Recover()
			defer time.Sleep(RetryTime)
			log.Println("Connecting to server ...")

			sconn, err := net.Dial("tcp", config.Server)
			if err != nil {
				log.Println("Cann't connect to server")
				return
			}
			defer sconn.Close()

			sconn.(*net.TCPConn).SetKeepAlive(true)
			sconn.(*net.TCPConn).SetKeepAlivePeriod(TcpKeepAlivePeriod)

			cinfo, _ := json.Marshal(config)

			// 添加字节缓冲
			var buffer bytes.Buffer

			// 发送客户端信息
			// START info_len info
			buffer.Write([]byte{START})
			binary.Write(&buffer, binary.BigEndian, uint64(len(cinfo)))
			buffer.Write(cinfo)
			sconn.Write(buffer.Bytes())
			buffer.Reset()

			// 读取返回信息
			// SUCCESS / ERROR / BUSY
			var cmd = make([]byte, 1)
			io.ReadAtLeast(sconn, cmd, 1)
			switch cmd[0] {
			case ERROR_PWD:
				log.Println("Wrong password")
				isContinue = false
				return
			case ERROR_BUSY:
				log.Println("Port is occupied")
				isContinue = false
				return
			case ERROR_LIMIT_PORT:
				log.Println("Does not meet the port range")
				isContinue = false
				return
			}

			if cmd[0] != SUCCESS {
				log.Println("Unkown error")
				isContinue = false
				return
			}

			log.Println("Certification sucessful")
			for _, cc := range config.Map {
				log.Printf("%v->:%v\n", cc.Inner, cc.Outer)
			}

			cmd[0] = IDLE
			// 进入指令读取循环
			for {
				_, err := sconn.Read(cmd)
				if err != nil {
					return
				}

				switch cmd[0] {
				case NEWSOCKET:
					// 新建连接
					// 读取远端端口与ID
					sp := make([]byte, 3)
					io.ReadAtLeast(sconn, sp, 3)
					sport := uint16(sp[0]<<8) + uint16(sp[1])
					conn, err := net.Dial("tcp", config.Server)
					if err != nil {
						return
					}
					go doConn(conn, sport, sp)
				case IDLE:
					_, err := sconn.Write([]byte{SUCCESS})
					if err != nil {
						return
					}
				}
			}
		}()
	}
}
