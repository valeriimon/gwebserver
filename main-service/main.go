package mainservice

import (
	"log"
	"net"
)

func InitService(address string) {
	if lst, err := net.Listen("tcp", address); err == nil {
		log.Printf("Main Service Listening On: %s", lst.Addr())
	} else {
		log.Fatalln(err.Error())
	}
}
