package httpproxy

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
)

func New(port int) {
	log.Printf("Starting proxy on port " + string(port))

	ln, err := net.Listen("tcp", ":"+string(port))

	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	buffer := make([]byte, 4096)

	_, err := c.Read(buffer)
	if err != nil {
		log.Fatal("Error reading bytes")
	}

	var lines []string

	j := 0
	for i := 0; i < len(buffer); i++ {
		if buffer[i] == '\n' {
			lines = append(lines, string(buffer[j:i-1]))
			j = i + 1
		}
	}

	request := lines[0]
	tokens := strings.Split(request, " ")
	log.Printf(c.RemoteAddr().String() + " -> " + tokens[1])
	client := &http.Client{}
	req, err := http.NewRequest(tokens[0], tokens[1], nil)

	if err != nil {
		log.Fatal(err)
		return
	}

	for i := 1; i < len(lines); i++ {
		if len(lines[i]) > 2 {
			headerTokens := strings.Split(lines[i], ": ")
			req.Header.Add(headerTokens[0], headerTokens[1])
		}
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
		return
	}

	content, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
		return
	}

	io.Copy(c, bytes.NewReader(content))

	c.Close()
}