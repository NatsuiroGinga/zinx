package znet

import (
	"fmt"
	"io"
	"net"
	"strings"
	"zinx/config"
	errs "zinx/lib/enum/err"
	"zinx/lib/logger"
	"zinx/ziface"
)

// Connection 链接模块
type Connection struct {
	server     ziface.IServer     // 当前链接隶属于哪个server
	conn       *net.TCPConn       // 当前链接的TCPConn
	connId     uint32             // 链接的ID
	isClosed   bool               // 当前的链接状态
	exitChan   chan struct{}      // 告知当前链接已经退出的/停止的channel
	msgHandler ziface.IMsgHandler // 消息管理模块
	msgChan    chan []byte        // 无缓冲管道，用于读、写两个goroutine之间的消息通信
}

// StartReader 链接的读业务方法
func (conn *Connection) startReader() {
	logger.Info("[Reader Goroutine is running...]")
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
			// 判断是否是OpError，如果是则判断是否是链接关闭
			if err == io.EOF ||
				err == io.ErrUnexpectedEOF ||
				strings.Contains(err.Error(), "use of closed network connection") {

				logger.Info("remote addr = ", conn.RemoteAddr(), " is closed")
			} else {
				logger.Error("conn Read error: ", err)
			}

			conn.exitChan <- struct{}{}
			break
		}
		// 拆包，得到msgId和dataLen放在msg中
		if msg, err = dataPack.Unpack(header); err != nil {
			logger.Error("unpack error: ", err)
			conn.exitChan <- struct{}{}
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
		// 将msg封装到request中
		request := NewRequest(conn, msg)
		// 将request发送给TaskQueue，由worker进行处理
		if config.ZinxProperties.WorkerPoolSize > 0 {
			conn.msgHandler.SendMsgToTaskQueue(request)
		} else {
			go func() {
				err := conn.msgHandler.Handle(request)
				if err != nil {
					logger.Error("handle error:", err)
				}
			}()
		}
	}
}

// StartWriter 链接的写业务方法, 从chan中获取数据，然后写给客户端
func (conn *Connection) startWriter() {
	logger.Info("[Writer Goroutine is running...]")
	defer func() {
		logger.Info("connId = ", conn.connId, " Writer is exit, remote addr is ", conn.RemoteAddr())
		conn.Stop()
	}()

	for {
		select {
		case data := <-conn.msgChan:
			if _, err := conn.conn.Write(data); err != nil {
				logger.Error("sendMsgToTaskQueue data error: ", err)
				return
			}
		case <-conn.exitChan:
			return
		}
	}
}

// Start 启动链接，让当前的链接准备开始工作
func (conn *Connection) Start() {
	logger.Info(fmt.Sprintf("connId = %d start work, remote addr is %s", conn.connId, conn.RemoteAddr()))
	// 1. 启动从当前链接的读数据业务
	go conn.startReader()
	// 2. 启动从当前链接写数据业务
	go conn.startWriter()
	// 3. 调用OnConnStart钩子函数
	conn.server.CallOnConnStart(conn)
}

// Stop 停止链接，结束当前链接的工作
func (conn *Connection) Stop() {
	logger.Info(fmt.Sprintf("connId = %d stop work, remote addr is %s", conn.connId, conn.RemoteAddr()))
	if conn.isClosed {
		return
	}
	conn.isClosed = true
	conn.server.CallOnConnStop(conn)
	// 1. 关闭socket链接
	_ = conn.conn.Close()
	// 2. 通知从缓冲队列读数据的业务，该链接已经关闭
	conn.exitChan <- struct{}{}
	// 3. 删除服务器中的链接
	err := conn.server.ConnManager().Remove(conn)
	if err != nil {
		logger.Error("remove conn error: ", err)
	}
	// 4. 关闭该链接全部管道
	close(conn.exitChan)
	close(conn.msgChan)
}

func (conn *Connection) TcpConnection() *net.TCPConn {
	return conn.conn
}

func (conn *Connection) ConnID() uint32 {
	return conn.connId
}

func (conn *Connection) RemoteAddr() net.Addr {
	return conn.conn.RemoteAddr()
}

// Send 发送数据，将数据发送给管道
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
	// 3. 写到管道
	conn.msgChan <- packedBytes

	return
}

// NewConnection 初始化链接模块的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	// 初始化Conn属性
	connection := &Connection{
		conn:       conn,
		connId:     connID,
		isClosed:   false,
		exitChan:   make(chan struct{}, 1),
		msgHandler: msgHandler,
		msgChan:    make(chan []byte),
		server:     server,
	}
	// 将新创建的Conn添加到链接管理中
	server.ConnManager().Add(connection)

	return connection
}
