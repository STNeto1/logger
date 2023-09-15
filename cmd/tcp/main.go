package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/stneto1/logger/pkg"
)

func main() {
	pkg.InitDB()

	log.Println("TCP Server")

	ln, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}

	msgChan := make(chan pkg.Message, 100_000)

	go handleMessageChannel(msgChan)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting connection", err)
			continue
		}

		go handleConnection(conn, msgChan)
	}
}

func handleConnection(conn net.Conn, msgChan chan pkg.Message) {
	defer conn.Close()

	// max message size 10kb
	buf := make([]byte, 10*1024)

	for {
		n, err := conn.Read(buf)

		if err != nil {
			if err == io.EOF {
				continue
			}

			fmt.Println("Error reading from connection", err)
			continue
		}

		dirtyMessage := string(buf[:n])
		tokens := strings.Split(dirtyMessage, "}{")

		for _, token := range tokens {
			correctToken := sanitizeToken(token)

			var payload pkg.Message
			if err := json.Unmarshal([]byte(correctToken), &payload); err != nil {
				fmt.Println("Error unmarshalling payload", err, correctToken)

				continue
			}

			msgChan <- payload
		}

	}
}

func handleMessageChannel(ch chan pkg.Message) {
	// queue := make([]pkg.Message, 0, 100)
	var queue []pkg.Message

	for {
		select {
		case msg := <-ch:
			if len(queue) == 100 {
				fmt.Println("batch size reached, sending to database")
				if err := pkg.DBCon.CreateMessages(queue); err != nil {
					fmt.Println("error on batch insert:", err)
				} else {
					fmt.Println("batch insert success")
				}

				queue = nil
			}

			queue = append(queue, msg)

		case <-time.After(1 * time.Second):
			fmt.Printf("pool size: %d\n", len(queue))
		}
	}
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
