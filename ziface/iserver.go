package ziface

// startStop interface defines the start and stop method
type startStop interface {
	Start() // Start
	Stop()  // Stop
}

// IServer interface defines the server interface
type IServer interface {
	startStop
	Serve() // start the server
}
