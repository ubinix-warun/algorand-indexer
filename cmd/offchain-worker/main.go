package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"encoding/json"

	"github.com/algorand/go-algorand/data/transactions"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
// Token
// KEY

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, byt, err := c.ReadMessage()
			if err != nil {
				log.Println("read: ", err)
				return
			}

			var dat transactions.ApplyData

			if err := json.Unmarshal(byt, &dat); err != nil {
				log.Println("unmarshal: ", err)
				// log.Printf(" >>> recv: %s", string(byt))
			}

			log.Printf("recv: %s", string(byt))
			log.Println("ApplicationID: ", dat.ApplicationID)
			log.Println("EvalDelta: ", dat.EvalDelta)
			
			// FIX: ADDRESS! < PrivateKey

			// globalState["reqm"]
			// globalState["requrl"]


			// globalState["resp"] < Result

			// 
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}