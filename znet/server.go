package znet

import (
	"fmt"
	"net"
	"time"
	"zinx/handler"
	"zinx/lib/logger"
	"zinx/ziface"
)

// 默认配置
const (
	IPVersion = "tcp4"
	IP        = "0.0.0.0"
	Port      = 8848
)

// Server defines the server struct
type Server struct {
	Name      string // server name
	IPVersion string // server bind ip version
	IP        string // server bind ip
	Port      int    // server bind port
}

func NewServer(name string) ziface.IServer {
	return &Server{
		Name:      name,
		IPVersion: IPVersion,
		IP:        IP,
		Port:      Port,
	}
}

func init() {
	logger.Setup(&logger.Settings{
		Path:       "./logs",      // 日志文件路径
		Name:       "zinx",        // 日志文件名称
		Ext:        "log",         // 日志文件后缀
		TimeFormat: time.DateOnly, // 时间格式
	})
}

func (server *Server) Start() {
	logger.Info(fmt.Sprintf("Server Listener at IP: %s, Port: %d, is starting...", server.IP, server.Port))
	go func() {
		// 1. 创建socket
		addr, err := net.ResolveTCPAddr(server.IPVersion, fmt.Sprintf("%s:%d", server.IP, server.Port))
		if err != nil {
			logger.Error(fmt.Sprintf("Resolve TCP Address failed: %s", err.Error()))
			return
		}
		// 2. 监听服务器地址
		listener, err := net.ListenTCP(server.IPVersion, addr)
		if err != nil {
			logger.Error(fmt.Sprintf("Listen TCP Address failed: %s", err.Error()))
			return
		}
		logger.Info(fmt.Sprintf("Server Listener at IP: %s, Port: %d, is started", server.IP, server.Port))
		cid := uint32(0)
		// 3. 阻塞等待客户端连接，处理客户端连接业务
		for {
			// 3.1 如果有客户端连接，阻塞返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				logger.Error(fmt.Sprintf("Accept TCP failed: %s", err.Error()))
				continue
			}
			dealConn := NewConnection(conn, cid, handler.EchoHandler)
			cid++
			// 3.2 处理客户端业务
			logger.Info(fmt.Sprintf("Accept a client, Address: %s", conn.RemoteAddr()))
			go dealConn.Start()
		}
	}()
}

func (server *Server) Stop() {
	logger.Info("Server is stopping...")
}

func (server *Server) Serve() {
	logger.Info("Server is serving...")
	server.Start() // 启动server服务
	select {}      // 阻塞
}
