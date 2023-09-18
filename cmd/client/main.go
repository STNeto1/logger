package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/stneto1/logger/pkg"
)

func main() {
	mode := flag.String("mode", "tcp", "mode to run the client")
	msgs := flag.Int("msgs", 1, "number of messages to send")
	flag.Parse()

	log.Println("Client -", *mode, *msgs)

	switch *mode {
	case "tcp":
		sendTcpMessage(*msgs)
	case "http":
		sendHttpMessage(*msgs)
	}
}

func createMessage() pkg.Message {
	var topic string
	err := faker.FakeData(&topic)
	if err != nil {
		log.Fatalln("error creating fake data:", err)
	}

	var data map[string]string
	err = faker.FakeData(&data)
	if err != nil {
		log.Fatalln("error creating fake data:", err)
	}

	msg, err := json.Marshal(data)
	if err != nil {
		log.Fatalln("error marshaling data:", err)
	}

	return pkg.Message{
		Topic: topic,
		Body:  json.RawMessage(msg),
	}
}

func sendTcpMessage(qty int) {
	var wg sync.WaitGroup
	wg.Add(qty)

	for i := 0; i < qty; i++ {
		go func(group *sync.WaitGroup, idx int) {
			defer group.Done()
			conn, err := net.Dial("tcp", "127.0.0.1:1234")

			if err != nil {
				log.Println("error dialing connection:", err)
			}

			time.Sleep(time.Millisecond * (50 + time.Duration(idx)))

			msg := createMessage()
			// struct to bytes
			msgBytes, err := json.Marshal(msg)
			if err != nil {
				log.Println("error on serialization:", err)
				return
			}

			if _, err := conn.Write(msgBytes); err != nil {
				log.Println("error writing to connection:", err)
				return
			}

			log.Printf("message %d sent\n", idx)

			if err := conn.Close(); err != nil {
				log.Println("error closing connection:", err)
			}
		}(&wg, i+1)
	}

	wg.Wait()
}

func sendHttpMessage(qty int) {
	var wg sync.WaitGroup
	wg.Add(qty)

	for i := 0; i < qty; i++ {
		go func(group *sync.WaitGroup, idx int) {
			defer group.Done()

			msg := createMessage()
			// struct to bytes
			msgBytes, err := json.Marshal(msg)
			if err != nil {
				log.Panicln("error on serialization:", err)
				return
			}

			time.Sleep(time.Millisecond * (50 + time.Duration(idx)))

			res, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/log", bytes.NewReader(msgBytes))
			if err != nil {
				log.Println("error creating request:", err)
				return
			}
			res.Header.Set("Content-Type", "application/json")

			client := http.Client{}
			if _, err := client.Do(res); err != nil {
				log.Println("error sending request:", err)
				return
			}

			log.Printf("message %d sent\n", idx)

			// body, _ := io.ReadAll(res.Body)
			// log.Println(string(body))
			// time.Sleep(time.Millisecond)

		}(&wg, i+1)
	}

	wg.Wait()
}
