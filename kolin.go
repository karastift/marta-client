// client

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const Port = 2222
const cmdPrefix = "&"

var conn net.Conn

func main() {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", Port))

	// check if initiating the listener failed
	if err != nil {
		fmt.Println("Failed initiating the listener.")
		fmt.Println(err)
		os.Exit(1)
	}

	// defer .Close()
	// it will be executed at the end of the main function
	defer listener.Close()

	// main loop to check of incoming connections
	for {

		// accepting an incoming connection
		conn, err = listener.Accept()

		if err != nil {
			fmt.Println("Failed accepting a connection.")
			fmt.Println(err)
		}

		// read stream until \n
		netData, err := bufio.NewReader(conn).ReadString('\n')

		// if marta is ofline, a eof error will be thrown
		if err != nil {
			// wait 5 seconds and try to reconnect
			fmt.Println("Marta is offline. Trying to reconnect in 5 seconds.")
			time.Sleep(5 * time.Second)
		} else {
			// handle data
			handleData(netData)
		}

		// close connection
		conn.Close()
	}
}

func handleData(rawStr string) {
	// print rawStr without new line at the end
	fmt.Println("Received: '" + string(rawStr)[0:len(rawStr)-1] + "'")

	// if data is a command
	if string(rawStr[0]) == cmdPrefix {

		// remove first char and split rawStr after every " "
		data := strings.Split(rawStr[1:len(rawStr)-1], " ")
		cmd := data[0]

		fmt.Println("'" + cmd + "'")

		switch cmd {
		case "ping":
			handlePing()
		}
	}
}

func handlePing() {
	fmt.Println("Got pinged.")
	sendToConn("Pong.")
}

func sendToConn(str string) {
	_, err := conn.Write([]byte(str + "\n"))

	if err != nil {
		fmt.Println("Failed responding to ping command.")
		fmt.Println(err)
	}
}
