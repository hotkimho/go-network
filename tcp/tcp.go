package tcp

import (
	"context"
	"fmt"
	"io"
	"net"
	"syscall"
	"time"
)

func Listener() {
	//리스너에 IP와 포트 번호가 바인딩됨
	listener, err := net.Listen("tcp", "127.0.0.1:3000")
	if err != nil {
		fmt.Println("listener(): Error net.Listen func")
		return
	}
	defer func() { _ = listener.Close() }()
	fmt.Printf("연결된 주소 %s", listener.Addr())
	for {
		//수신 연결, 실패 시 err 반환
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go func(c net.Conn) {
			defer c.Close()
			fmt.Println("go routine : ", c.LocalAddr())
		}(conn)
	}
}

func ListenerAndDial() {
	//포트가 없으면 랜덤으로!!
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		fmt.Println(err)
		return
	}
	done := make(chan struct{})
	go func() {
		defer func() { done <- struct{}{} }()
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println(err)
				return
			}
			go func(c net.Conn) {
				defer func() {
					c.Close()
					done <- struct{}{}
				}()

				buf := make([]byte, 1024)
				for {
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							fmt.Println(err)
						}
						return
					}
					fmt.Printf("received: %q\n", buf[:n])
				}
			}(conn)
		}
	}()
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		fmt.Println(err)
	}
	conn.Close()
	<-done
	listener.Close()
	<-done
}

func DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	d := net.Dialer{
		Control: func(_, addr string, _ syscall.RawConn) error {
			return &net.DNSError{
				Err:         "connection timed out",
				Name:        addr,
				Server:      "127.0.0.1",
				IsTimeout:   true,
				IsTemporary: true,
			}
		},
		Timeout: timeout,
	}
	return d.Dial(network, address)
}

func TestDialTimeout() {
	c, err := DialTimeout("tcp", "10.0.0.1:http", 5*time.Second)
	if err == nil {
		c.Close()
		fmt.Println("connection did not time out")
		return
	}
	nErr, ok := err.(net.Error)
	if !ok {
		fmt.Println(err)
		return
	}
	if !nErr.Timeout() {
		fmt.Println("error is not a timeout")
	}
}

func DialContext() {
	deadLine := time.Now().Add(3 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadLine)
	defer cancel()

	var d net.Dialer
	d.Control = func(_, _ string, _ syscall.RawConn) error {
		//context 의 데드라인보다 오랫동안(+1초) 대기합니다.
		return nil
	}
	conn, err := d.DialContext(ctx, "tcp", "127.0.0.1:3000")
	if err == nil {
		conn.Close()
		fmt.Println("타임아웃이 발생하지 않았습니다")
		return
	}
	nErr, ok := err.(net.Error)
	if !ok {
		fmt.Println("타입 변환이 실패했습니다.")
	} else {
		if !nErr.Timeout() {
			fmt.Printf("타임아웃이 발생하지 않았습니다: %v\n", err)
		}
	}
	if ctx.Err() != context.DeadlineExceeded {
		fmt.Printf("타임 아웃 : %v\n", ctx.Err())
	}
}

func DialContextCancel() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sync := make(chan struct{})

	go func() {
		defer func() {
			fmt.Println("고루틴 종료")
			sync <- struct{}{}
		}()
		var d net.Dialer
		d.Control = func(_, _ string, _ syscall.RawConn) error {
			time.Sleep(time.Second * 1)
			return nil
		}
		fmt.Println("start DialContext")
		conn, err := d.DialContext(ctx, "tcp", "127.0.0.1:3000")
		if err != nil {
			fmt.Println("DialContext error :", err)
			return
		}
		conn.Close()
		fmt.Println("connection did not time out")
	}()
	fmt.Println("start cancel")
	cancel()
	<-sync
	if ctx.Err() != context.Canceled {
		fmt.Printf("expected canceld context: actual: %q", ctx.Err())
	}
}
