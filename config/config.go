package config

import (
	"encoding/json"
	"fmt"
	"os"
	"zinx/ziface"
)

/* 配置文件格式
{
  "server": {
    "name": "zinx",
    "host": "127.0.0.1",
    "port": 8848
  },
  "zinx": {
    "max-connections": 1000,
    "max-package-size": 4096
  }
}
*/

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
	version        string
	MaxConnections int    `json:"max-connections"`  // 当前服务器主机允许的最大连接数
	MaxPackageSize uint32 `json:"max-package-size"` // 框架的数据包的最大值
}

func (properties *zinxProperties) String() string {
	return fmt.Sprintf("[Zinx] Version: %s, MaxConnections: %d, MaxPackageSize: %d", properties.version, properties.MaxConnections, properties.MaxPackageSize)
}

// 初始化配置参数
func init() {
	filename := "config/zinx.json"
	ServerProperties = &serverProperties{
		Port:      8848,
		Name:      "ZinxServerApp",
		Host:      "0.0.0.0",
		TCPServer: nil,
	}

	ZinxProperties = &zinxProperties{
		MaxConnections: 1000,
		MaxPackageSize: 4 << 10,
		version:        "v0.5",
	}
	if fileExists(filename) {
		loadFile(filename)
	}
}

// fileExists 判断文件是否存在
func fileExists(filename string) bool {
	stat, err := os.Stat(filename)
	return err == nil && !stat.IsDir()
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
