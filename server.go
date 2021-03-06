package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"time"
)

type Server struct {
	Addr string
	Port int
	Conn net.Conn
}

// Returns pointer to an instance of Server.
func NewServer(addr string, port int) *Server {
	return &Server{
		Addr: addr,
		Port: port,
	}
}

// Checks if connection is not `nil`.
func (server *Server) IsLoggedIn() bool {
	return server.Conn != nil
}

// Logs into server.
func (server *Server) Login() (net.Conn, error) {
	// connect to server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", server.Addr, server.Port))

	if err != nil {
		return nil, err
	}

	server.Conn = conn

	// send login code and info about client

	info := NewInfo()

	conn.Write([]byte("marta login|" + string(info.Json()) + "\n"))

	// only wait 5 seconds for marta to responds
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	// receive login response
	status, err := bufio.NewReader(conn).ReadString('\n')

	// reset deadline, because client should listen (listenToServer() in kolin.go)
	// time.Time{} is apperently a "zero" value
	conn.SetReadDeadline(time.Time{})

	if err != nil {
		return nil, errors.New("login failed: failed to receive login response")
	}

	// check login response
	if checkResponse(string(status)) {
		return conn, nil
	} else {
		return nil, errors.New("login failed: wrong login response")
	}
}

// Disconnect from server and close connection.
func (server *Server) Disconnect() {
	server.Conn.Close()
}

// Check login response from server.
func checkResponse(res string) bool {
	return res == "marta logged in\n"
}
