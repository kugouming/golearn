package types

import "time"

const (
	_                uint8 = iota // 初始化无意义
	START                         // START 第一次连接服务器
	NEWSOCKET                     // NEWSOCKET 新连接
	NEWCONN                       // NEWCONN 新连接发送到服务端命令
	ERROR                         // ERROR 处理失败
	SUCCESS                       // SUCCESS 处理成功
	IDLE                          // IDLE 空闲命令 什么也不做
	KILL                          // KILL 退出命令
	ERROR_PWD                     // ERROR_PWD 密码错误
	ERROR_BUSY                    // ERROR_BUSY 端口被占用
	ERROR_LIMIT_PORT              // ERROR_LIMIT_PORT 不满足端口范围
)

const (
	RetryTime          = time.Second      // RetryTime 断线重连时间
	TcpKeepAlivePeriod = 30 * time.Second // 保持请求链接时间
	WaitTimeOut        = 30 * time.Second // 连接等待超时时间
	WaitMax            = 10               // 等待队列长度
)
