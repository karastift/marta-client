// client

package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/karastift/marta-client/world"
)

const marta = "127.0.0.1"
const port = 2222

var server *world.Server

func main() {
	server = world.NewServer(marta, port)

	loginToServer()

	err := listenToServer(handleData)

	if err != nil {
		panic(err)
	}
}

// Tries to log into server. Will try to reconnect every 5 seconds.
func loginToServer() net.Conn {
	for {

		log("Trying to log into marta.")

		conn, err := server.Login()

		if err != nil {
			log("Login failed. Retrying in 5 seconds.")
			time.Sleep(5 * time.Second)
			continue
		}

		log("Successfully logged into marta.")

		return conn
	}
}

// Listens to data from server and passes it into the `handleData` function.
func listenToServer(handleData func(data string)) error {

	if !server.IsLoggedIn() {
		panic(errors.New("failed listening to marta: not logged into marta"))
	}

	for {
		data, err := bufio.NewReader(server.Conn).ReadString('\n')

		if err != nil {
			// marta is offline
			// try to reconnect
			loginToServer()
			continue
		}

		go handleData(data)
	}
}

func handleData(data string) {
	log("Received: '" + data[:len(data)-1] + "'")
	// server.Conn.Write([]byte("you sent " + data + "\n"))
}

func log(str string) {
	fmt.Println("[" + time.Now().Format(time.ANSIC) + "] " + str)
}
