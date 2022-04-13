package main

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func getCurrency() {
	req, err := http.NewRequest("GET", "https://www.cbr-xml-daily.ru/daily_json.js", nil)
	if err != nil {
		// handle error
		log.Fatal(err.Error())
	}

	filename := "Last-Modified.txt"

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Println("File " + filename + " does not exist")
		f, err := os.Create(filename)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				log.Fatal(err.Error())
			}
		}(f)
	} else {
		b, err := ioutil.ReadFile(filename)
		// can file be opened?
		if err != nil {
			log.Fatal(err.Error())
		}
		if len(string(b)) > 0 {
			req.Header.Set("If-Modified-Since", string(b))
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// handle error
		log.Fatal(err.Error())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}(resp.Body)

	fmt.Println("Status: " + resp.Status)
	fmt.Println("StatusCode: " + strconv.Itoa(resp.StatusCode))

	if resp.StatusCode == 304 {
		return
	}

	/*var headers map[string]string
	headers = make(map[string]string)
	for key := range resp.Header {
		//fmt.Println(key + " : " + resp.Header.Get(key))
		headers[key] = resp.Header.Get(key)
	}
	h, err := json.MarshalIndent(headers, " ", " ")
	if err != nil {
		// handle error
		log.Fatal(err.Error())
	}
	fmt.Println(string(h))*/

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}(file)
	_, err = file.WriteString(resp.Header.Get("Last-Modified"))
	if err != nil {
		// handle error
		log.Fatal(err.Error())
	}
	var result map[string]map[string]interface{}

	_ = json.NewDecoder(resp.Body).Decode(&result)
	/*if err != nil {
		// handle error
		log.Fatal(err.Error())
	}*/
	for valute := range result["Valute"] {
		if valute != "USD" && valute != "EUR" {
			delete(result["Valute"], valute)
			result["Valute"]["now"] = time.Now()
		}
	}

	b, err := json.MarshalIndent(result["Valute"], " ", " ")
	if err != nil {
		// handle error
		log.Fatal(err.Error())
	}
	fmt.Println(string(b))

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		// handle error
		log.Fatal(err.Error())
	}
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		// handle error
		log.Fatal(err.Error())
	}
	/*recvCh := make(chan interface{})
	_, err = ec.BindRecvChan("hello", recvCh)
	if err != nil {
		return
	}*/

	sendCh := make(chan interface{})
	err = ec.BindSendChan("hello", sendCh)
	if err != nil {
		return
	}
	me := result["Valute"]
	// Send via Go channels
	sendCh <- me

	/*// Receive via Go channels
	who := <-recvCh
	fmt.Print("who: ")
	fmt.Println(who)*/

	//ec.Close()

	//nc.Close()

}

func main() {
	for true {
		getCurrency()
		time.Sleep(120 * time.Second)
	}

}
