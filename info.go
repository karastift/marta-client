package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"os/user"
	"runtime"
	"strings"
)

type Info struct {
	Username      string
	Name          string
	Uid           string
	HomeDir       string
	Os            string
	Device        string
	MacAdress     string
	Administrator bool
	LocalAddress  string
}

// Returns a pointer to an instance of Info with all information stored inside.
func NewInfo() *Info {

	info := Info{}

	user, err := user.Current()

	if err != nil {
		panic(err)
	}

	info.Username = user.Username
	info.Name = user.Name
	info.Uid = user.Uid
	info.HomeDir = user.HomeDir

	info.Os = strings.Title(runtime.GOOS)

	info.LocalAddress = getIP()
	info.MacAdress = getMacAddr()

	return &info
}

func (info *Info) Encode() string {
	return base64.StdEncoding.EncodeToString([]byte(info.String()))
}

func (info *Info) Json() []byte {
	j, err := json.Marshal(info)

	if err != nil {
		panic(err)
	}

	return j
}

func (info *Info) String() string {
	return fmt.Sprintf(`Username	%s
Name		%s
Uid		%s
HomeDir		%s
Os		%s
Device		%s
MacAdress	%s
Administrator	%t
LocalAddress	%s`, info.Username, info.Name, info.Uid, info.HomeDir, info.Os, info.Device, info.MacAdress, info.Administrator, info.LocalAddress)
}

func getIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")

	if err != nil {
		logger.Println("Failed to get IP.")
		return ""
	}

	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.String()
}

func getMacAddr() (addr string) {
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, i := range interfaces {
			if i.Flags&net.FlagUp != 0 && !bytes.Equal(i.HardwareAddr, nil) {
				// Don't use random as we have a real address
				addr = i.HardwareAddr.String()
				break
			}
		}
	}
	return
}
