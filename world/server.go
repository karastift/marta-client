package world

import (
	"bufio"
	"errors"
	"fmt"
	"net"
)

type Server struct {
	Addr string
	Port int
	Conn net.Conn
}

func NewServer(addr string, port int) *Server {
	return &Server{
		Addr: addr,
		Port: port,
	}
}

func (server *Server) IsLoggedIn() bool {
	return server.Conn != nil
}

func (server *Server) Login() (net.Conn, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", server.Addr, server.Port))

	if err != nil {
		return nil, err
	}

	server.Conn = conn

	// send login code
	fmt.Fprintf(conn, "marta login\n")

	// receive login response
	status, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		return nil, errors.New("login failed: failed to receive login response")
	}

	// check login response
	if checkResponse(string(status)) {
		return conn, nil
	} else {
		return nil, errors.New("login failied: wrong login response")
	}
}

func (server *Server) Disconnect() {
	server.Conn.Close()
}

func checkResponse(res string) bool {
	return res == "marta logged in\n"
}
