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

			msg := pkg.Message{
				Topic: "test",
				Body:  json.RawMessage(`"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Fusce dapibus nunc nec ullamcorper iaculis. Fusce elit libero, cursus eget luctus at, maximus sed nunc. Nunc et tincidunt mi, non semper lacus. Donec pretium placerat risus. Nulla ornare velit nec orci imperdiet aliquet ut eget ligula. Aliquam elementum ipsum id magna tempor, elementum commodo odio rutrum. Aliquam erat volutpat. Ut posuere interdum turpis, nec blandit sapien semper at. Vivamus scelerisque, dolor eget auctor ultricies, eros augue dignissim orci, at iaculis nunc mi sed purus. Suspendisse potenti."`),
			}

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
			time.Sleep(time.Millisecond)

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

	msg := pkg.Message{
		Topic: "test",
		Body:  json.RawMessage(`"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Fusce dapibus nunc nec ullamcorper iaculis. Fusce elit libero, cursus eget luctus at, maximus sed nunc. Nunc et tincidunt mi, non semper lacus. Donec pretium placerat risus. Nulla ornare velit nec orci imperdiet aliquet ut eget ligula. Aliquam elementum ipsum id magna tempor, elementum commodo odio rutrum. Aliquam erat volutpat. Ut posuere interdum turpis, nec blandit sapien semper at. Vivamus scelerisque, dolor eget auctor ultricies, eros augue dignissim orci, at iaculis nunc mi sed purus. Suspendisse potenti."`),
	}

	// struct to bytes
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Panicln("error on serialization:", err)
		return
	}

	for i := 0; i < qty; i++ {
		go func(group *sync.WaitGroup, idx int) {
			defer group.Done()

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
