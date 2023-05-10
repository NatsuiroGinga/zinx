package znet

import (
	"bufio"
	"fmt"
	"net"
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

func (server *Server) Start() {
	go logger.Info(fmt.Sprintf("Server Listener at IP: %s, Port: %d, is starting...", server.IP, server.Port))
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
	go logger.Info(fmt.Sprintf("Server Listener at IP: %s, Port: %d, is started", server.IP, server.Port))
	// 3. 阻塞等待客户端连接，处理客户端连接业务
	for {
		// 3.1 如果有客户端连接，阻塞返回
		conn, err := listener.AcceptTCP()
		if err != nil {
			logger.Error(fmt.Sprintf("Accept TCP failed: %s", err.Error()))
			continue
		}
		// 3.2 处理客户端业务
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	reader := bufio.NewReader(conn)
	buf := make([]byte, 512)
	var (
		cnt int
		err error
	)

	for {
		cnt, err = reader.Read(buf)
		if err != nil {
			logger.Error(fmt.Sprintf("Read TCP failed: %s", err.Error()))
			continue
		}

		logger.Info(fmt.Sprintf("Receive from client, cnt: %d, buf: %s", cnt, buf))

		if _, err = conn.Write(buf[:cnt]); err != nil {
			logger.Error(fmt.Sprintf("Write TCP failed: %s", err.Error()))
			continue
		}
	}
}

func (server *Server) Stop() {
	logger.Info("Server is stopping...")
}

func (server *Server) Serve() {
	logger.Info("Server is serving...")
	server.Start() // 启动server服务
}
