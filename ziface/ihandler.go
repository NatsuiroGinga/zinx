package ziface

// IMsgHandler 消息管理抽象层
type IMsgHandler interface {
	Handle(request IRequest) error                     // 处理消息
	RegisterRouter(msgId uint32, router IRouter) error // 为消息添加具体的处理逻辑
}
