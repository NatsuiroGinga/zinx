package ziface

import (
	"net"
)

type IConnection interface {
	startStop                       // 开启链接，停止链接
	GetTCPConnection() *net.TCPConn // 获取当前连接的socket TCPConn
	GetConnID() uint32              // 获取当前连接ID
	RemoteAddr() net.Addr           // 获取远程客户端地址信息
	Send(data []byte) error         // 直接将数据发送给远程的TCP连接
}

// HandleFunc 定义一个统一处理连接业务的接口
//
// conn 当前连接, data 数据字节数组, n 实际数据长度
type HandleFunc func(conn *net.TCPConn, data []byte, n int) error
