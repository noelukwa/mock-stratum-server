package server_test

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"testing"
	"time"

	"gitlab.com/0xjonin/stratum/pkg/server"
	"gitlab.com/0xjonin/stratum/pkg/testutil"
)

func TestAuthorize(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fmt.Printf("Starting server\n")
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	go testutil.NewTCP(ctx, ln)
	time.Sleep(time.Second * 3)
	testAuthorize(t, ln.Addr().String())
}

func testAuthorize(t *testing.T, addr string) {
	fmt.Printf("Connecting to server\n")
	client, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}

	req := server.Request{
		Method: "mining.authorize",
		Params: []interface{}{"test", "test"},
		Id:     "1",
	}

	rawmsg, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}

	rawmsg = append(rawmsg, '\n')
	fmt.Printf("Sending request\n")
	_, err = client.Write(rawmsg)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Reading response\n")
	for {
		fmt.Printf("Creating NewReader for socket\n")
		reader := bufio.NewReader(client)

		fmt.Printf("Reading from socket\n")
		rawmessage, err := reader.ReadString('\n')
		if err != nil {
			t.Fatal(err)
		}

		r := server.Response{}
		err = json.Unmarshal([]byte(rawmessage), &r)
		if err != nil {
			fmt.Println(string(rawmessage))
			t.Fatal(err)
		}

		if r.Id == "1" {
			break
		} else {
			t.Fatal("got a wrong message: ", rawmessage)
		}
	}
}
