package znet

import (
	"fmt"
	"net"
	"time"
	"zinx/config"
	errs "zinx/lib/enum/err"
	"zinx/lib/logger"
	"zinx/lib/util"
	"zinx/ziface"
)

// IP_VERSION 默认配置
const IP_VERSION = "tcp4"

// Server defines the server struct
type Server struct {
	name       string             // server name
	ipVersion  string             // server bind ip version
	ip         string             // server bind ip
	port       int                // server bind port
	msgHandler ziface.IMsgHandler // handler for msg
}

func (server *Server) RegisterRouter(msgID uint32, router ziface.IRouter) (err error) {
	err = server.msgHandler.RegisterRouter(msgID, router)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info("Add Router success!")

	return
}

func NewServer() *Server {
	return &Server{
		name:      config.ServerProperties.Name,
		ipVersion: IP_VERSION,
		ip:        config.ServerProperties.Host,
		port:      config.ServerProperties.Port,
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
		logger.Info(fmt.Sprintf("Server %s is started", server.name))
		cid := uint32(0)

		// 3. 阻塞等待客户端连接，处理客户端连接业务
		for {
			// 3.1 如果有客户端连接，阻塞返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				logger.Error(util.NewErrorWithPattern(errs.ACCEPT_TCP_FAILED, err.Error()))
				continue
			}
			dealConn := NewConnection(conn, cid, server.msgHandler)
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
