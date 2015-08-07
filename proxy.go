package httpproxy

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

func New(port int) {
	portNum := strconv.Itoa(port)
	log.Printf("Starting proxy on port " + portNum)

	ln, err := net.Listen("tcp", ":"+portNum)

	if err != nil {
		log.Fatal("New proxy listen: ")
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("New proxy accept: ")
			log.Fatal(err)
		} else {
			go handleConnection(conn)
		}
	}
}

func handleConnection(c net.Conn) {
	buffer := make([]byte, 4096)

	read, err := c.Read(buffer)
	if err != nil {
		log.Fatal("Reading buffer: ")
		log.Fatal(err)
		return
	}

	if read <= 0 {
		log.Fatal("Read zero bytes")
		return
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
	log.Printf("LINE 0: " + request)
	tokens := strings.Split(request, " ")
	log.Printf(c.RemoteAddr().String() + " <-> " + tokens[1])
	client := &http.Client{}
	req, err := http.NewRequest(tokens[0], tokens[1], nil)

	if err != nil {
		log.Fatal("Request generation: ")
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
		log.Fatal("Client execute:")
		log.Fatal(err)
		return
	}

	content, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal("Content retrieval: ")
		log.Fatal(err)
		return
	}

	io.Copy(c, bytes.NewReader(content))

	c.Close()
}
