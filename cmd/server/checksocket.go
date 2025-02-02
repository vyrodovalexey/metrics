package main

import (
	"errors"
	"go.uber.org/zap"
	"net"
	"syscall"
	"time"
)

func IsPortAvailable(addr string, logger *zap.SugaredLogger) error {
	var err error
	for i := 0; i <= 3; i++ {
		var listener net.Listener
		listener, err = net.Listen("tcp", addr)
		if err != nil {
			// Check if the error is due to the port being already in use
			if errors.Is(err, syscall.EADDRINUSE) {
				logger.Infof("%v", err)
			} else {
				// Return other errors (e.g., permission denied)
				return err
			}
		}
		defer listener.Close()
		if i == 0 {
			<-time.After(1 * time.Second)
		} else {
			<-time.After(time.Duration(i*2+1) * time.Second)
		}
	}
	return err
}
