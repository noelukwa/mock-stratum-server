package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	name = "basic/4.0.0"
	user = "local.test"
)

type Request struct {
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
	Id     int           `json:"id"`
}

type Response struct {
	Id     int         `json:"id"`
	Result interface{} `json:"result"`
	Error  interface{} `json:"error"`
}

func main() {
	server := flag.String("server", "127.0.0.1:8080", "stratum server address")
	flag.Parse()

	if *server == "" {
		log.Fatalln("server address is required")
	}

	c, err := net.Dial("tcp", *server)
	if err != nil {
		log.Fatalln(err)
	}
	defer c.Close()

	log.Println("connected to", c.LocalAddr())
	if err := sendAuthRequest(c); err != nil {
		log.Fatalln(err)
	}

	if err := sendSubRequest(c); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("now listening for notifications")

out:
	for {

		select {
		case <-time.After(time.Second * 10):
			break out
		default:
			reader := bufio.NewReader(c)
			rawmessage, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalln(err)
			}
			log.Printf("RECEIVED: %+s\n", rawmessage)

		}
	}
}

func sendAuthRequest(c net.Conn) error {
	log.Println("sending auth request")
	var req Request
	req.Method = "mining.authorize"
	req.Params = []interface{}{user, ""}
	req.Id = 2

	raw, err := json.Marshal(req)
	if err != nil {
		return err
	}
	raw = append(raw, '\n')
	_, err = c.Write(raw)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(c)
	rawmessage, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	log.Printf("RECEIVED: %+s\n", rawmessage)
	return nil
}

func sendSubRequest(c net.Conn) error {
	log.Println("sending subscribe request")
	var req Request
	req.Method = "mining.subscribe"
	req.Params = []interface{}{name}
	req.Id = 2

	raw, err := json.Marshal(req)
	if err != nil {
		return err
	}
	raw = append(raw, '\n')
	_, err = c.Write(raw)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(c)
	rawmessage, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	log.Printf("RECEIVED: %+s\n", rawmessage)
	return nil
}
