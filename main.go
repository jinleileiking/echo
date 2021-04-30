package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	// "github.com/cloudflare/tableflip"
)

func echo(w http.ResponseWriter, req *http.Request) {

	fmt.Printf("incoming URL: %#v\n ", req.URL)
	if t, ok := req.URL.Query()["t"]; ok {
		if tInt, err := strconv.Atoi(t[0]); err != nil {
			panic(err)
		} else {
			time.Sleep(time.Duration(tInt) * time.Second)
		}
	} else {
		time.Sleep(time.Duration(*sleep) * time.Second)
	}

	k, _ := req.URL.Query()["k"]
	fmt.Fprintf(w, fmt.Sprintf("%s", k[0]))
}

var sleep = flag.Int("sleep", 5, "Sleep time")
var isServer = flag.Bool("server", true, "Use as server")
var url = flag.String("url", "http://localhost:8094/echo?k=1", "server url")
var port = flag.Int("port", 8094, "server port")

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	flag.Parse()

	if *isServer {
		log.Println("Runing server")
		log.Printf("Sleep: %d", *sleep)
		server()
	} else {
		if *url == "" {
			log.Println("need param url")
			return
		}
		log.Println("Runing client")
		client()
	}
}

func server() {

	mux := http.NewServeMux()

	mux.HandleFunc("/echo", echo)

	srv := http.Server{Addr: ":8094", Handler: mux}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT)
		<-sigint

		// We received an interrupt signal, shut down.
		log.Println("Begin exit")
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed

	// http.HandleFunc("/echo", echo)
	// if err := http.ListenAndServe(":8094", nil); err != nil {
	// 	panic(err)
	// }

	// upg, _ := tableflip.New(tableflip.Options{})
	// defer upg.Stop()

	// go func() {
	// 	sig := make(chan os.Signal, 1)

	// 	// ctrl c
	// 	signal.Notify(sig, syscall.SIGINT)
	// 	for range sig {
	// 		log.Println("Got Signal")
	// 		upg.Upgrade()
	// 	}
	// }()

	// // Listen must be called before Ready
	// ln, err := upg.Listen("tcp", ":8094")

	// if err != nil {
	// 	panic(err)
	// }

	// defer ln.Close()

	// go http.Serve(ln, nil)

	// if err := upg.Ready(); err != nil {
	// 	panic(err)
	// }

	// <-upg.Exit()

	log.Println("server exited")

	return
}

func client() {

	// var netTransport = &http.Transport{
	// 	Dial: (&net.Dialer{
	// 		Timeout: 5 * time.Second,
	// 	}).Dial,
	// 	TLSHandshakeTimeout: 500 * time.Second,
	// }
	// var netClient = &http.Client{
	// 	Timeout:   time.Second * 10,
	// 	Transport: netTransport,
	// }

	// Xs send 100 requests
	for {
		for i := 0; i < 100; i++ {
			go func() {
				// response, err := netClient.Get(*url)

				// default uses connection pool and persistConn
				response, err := http.Get(*url)
				if err != nil {
					panic(err)
				}

				buf, err := ioutil.ReadAll(response.Body)
				if err != nil {
					panic(err)
				}

				log.Printf("Ret:%s\n", string(buf))
			}()
		}
		log.Println("Send done... Begin sleep")
		time.Sleep(time.Duration(1) * time.Second)
		log.Println("Sleep done")
	}
	return
}
