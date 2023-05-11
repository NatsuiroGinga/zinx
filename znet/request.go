package znet

import (
	"zinx/ziface"
)

// Request 封装请求的数据
type Request struct {
	conn    ziface.IConnection // 已经和客户端建立好的连接
	message ziface.IMessage    // 客户端请求的消息
}

func NewRequest(conn ziface.IConnection, message ziface.IMessage) *Request {
	return &Request{conn: conn, message: message}
}

func (request *Request) Connection() ziface.IConnection {
	return request.conn
}

func (request *Request) Message() ziface.IMessage {
	return request.message
}
