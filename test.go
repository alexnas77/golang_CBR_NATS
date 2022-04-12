package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
)

func main() {

	nc, _ := nats.Connect(nats.DefaultURL)
	ec, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	defer ec.Close()
	defer nc.Close()

	type person struct {
		Name    string
		Address string
		Age     int
	}

	recvCh := make(chan *person)
	_, err := ec.BindRecvChan("hello", recvCh)
	if err != nil {
		return
	}

	sendCh := make(chan *person)
	err = ec.BindSendChan("hello", sendCh)
	if err != nil {
		return
	}

	me := &person{Name: "derek", Age: 22, Address: "140 New Montgomery Street"}

	// Send via Go channels
	sendCh <- me

	// Receive via Go channels
	who := <-recvCh
	fmt.Print("who: ")
	fmt.Println(who)
}
