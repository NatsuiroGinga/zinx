package ziface

// startStop interface defines the start and stop method
type startStop interface {
	Start() // Start
	Stop()  // Stop
}

// IServer interface defines the server interface
type IServer interface {
	startStop
	Serve()                                            // start the server and serve
	RegisterRouter(msgID uint32, router IRouter) error // register router
	ConnManager() IConnManager                         // connection manager
	SetOnConnStart(hook Hook)                          // set hook function for connection start
	SetOnConnStop(hook Hook)                           // set hook function for connection stop
	CallOnConnStart(conn IConnection)                  // call hook function for connection start
	CallOnConnStop(conn IConnection)                   // call hook function for connection stop
}

// Hook 定义 Server 的 Hook 函数原型
type Hook func(conn IConnection)
