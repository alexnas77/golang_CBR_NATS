package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

func subscribeWait() {
	nc, _ := nats.Connect(nats.DefaultURL)
	fmt.Println("NATS connected")
	ec, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	//defer ec.Close()
	//defer nc.Close()

	recvCh := make(chan interface{})
	_, err := ec.BindRecvChan("hello", recvCh)
	if err != nil {
		return
	}
	fmt.Println("NATS chanel bind")

	// Receive via Go channels
	who := <-recvCh
	log.Print("who: ")
	fmt.Println(who)
	ec.Close()
	nc.Close()
}

func main() {
	for true {
		subscribeWait()
		time.Sleep(10 * time.Second)
	}
}
