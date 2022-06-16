package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

//시간을 가져와 현재시간과 비교
func HandTime() {
	res, err := http.Head("https://www.time.gov")
	if err != err {
		fmt.Println(err)
	}
	_ = res.Body.Close()
	fmt.Println("now:", time.Now())
	now := time.Now().Round(time.Second)
	fmt.Println("after now:", now)
	date := res.Header.Get("Date")
	if date == "" {
		fmt.Println("no Date Header receivedd from time")
	}

	dt, err := time.Parse(time.RFC1123, date)
	fmt.Println("dt:", dt)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("time.gov: %s (skew %s)\n", dt, now.Sub(dt))
}

func BlockIndefinitely() {
	http.HandleFunc("/", func(writer http.ResponseWriter, r *http.Request) {
		fmt.Println("Request Get")
		select {}
	})
	_ = http.ListenAndServe("127.0.0.1:3000", nil)
	fmt.Println("요청이 끝났습니다")
}

func BlockIndefinitelyTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/", nil)
	if err != nil {
		fmt.Println(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if !errors.Is(err, context.DeadlineExceeded) {
			fmt.Println(err)
		}
		return
	}
	res.Body.Close()
}
