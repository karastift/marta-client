package main

import (
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

	mac, err := getMacAddr()

	if err != nil {
		panic(err)
	}

	info.MacAdress = mac[0]

	ip, err := getIP()

	if err != nil {
		panic(err)
	}

	info.LocalAddress = ip.String()

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

func getMacAddr() ([]string, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var as []string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}
	return as, nil
}

func getIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}
