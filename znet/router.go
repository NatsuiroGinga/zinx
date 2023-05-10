package znet

import (
	"zinx/ziface"
)

// BaseRouter 基础路由, 供其他路由继承
type BaseRouter struct {
}

func (router *BaseRouter) PreHandle(ziface.IRequest) {}

func (router *BaseRouter) Handle(ziface.IRequest) {}

func (router *BaseRouter) PostHandle(ziface.IRequest) {}
