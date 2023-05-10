package handler

import (
	"errors"
	"net"
	"zinx/lib/logger"
)

func EchoHandler(conn *net.TCPConn, data []byte, n int) error {
	logger.Info("[Conn Handle] Echo...")
	if _, err := conn.Write(data[:n]); err != nil {
		logger.Error("Write back buf err", err)
		return errors.New("write back buf err")
	}
	return nil
}
