// client

package main

import (
	"fmt"
	"time"

	"github.com/karastift/marta-client/world"
)

const marta = "127.0.0.1"
const port = 2222

var server *world.Server

func main() {
	server = world.NewServer(marta, port)

	_, err := server.Connect()

	if err != nil {
		fmt.Println("Marta is offline.")
		return
	}

	server.Listen(handleData)
}

func handleData(data string) {
	fmt.Println("[" + time.Now().Format(time.ANSIC) + "]" + " Received: '" + data[:len(data)-1] + "'")

	server.Conn.Write([]byte("Beat me out of me\n"))

	fmt.Println("Wrote")
}
