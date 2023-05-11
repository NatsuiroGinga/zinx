package ziface

// IRequest 接口
type IRequest interface {
	Connection() IConnection // 获取请求连接信息
	Message() IMessage       // 获取请求消息
}
