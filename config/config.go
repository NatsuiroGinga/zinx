package config

import (
	"encoding/json"
	"fmt"
	"os"
	"zinx/ziface"
)

// ServerProperties 服务器配置参数
var ServerProperties *serverProperties

// serverProperties 服务器配置参数
type serverProperties struct {
	Port      int            `json:"port"`       // 服务器端口
	Name      string         `json:"name"`       // 服务器名称
	Host      string         `json:"host"`       // 服务器IP
	TCPServer ziface.IServer `json:"tcp-server"` // 当前Zinx全局的Server对象
}

func (properties *serverProperties) String() string {
	return fmt.Sprintf("[Server] Name: %s, Host: %s, Port: %d", properties.Name, properties.Host, properties.Port)
}

// ZinxProperties 框架配置参数
var ZinxProperties *zinxProperties

// zinxProperties 框架配置参数
type zinxProperties struct {
	MaxConnections int    `json:"max-connections"`  // 当前服务器主机允许的最大连接数
	MaxPackageSize int    `json:"max-package-size"` // 框架的数据包的最大值
	Version        string `json:"version"`          // Zinx版本
}

func (properties *zinxProperties) String() string {
	return fmt.Sprintf("[Zinx] Version: %s, MaxConnections: %d, MaxPackageSize: %d", properties.Version, properties.MaxConnections, properties.MaxPackageSize)
}

// 初始化配置参数
func init() {
	ServerProperties = &serverProperties{
		Port:      8848,
		Name:      "ZinxServerApp",
		Host:      "0.0.0.0",
		TCPServer: nil,
	}

	ZinxProperties = &zinxProperties{
		MaxConnections: 1000,
		MaxPackageSize: 4 << 10,
		Version:        "v0.4",
	}

	loadFile("config/zinx.json")
}

// LoadFile 加载配置文件
func loadFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &struct {
		*serverProperties `json:"server"`
		*zinxProperties   `json:"zinx"`
	}{
		ServerProperties,
		ZinxProperties,
	})
	if err != nil {
		panic(err)
	}
}