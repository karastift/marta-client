// client

package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"time"
)

const marta = "127.0.0.1"
const port = 2222

var server *Server

func main() {

	server = NewServer(marta, port)

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

// Handles the data coming from the server.
func handleData(data string) {

	cmd, err := NewCommand(data)

	// if command is malformed, just answer with \n so the server knows that the client didnt time out
	if err != nil {
		log(err.Error())
		send("understood but not understood\n")
		return
	}

	switch cmd.CmdType {
	case "ping":
		ping(cmd)
	case "info":
		info(cmd)
	}
}

// Responds to "!ping" command.
func ping(cmd *Command) error {
	log("Got pinged by marta.")

	send("Pong\n")

	return nil
}

// Responds to "!info" command.
func info(cmd *Command) error {
	log("Info requested by marta.")

	info := NewInfo()

	sendBytes(append(info.Json(), '\n'))

	return nil
}

// Converts data to bytes and sends it to the server.
func sendBytes(data []byte) {
	server.Conn.Write(data)
}

// Converts data to bytes and sends it to the server.
func send(data string) {
	server.Conn.Write([]byte(data))
}

// Logs.
func log(str string) {
	fmt.Println("[" + time.Now().Format(time.ANSIC) + "] " + str)
}
