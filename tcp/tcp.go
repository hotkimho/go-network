package tcp

import (
	"fmt"
	"io"
	"net"
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
	listener, err := net.Listen("tc", "127.0.0.1:")
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
