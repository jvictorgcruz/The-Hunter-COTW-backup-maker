package instance

import (
	"net"
	"time"
)

const port = "127.0.0.1:45321"

func TryLock() (bool, <-chan string, error) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return false, nil, nil
	}

	ch := make(chan string, 10)
	go func() {
		defer listener.Close()
		defer close(ch)
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			buf := make([]byte, 10)
			n, err := conn.Read(buf)
			if err == nil {
				ch <- string(buf[:n])
			}
			conn.Close()
		}
	}()

	return true, ch, nil
}

func NotifyExisting() {
	conn, err := net.DialTimeout("tcp", port, 2*time.Second)
	if err != nil {
		return
	}
	defer conn.Close()
	_, _ = conn.Write([]byte("show"))
}
