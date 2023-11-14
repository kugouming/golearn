package types

import (
	"net"
	"sync"
	"time"
)

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
	Server *ServerConfig `json:"server"` // Server 端配置
	Client *ClientConfig `json:"client"` // Client 端配置
}

// Worker 工作线程结构体
type Worker struct {
	Conn     net.Conn // 客户端连接
	LastTime int64    // 客户端连接超时时间
}

type Resource struct {
	Listener   net.Listener     // Net 监听句柄
	WaitWorker [WaitMax]*Worker // 工作负载
	Running    bool             // 是否在运行
	Mux        sync.Mutex       // 工作负载锁
}

// NewConn 建立新连接
func (r *Resource) NewConn(conn net.Conn) (bool, uint8) {
	r.Mux.Lock()
	defer r.Mux.Unlock()

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
