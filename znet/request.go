package znet

import (
	"zinx/ziface"
)

type Request struct {
	conn ziface.IConnection
	data []byte
}

func NewRequest(conn ziface.IConnection, data []byte) *Request {
	return &Request{conn: conn, data: data}
}

func (request *Request) Connection() ziface.IConnection {
	return request.conn
}

func (request *Request) Data() []byte {
	return request.data
}
