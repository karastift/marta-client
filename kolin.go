// client

package main

import (
	"bufio"
	"errors"
	"fmt"
	"time"

	"github.com/karastift/marta-client/world"
)

const marta = "127.0.0.1"
const port = 1111

func main() {
	server := world.NewServer(marta, port)

	_, err := server.Login()

	if err != nil {
		panic(errors.New("failed logging into marta: marta is offline"))
	}
	log("Successfully logged into marta.")

	err = listenTo(server, handleData)

	if err != nil {
		panic(err)
	}
}

func listenTo(server *world.Server, handleData func(data string)) error {

	if !server.IsLoggedIn() {
		panic(errors.New("failed listening to marta: not logged into marta"))
	}

	for {
		data, err := bufio.NewReader(server.Conn).ReadString('\n')

		if err != nil {
			// TODO: maybe add reconnection logic later
			return errors.New("failed listening to marta: marta is offline")
		}

		go handleData(data)
	}
}

func handleData(data string) {
	log("Received: '" + data[:len(data)-1] + "'")
}

func log(str string) {
	fmt.Println("[" + time.Now().Format(time.ANSIC) + "] " + str)
}
