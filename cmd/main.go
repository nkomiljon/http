package main

import (
	"time"
	"io/ioutil"
	"fmt"
	"bytes"
	"strings"
	"strconv"
	"io"
	"os"
	"net"
	"log"

)

func main()  {
	host := "0.0.0.0"
	port := "9999"
	if err := execute(host, port); err != nil {
		os.Exit(1)
	}
	log.Print("server closed")
}

func execute(host, port string) (err error)  {
	listener, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		log.Print(err)
		return
	}

	defer func()  {
		if cerr := listener.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Print(cerr)
		}
	}()

for {
	conn, err := listener.Accept()
	 if err != nil {
		 log.Print(err)
		 continue 
	 }
	
	err = handle(conn)
	err = conn.Close()
	if err != nil {
		log.Print(err)
		continue
	}
}

	return
}

func handle(conn net.Conn) (err error)  {
	defer func()  {
		if cerr := conn.Close();
		cerr != nil {
			if err == nil {
				err = cerr 
				return
			}
			log.Print(cerr)
		}
	}()

    
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != io.EOF {
		log.Printf("EOF")
	}
	if err != nil {
		return
	}

	data := buf[:n]
	requestLineDelim := []byte{'\r', '\n'}
	requestLineEnd  := bytes.Index(data, requestLineDelim)
    if requestLineEnd == -1 {
		return 
	}
	requestLine := string(data[:requestLineEnd])
	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		return
	}
	method, path, version := parts[0], parts[1], parts[2]
	if method != "GET" {
		return
	}
	if version != "HTTP/1.1" {
		return
	}

	if path == "/" {
		body, err := ioutil.ReadFile("static/index.html")
		if err != nil {
			return fmt.Errorf("can't open index.htmml: %w,", err)
		}
			marker := "{{year}}"
			year := time.Now().Year()
			body = bytes.ReplaceAll(body, []byte(marker), []byte(strconv.Itoa(year)))

			_,err = conn.Write([]byte(
			"HTTP/1.1 200 OK\r\n"+
			"Content-Length: " +strconv.Itoa(len(body))+"\r\n"+
			"Content-Type: text/html\r\n"+
			"Connection: close\r\n"+
			"\r\n"+
			string(body),
		))
		if err != nil {
			return err
		}
	}
	return nil
 }
