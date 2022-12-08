package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
)

const LISTEN_PORT = 46374

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

func main() {
	http.HandleFunc("/", content)
	http.ListenAndServe(":"+strconv.Itoa(LISTEN_PORT), nil)
}

func main3() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	fmt.Println("new client")

	proxy, err := net.Dial("tcp", "127.0.0.1:80")
	if err != nil {
		panic(err)
	}

	fmt.Println("proxy connected")
	go copyIO(conn, proxy)
	go copyIO(proxy, conn)
}

func copyIO(src, dest net.Conn) {
	defer src.Close()
	defer dest.Close()
	io.Copy(src, dest)
}
