package znet

import (
	"fmt"
	"sync"
	"time"
	"zinx/config"
	errs "zinx/lib/enum/err"
	"zinx/lib/logger"
	"zinx/lib/util"
	"zinx/ziface"
)

// MsgHandler 消息管理抽象层
type MsgHandler struct {
	routers        sync.Map               /*map[uint32]ziface.IRouter*/ // 消息ID和路由的映射关系
	taskQueue      []chan ziface.IRequest // 任务队列
	workerPoolSize uint32                 // 业务工作Worker池的数量
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		workerPoolSize: config.ZinxProperties.WorkerPoolSize,
		taskQueue:      make([]chan ziface.IRequest, config.ZinxProperties.WorkerPoolSize),
	}
}

// Handle 处理消息
func (handler *MsgHandler) Handle(request ziface.IRequest) (err error) {
	if router, ok := handler.routers.Load(request.Message().ID()); ok {
		r, exist := router.(ziface.IRouter)
		if !exist {
			err = util.NewErrorWithPattern(errs.MESSAGE_NOT_REGISTERED, request.Message().ID())
			logger.Error(err)
			return
		}
		r.PreHandle(request)
		r.Handle(request)
		r.PostHandle(request)
		return
	}
	err = util.NewErrorWithPattern(errs.MESSAGE_NOT_REGISTERED, request.Message().ID())
	logger.Error(err)
	return
}

// RegisterRouter 注册路由
func (handler *MsgHandler) RegisterRouter(msgId uint32, router ziface.IRouter) (err error) {
	if _, ok := handler.routers.Load(msgId); ok {
		err = util.NewErrorWithPattern(errs.MESSAGE_REGISTERED, msgId)
		logger.Error(err)
		return
	}
	handler.routers.Store(msgId, router)
	logger.Info(fmt.Sprintf("msgId: %d register successfully", msgId))
	return
}

// StartWorkerPool 启动Worker池
func (handler *MsgHandler) StartWorkerPool() {
	logger.Info("start worker pool")
	for i := uint32(0); i < handler.workerPoolSize; i++ {
		// 1. 给当前的worker对应的channel消息队列开辟空间
		handler.taskQueue[i] = make(chan ziface.IRequest, config.ZinxProperties.MaxWorkerTaskLen())
		// 2. 启动当前的Worker，阻塞等待消息从channel传递进来
		go handler.startOneWorker(i, handler.taskQueue[i])
	}
}

// startOneWorker 启动一个Worker工作流程
func (handler *MsgHandler) startOneWorker(workerId uint32, taskQueue chan ziface.IRequest) {
	logger.Info(fmt.Sprintf("worker ID: %d is started", workerId))
	// 不断的阻塞等待对应消息队列的消息
	for request := range taskQueue {
		err := handler.Handle(request)
		if err != nil {
			logger.Error(err)
		}
	}
}

// SendMsgToTaskQueue 将消息交给TaskQueue，由worker进行处理
func (handler *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	// hash取模得到当前的workerId
	workerId := request.Connection().ConnID() % handler.workerPoolSize
	logger.Info(fmt.Sprintf("Add ConnID: %d, request msgID: %d to workerID: %d", request.Connection().ConnID(), request.Message().ID(), workerId))
	// 轮询 查找空闲的worker, 将消息发送给worker的任务队列, 由worker进行处理
	// 防止超时
	timer := time.NewTimer(time.Second * 3)
	for {
		select {
		case handler.taskQueue[workerId] <- request: // 将消息发送给worker的任务队列
			logger.Info(fmt.Sprintf("Add request to worker ID: %d successfully", workerId))
			return
		case <-timer.C: // 说明当前的worker_pool满了
			logger.Error(errs.WORKER_POOL_FULL)
			timer.Stop()
			return
		default: // 选择下一个worker
			logger.Warn("worker pool is full, try next worker")
			workerId = (workerId + 1) % handler.workerPoolSize
		}
	}
}
