package ziface

// IRequest 接口
type IRequest interface {
	Connection() IConnection // 获取请求连接信息
	Data() []byte            // 获取请求消息的数据
}
