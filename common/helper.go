package common

import (
	"crypto/tls"
	"flag"
	"log"
	"net"
	"net/http"
	"time"
)

func MakeHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

func RecoverFromPanic() any {
	if err := recover(); err != nil {
		log.Printf("Recovered from panic: %v\n", err)
		return err
	}

	return ""
}

func GetFreePort() uint {
	var l *net.TCPListener
	for {
		addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
		if err != nil {
			return GetFreePort()
		} else {
			l, err = net.ListenTCP("tcp", addr)
			if err != nil {
				return GetFreePort()
			}
			defer l.Close()

			break
		}
	}

	return uint(l.Addr().(*net.TCPAddr).Port)
}

func GetFreePortsLength() int {
	var freePorts []uint
	for {
		freePort := GetFreePort()
		for _, port := range freePorts {
			if port == freePort {
				return len(freePorts)
			}
		}
		freePorts = append(freePorts, freePort)
	}
}

func IsTest() bool {
	return flag.Lookup("test.v") != nil
}
