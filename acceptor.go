package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var quit = make(chan struct{})

func main() {
	fromPort, toPort, die := processFlags()
	listeners := createHttpListeners(fromPort, toPort, die)
	startHttpServers(listeners)
	handlePosix(listeners)
}

// Process the flags on load
func processFlags() (int64, int64, bool) {
	fromPortFlag := flag.Int("from", 1024, "From port")
	tilPortFlag := flag.Int("until", 2048, "Until port")
	dieFlag := flag.Bool("die", true, "Die when port cannot be opened")
	flag.Parse()

	fromPort := *fromPortFlag
	tilPort := *tilPortFlag
	die := *dieFlag

	if fromPort < 1 {
		log.Panicf("From port must be between 1 and 65535")
	}

	if tilPort > 65535 {
		log.Panicf("Until port must be between 1 and 65535")
	}

	return int64(fromPort), int64(tilPort), die
}

// Builds a slice of http listeners
func createHttpListeners(fromPort int64, toPort int64, die bool) []net.Listener {
	listeners := make([]net.Listener, toPort-fromPort+1)

	for listenPort := fromPort; listenPort <= toPort; listenPort++ {
		addr := ":" + strconv.FormatInt(int64(listenPort), 10)
		log.Printf("trying to bind %s", addr)
		l, err := net.Listen("tcp", addr)
		if err != nil {
			log.Printf("can't create listener number {%d} - %+v", listenPort, err)
			if die == true {
				os.Exit(1)
			}

			continue
		}

		listeners = append(listeners, l)
	}

	return listeners
}

// Starts http servers inside goroutines based on listeners
func startHttpServers(listeners []net.Listener) {
	http.HandleFunc("/", content)
	for _, l := range listeners {
		if l != nil {
			go http.Serve(l, nil)
		}
	}
}

// Serves a basic html page
func content(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "<!doctype html>")
	fmt.Fprintln(w, "<html lang=en>")
	fmt.Fprintln(w, "<head>")
	fmt.Fprintln(w, "<meta charset=utf-8>")
	fmt.Fprintln(w, "</head>")
	fmt.Fprintln(w, "<body>")
	fmt.Fprintf(w, "<p>Successful connection to `%v`</p>\n", req.Host)
	fmt.Fprintln(w, "<p>Headers:</p>")
	fmt.Fprintln(w, "<pre>")

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}

	fmt.Fprintf(w, "</pre>")

	fmt.Fprintf(w, "</body>")
	fmt.Fprintf(w, "</html>")
}

// Close the listeners
func closeListeners(listeners []net.Listener) {
	log.Printf("closing listeners...")
	for _, l := range listeners {
		if l != nil {
			if err := l.Close(); err != nil {
				log.Printf("can't close the listener:%+v", err)
			}
		}
	}
}

// https://gobyexample.com/signals
func handlePosix(listeners []net.Listener) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		closeListeners(listeners)
		done <- true
	}()

	fmt.Println("awaiting signal (e.g ctrl-c)")
	<-done
	fmt.Println("exiting")
}
