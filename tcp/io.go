package tcp

import (
	"io"
	"net"
)

func proxyConn(source, destination string) error {
	connSource, err := net.Dial("tpc", source)
	if err != nil {
		return err
	}
	defer connSource.Close()

	connDestination, err := net.Dial("tcp", destination)
	if err != nil {
		return err
	}
	defer connDestination.Close()

	go func() {
		_, _ = io.Copy(connSource, connDestination)
	}()
	_, err = io.Copy(connDestination, connSource)
	return err
}
