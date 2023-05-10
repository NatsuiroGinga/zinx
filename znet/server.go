package znet

import (
	"log"
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
	log.Printf("[Zinx] Server Name: %s starting...", server.Name)
}

func (server *Server) Stop() {
	log.Printf("[Zinx] Server Name: %s stoping...", server.Name)
}

func (server *Server) Serve() {
	log.Printf("[Zinx] Server Name: %s serving...", server.Name)
}
