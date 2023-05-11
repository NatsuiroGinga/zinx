package znet

import (
	"io"
	"net"
	errs "zinx/lib/enum/err"
	"zinx/lib/logger"
	"zinx/ziface"
)

// Connection 链接模块
type Connection struct {
	conn       *net.TCPConn       // 当前链接的TCPConn
	connId     uint32             // 链接的ID
	isClosed   bool               // 当前的链接状态
	exitChan   chan struct{}      // 告知当前链接已经退出的/停止的channel
	msgHandler ziface.IMsgHandler // 消息管理模块
}

// StartReader 链接的读业务方法
func (conn *Connection) startReader() {
	logger.Info("Reader Goroutine is running...")
	defer func() {
		logger.Info("connId = ", conn.connId, " Reader is exit, remote addr is ", conn.RemoteAddr())
		conn.Stop()
	}()

	var (
		msg ziface.IMessage
		err error
	)
	header := make([]byte, dataPack.HeadLen())

	for {
		// 读取客户端的数据到buf中
		if _, err = io.ReadFull(conn.conn, header); err != nil {
			logger.Error("conn Read error: ", err)
			break
		}
		// 拆包，得到msgId和dataLen放在msg中
		if msg, err = dataPack.Unpack(header); err != nil {
			logger.Error("unpack error: ", err)
			break
		}
		// 根据dataLen再次读取data，放在msg.Data中
		if msg.DataLen() > 0 {
			data := make([]byte, msg.DataLen())
			if _, err = io.ReadFull(conn.conn, data); err != nil {
				logger.Error("read msg data error : ", err)
				break
			}
			msg.SetData(data)
		}
		// 调用当前链接所绑定的router的PreHandle方法等
		request := NewRequest(conn, msg)

		// 从路由中，找到注册绑定的Conn对应的router调用
		go func() {
			err = conn.msgHandler.Handle(request)
			if err != nil {
				logger.Error("connId = ", conn.connId, " handle error: ", err)
				return
			}
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

func (conn *Connection) Send(msgId uint32, data []byte) (err error) {
	// 1. 判断当前链接是否已经关闭
	if conn.isClosed {
		return errs.CONNECT_CLOSED
	}
	// 2. 将data进行封包，并且发送
	message := NewMessage(msgId, data)
	var packedBytes []byte
	if packedBytes, err = dataPack.Pack(message); err != nil {
		logger.Error("pack error: ", err)
		return
	}
	// 3. 写回客户端
	if _, err = conn.conn.Write(packedBytes); err != nil {
		logger.Error("write error: ", err)
		return
	}

	return
}

// NewConnection 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	return &Connection{
		conn:       conn,
		connId:     connID,
		isClosed:   false,
		exitChan:   make(chan struct{}, 1),
		msgHandler: msgHandler,
	}
}
