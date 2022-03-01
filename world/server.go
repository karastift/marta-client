package world

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

type Server struct {
	Addr    string
	Port    int
	Conn    net.Conn
	pausing bool
}

func NewServer(addr string, port int) *Server {
	return &Server{
		Addr:    addr,
		Port:    port,
		pausing: false,
	}
}

func (server *Server) Send(data []byte) ([]byte, error) {

	server.Conn.Write(data)

	responseData, err := bufio.NewReader(server.Conn).ReadBytes('\n')

	return []byte(responseData), err
}

func (server *Server) Listen(handleData func(data string)) {
	for {
		if server.pausing {
			time.Sleep(5 * time.Second)
			continue
		}
		data, err := bufio.NewReader(server.Conn).ReadString('\n')

		if err != nil {
			fmt.Println("Marta is offline.")
			break
		}

		go handleData(data)
	}
}

func (server *Server) PauseListening() {
	server.pausing = true
}

func (server *Server) ResumeListening() {
	server.pausing = false
}

func (server *Server) IsConnected() bool {
	return server.Conn != nil
}

func (server *Server) Connect() (net.Conn, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", server.Addr, server.Port))

	if err != nil {
		return nil, err
	}
	server.Conn = conn

	response, err := server.Send([]byte("marta login\n"))

	if err != nil {
		fmt.Println("Login failed.")
	}

	if checkResponse(string(response)) {
		fmt.Println("[" + time.Now().Format(time.ANSIC) + "] Successfully connected to marta.")

	}

	return conn, err
}

func (server *Server) Disconnect() {
	server.Conn.Close()
}

func checkResponse(res string) bool {
	return res == "marta logged in\n"
}
