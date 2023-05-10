package znet

import (
	"net"
	"zinx/lib/logger"
	"zinx/ziface"
)

// Connection 链接模块
type Connection struct {
	Conn     *net.TCPConn      // 当前链接的TCPConn
	ConnID   uint32            // 链接的ID
	isClosed bool              // 当前的链接状态
	handler  ziface.HandleFunc // 当前链接所绑定的处理业务方法
	ExitChan chan struct{}     // 告知当前链接已经退出的/停止的channel
}

// StartReader 链接的读业务方法
func (conn *Connection) StartReader() {
	logger.Info("Reader Goroutine is running...")
	defer func() {
		logger.Info("ConnID = ", conn.ConnID, " Reader is exit, remote addr is ", conn.RemoteAddr().String())
		conn.Stop()
	}()

	buf := make([]byte, 512)
	var (
		n   int
		err error
	)
	for {
		// 读取客户端的数据到buf中，最大512字节
		if n, err = conn.Conn.Read(buf); err != nil {
			logger.Error("recv buf err", err)
			continue
		}
		// 调用当前链接所绑定的handle方法
		if err = conn.handler(conn.Conn, buf, n); err != nil {
			logger.Error("ConnID = ", conn.ConnID, "handle is error", err)
			break
		}
	}
}

// Start 启动链接，让当前的链接准备开始工作
func (conn *Connection) Start() {
	logger.Info("Conn Start()...ConnID = ", conn.ConnID)
	// 1. 启动从当前链接的读数据业务
	go conn.StartReader()
	// 2. 启动从当前链接写数据业务
}

// Stop 停止链接，结束当前链接的工作
func (conn *Connection) Stop() {
	logger.Info("Conn Stop()...ConnID = ", conn.ConnID)
	if conn.isClosed {
		return
	}
	conn.isClosed = true
	// 1. 关闭socket链接
	_ = conn.Conn.Close()
	// 2. 关闭该链接全部管道
	close(conn.ExitChan)
}

func (conn *Connection) GetTCPConnection() *net.TCPConn {
	//TODO implement me
	panic("implement me")
}

func (conn *Connection) GetConnID() uint32 {
	//TODO implement me
	panic("implement me")
}

func (conn *Connection) RemoteAddr() net.Addr {
	//TODO implement me
	panic("implement me")
}

func (conn *Connection) Send(data []byte) error {
	//TODO implement me
	panic("implement me")
}

// NewConnection 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, handler ziface.HandleFunc) *Connection {
	return &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		handler:  handler,
		ExitChan: make(chan struct{}, 1),
	}
}
