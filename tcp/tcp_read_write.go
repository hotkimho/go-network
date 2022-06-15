package tcp

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"io"
	"net"
	"reflect"
)

//뎅이터를 특정 버퍼만큼 쓰고, 읽어오기
func ReadIntoBuffer() {
	//20개의 Byte크기의 슬라이르 생성
	payload := make([]byte, 20)
	_, err := rand.Read(payload) //랜덤한 값으로 채웁니다
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("payload:", payload)
	//리스너 생성(localhost:3000)
	listener, err := net.Listen("tcp", "127.0.0.1:3000")
	if err != nil {
		fmt.Println(err)
		return
	}
	//고루틴에서 리스너가 요쳥을 수락합니다.
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		defer conn.Close()
		fmt.Println("연결이 성공했습니다", listener.Addr().String())
		_, err = conn.Write(payload)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("write payload")
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		fmt.Println(err)
		return
	}
	buf := make([]byte, 15)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		fmt.Println("읽은 값:", buf)
		fmt.Printf("read %d bytes\n", n)
	}
	conn.Close()
}

func ReadScanner() {
	//20개의 Byte크기의 슬라이르 생성
	const payload = "The bigger the interface, the weaker the abstraction."

	fmt.Println("payload:", payload)
	//리스너 생성(localhost:3000)
	listener, err := net.Listen("tcp", "127.0.0.1:3000")
	if err != nil {
		fmt.Println(err)
		return
	}
	//고루틴에서 리스너가 요쳥을 수락합니다.
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		defer conn.Close()
		fmt.Println("연결이 성공했습니다", listener.Addr().String())
		_, err = conn.Write([]byte(payload))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("write payload")
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	scanner.Split(bufio.ScanWords)

	var words []string

	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	err = scanner.Err()
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
	}
	result := []string{"The", "bigger", "the", "interface,", "the",
		"weaker", "the", "abstraction."}
	if !reflect.DeepEqual(words, result) {
		fmt.Println("두 문자열이 다릅니다.")
	}
	fmt.Printf("Scanned words: %#v\n", words)
}
