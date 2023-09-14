package main

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"strings"

	"github.com/stneto1/logger/pkg"
)

func main() {
	log.Println("TCP Server")

	ln, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting connection", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// max message size 10kb
	buf := make([]byte, 10*1024)

	for {
		n, err := conn.Read(buf)

		if err != nil {
			if err == io.EOF {
				continue
			}

			log.Println("Error reading from connection", err)
			continue
		}

		dirtyMessage := string(buf[:n])
		tokens := strings.Split(dirtyMessage, "}{")

		for _, token := range tokens {
			correctToken := sanitizeToken(token)

			var payload pkg.Message
			if err := json.Unmarshal([]byte(correctToken), &payload); err != nil {
				log.Println("Error unmarshalling payload", err)

				continue
			}

			log.Printf("Message received: %s", correctToken)
		}

	}

	// for {
	// 	buf := make([]byte, 1024)
	// 	n, err := conn.Read(buf)
	//
	// 	if err != nil {
	// 		if err == io.EOF {
	// 			continue
	// 		}
	//
	// 		log.Println("Error reading from connection", err)
	// 		continue
	// 	}
	//
	// 	var payload pkg.Message
	// 	if err := json.Unmarshal(buf[:n], &payload); err != nil {
	// 		log.Println("Error unmarshalling payload", err)
	// 		continue
	// 	}
	//
	// 	log.Printf("Message received: %s", buf[:n])
	// }
}

func sanitizeToken(token string) string {
	if !strings.HasPrefix(token, "{") {
		token = "{" + token
	}

	if !strings.HasSuffix(token, "}") {
		token = token + "}"
	}

	return token
}
