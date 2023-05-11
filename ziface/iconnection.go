package ziface

import (
	"net"
)

// IConnection 定义连接模块的抽象层
type IConnection interface {
	startStop                             // 开启链接，停止链接
	TcpConnection() *net.TCPConn          // 获取当前连接的socket TCPConn
	ConnID() uint32                       // 获取当前连接ID
	RemoteAddr() net.Addr                 // 获取远程客户端地址信息
	Send(msgId uint32, data []byte) error // 直接将数据发送给远程的TCP连接
}
