package ziface

// IServer interface defines the server interface
type IServer interface {
	Start()
	Stop()
	Serve()
}
