package main

import (
	"fmt"
	"net/http"
	"time"
)

func echo(w http.ResponseWriter, req *http.Request) {

	fmt.Printf("incoming URL: %#v\n ", req.URL)
	time.Sleep(time.Duration(60) * time.Second)
	k, _ := req.URL.Query()["k"]
	fmt.Fprintf(w, fmt.Sprintf("%s\n", k[0]))
}

func main() {

	http.HandleFunc("/echo", echo)

	if err := http.ListenAndServe(":8094", nil); err != nil {
		panic(err)
	}
}
