package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"github.com/gorilla/websocket"
	"os/signal"
	"io"
)


func main() {

	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("\n  Usage: wscat <<ws:// or wss:// URL>>")
		os.Exit(-1)
	}


	addr := os.Args[len(os.Args)-1]
	if !strings.HasPrefix(addr, "ws://") && !strings.HasPrefix(addr, "wss://") {
		addr = "ws://"+addr
	}

	c, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	defer c.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)


	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				fmt.Printf("<<server: %s>>\n", err)
				close(interrupt)
				return
			}

			io.WriteString(os.Stdout,fmt.Sprintf("%s\n",message))
		}
	}()


  // Handle program interrupt
	select {
	case <-interrupt:
		c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure,""))
		c.Close()
		os.Exit(0)
	}


	}