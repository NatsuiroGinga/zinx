package znet

import (
	"net"
	"zinx/config"
	"zinx/lib/logger"
	"zinx/lib/util"
	"zinx/ziface"
)

// Connection 链接模块
type Connection struct {
	conn     *net.TCPConn   // 当前链接的TCPConn
	connId   uint32         // 链接的ID
	isClosed bool           // 当前的链接状态
	exitChan chan struct{}  // 告知当前链接已经退出的/停止的channel
	router   ziface.IRouter // 当前链接所绑定的router对象
}

// StartReader 链接的读业务方法
func (conn *Connection) startReader() {
	logger.Info("Reader Goroutine is running...")
	defer func() {
		logger.Info("connId = ", conn.connId, " Reader is exit, remote addr is ", conn.RemoteAddr())
		conn.Stop()
	}()

	buf := make([]byte, config.ZinxProperties.MaxPackageSize)
	var (
		n   int
		err error
	)
	for {
		// 读取客户端的数据到buf中，最大512字节
		if n, err = conn.conn.Read(buf); err != nil {
			logger.Error("recv buf err", err)
			return
		}
		// 调用当前链接所绑定的router的PreHandle方法等
		request := NewRequest(conn, buf[:n])
		go func() {
			conn.router.PreHandle(request)
			conn.router.Handle(request)
			conn.router.PostHandle(request)
		}()
	}
}

// Start 启动链接，让当前的链接准备开始工作
func (conn *Connection) Start() {
	logger.Info("conn Start()...connId = ", conn.connId)
	// 1. 启动从当前链接的读数据业务
	go conn.startReader()
	// 2. 启动从当前链接写数据业务
}

// Stop 停止链接，结束当前链接的工作
func (conn *Connection) Stop() {
	logger.Info("conn Stop()...connId = ", conn.connId)
	if conn.isClosed {
		return
	}
	conn.isClosed = true
	// 1. 关闭socket链接
	_ = conn.conn.Close()
	// 2. 关闭该链接全部管道
	close(conn.exitChan)
}

func (conn *Connection) TcpConnection() *net.TCPConn {
	//TODO implement me
	panic("implement me")
}

func (conn *Connection) ConnID() uint32 {
	//TODO implement me
	panic("implement me")
}

func (conn *Connection) RemoteAddr() net.Addr {
	//TODO implement me
	panic("implement me")
}

func (conn *Connection) Send(data []byte) error {
	if _, err := conn.conn.Write(data); err != nil {
		logger.Error("Send data error: ", err)
		return err
	}
	logger.Info("Send data success:", util.Bytes2String(data))
	return nil
}

// NewConnection 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	return &Connection{
		conn:     conn,
		connId:   connID,
		isClosed: false,
		exitChan: make(chan struct{}, 1),
		router:   router,
	}
}
