package tcp

import (
	"fmt"
	"net"
)

func listener() {
	//리스너에 IP와 포트 번호가 바인딩됨
	listener, err := net.Listen("tcp", "127.0.0.1:3000")
	if err != nil {
		fmt.Println("listener(): Error net.Listen func")
	}

	defer func() { _ = listener.Close() }()

	fmt.Printf("bound to %q", listener.Addr())
}
