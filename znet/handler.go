package znet

import (
	"fmt"
	errs "zinx/lib/enum/err"
	"zinx/lib/logger"
	"zinx/lib/util"
	"zinx/ziface"
)

// MsgHandler 消息管理抽象层
type MsgHandler struct {
	routers map[uint32]ziface.IRouter
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{routers: make(map[uint32]ziface.IRouter)}
}

func (handler *MsgHandler) Handle(request ziface.IRequest) (err error) {
	if router, ok := handler.routers[request.Message().ID()]; ok {
		router.PreHandle(request)
		router.Handle(request)
		router.PostHandle(request)
		return
	}
	err = util.NewErrorWithPattern(errs.MESSAGE_NOT_REGISTERED, request.Message().ID())
	logger.Error(err)
	return
}

func (handler *MsgHandler) RegisterRouter(msgId uint32, router ziface.IRouter) (err error) {
	if _, ok := handler.routers[msgId]; ok {
		err = util.NewErrorWithPattern(errs.MESSAGE_REGISTERED, msgId)
		logger.Error(err)
		return
	}
	handler.routers[msgId] = router
	logger.Info(fmt.Sprintf("msgId: %d register successfully", msgId))

	return
}
