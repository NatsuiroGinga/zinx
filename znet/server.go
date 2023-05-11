package znet

import (
	"fmt"
	"net"
	"time"
	"zinx/config"
	"zinx/lib/logger"
	"zinx/ziface"
)

// IP_VERSION 默认配置
const IP_VERSION = "tcp4"

// Server defines the server struct
type Server struct {
	name      string         // server name
	ipVersion string         // server bind ip version
	ip        string         // server bind ip
	port      int            // server bind port
	router    ziface.IRouter // router
}

func (server *Server) RegisterRouter(router ziface.IRouter) ziface.IServer {
	server.router = router
	logger.Info("Add Router success!")
	return server
}

func NewServer() *Server {
	return &Server{
		name:      config.ServerProperties.Name,
		ipVersion: IP_VERSION,
		ip:        config.ServerProperties.Host,
		port:      config.ServerProperties.Port,
	}
}

// NewServerWithRouter creates a server with router
func NewServerWithRouter(name string, router ziface.IRouter) (server *Server) {
	return NewServer().RegisterRouter(router).(*Server)
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
	logger.Info(config.ServerProperties, "is starting...")
	logger.Info(config.ZinxProperties)

	go func() {
		// 1. 创建socket
		addr, err := net.ResolveTCPAddr(server.ipVersion, fmt.Sprintf("%s:%d", server.ip, server.port))
		if err != nil {
			logger.Error(fmt.Sprintf("Resolve TCP Address failed: %s", err.Error()))
			return
		}

		// 2. 监听服务器地址
		listener, err := net.ListenTCP(server.ipVersion, addr)
		if err != nil {
			logger.Error(fmt.Sprintf("Listen TCP Address failed: %s", err.Error()))
			return
		}
		logger.Info(fmt.Sprintf("Server Listener at ip: %s, port: %d, is started", server.ip, server.port))
		cid := uint32(0)

		// 3. 阻塞等待客户端连接，处理客户端连接业务
		for {
			// 3.1 如果有客户端连接，阻塞返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				logger.Error(fmt.Sprintf("Accept TCP failed: %s", err.Error()))
				continue
			}
			dealConn := NewConnection(conn, cid, server.router)
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
