package main

import (
	"bufio"
	"log"
	"net/url"
	"os"
	"sync"
	"time"

	"os/signal"

	"github.com/gorilla/websocket"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	destList := []string{
		"pd_tacaruna",
		"sjc",
	}

	var wg sync.WaitGroup
	for _, dest := range destList {
		wg.Add(1)
		go func(dest string) {
			defer wg.Done()

			u := url.URL{Scheme: "ws", Host: "localhost:8082", Path: "/ws"}
			log.Printf("Connecting to %s", u.String())

			c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
			if err != nil {
				log.Printf("dial error for %s: %v", dest, err)
				return
			}
			defer c.Close()

			file, err := os.Open("cmd/car_simulator/" + dest + ".txt")
			if err != nil {
				log.Printf("Error opening file for %s: %v", dest, err)
				return
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				select {

				case <-interrupt:
					log.Println("Interrupt received, closing connection")
					err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
					if err != nil {
						log.Printf("write close error for %s: %v", dest, err)
					}
					return
				default:
					err := c.WriteMessage(websocket.TextMessage, []byte(scanner.Text()))
					if err != nil {
						log.Printf("write error for %s: %v", dest, err)
						return
					}
					time.Sleep(time.Second) // control the rate of messages
				}
			}
			if err := scanner.Err(); err != nil {
				log.Printf("Error reading data for %s: %v", dest, err)
			}
		}(dest)
	}

	wg.Wait()
	log.Println("All connections closed")
}
