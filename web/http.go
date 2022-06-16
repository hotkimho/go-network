package web

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
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

//요청 타임아웃 문제
func BlockIndefinitely() {
	http.HandleFunc("/", func(writer http.ResponseWriter, r *http.Request) {
		fmt.Println("Request Get")
		select {}
	})
	_ = http.ListenAndServe("127.0.0.1:3000", nil)
	fmt.Println("요청이 끝났습니다")
}

//타임 아웃 해결하기(deadline context)
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
	res.Close = true
}

type User struct {
	FirstName string
	LastName  string
}

func HandlePostUser() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(r io.ReadCloser) {
			_, _ = io.Copy(ioutil.Discard, r)
			_ = r.Close()
		}(r.Body)

		if r.Method != http.MethodPost {
			http.Error(w, "not Post", http.StatusMethodNotAllowed)
			return
		}
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Decode Failed", http.StatusBadRequest)
			return
		}
		fmt.Println(user)
		w.WriteHeader(http.StatusAccepted)
	}
}

func PostUser() {
	res, err := http.Get("localhost")
	if err != nil {
		fmt.Println("Get Failed")
		return
	}
	if res.StatusCode != http.StatusAccepted {
		fmt.Printf("받는 요청 : %d; actual status %d\n", http.StatusAccepted, res.StatusCode)
	}

	buf := new(bytes.Buffer)
	user := User{FirstName: "Ho", LastName: "KIM"}
	_ = json.NewEncoder(buf).Encode(&user)

	res, _ = http.Post("localhost", "application/json", buf)
	if res.StatusCode != http.StatusAccepted {
		fmt.Println("error")
	}
	res.Close = true
}

var t = template.Must(template.New("hello").Parse("hello, {{.}}!"))

func DefaultHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func(r io.ReadCloser) {
				_, _ = io.Copy(ioutil.Discard, r)
				_ = r.Close()
			}(r.Body)

			var b []byte
			switch r.Method {
			case http.MethodGet:
				b = []byte("hi!!!!")
			case http.MethodPost:
				var err error
				b, err = ioutil.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "Internal server error",
						http.StatusInternalServerError)
					return
				}
			default:
				http.Error(w, "Method net allowed", http.StatusMethodNotAllowed)
				return
			}
			_ = t.Execute(w, string(b))
		},
	)
}

type Handlers struct {
	db  *sql.DB
	log *log.Logger
}

func (h *Handlers) Handler1() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := h.db.Ping()
		})
}
func (h *Handlers) Handler2() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := h.db.Ping()
		})
}
