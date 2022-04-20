package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gitlab.com/0xjonin/stratum/pkg/server"
)

func main() {
	port := flag.String("port", "8080", "port to listen on")

	dsn := flag.String("dsn", "postgresql://localhost:5434/luxor?user=luxor&password=luxor&sslmode=disable", "database connection string")

	flag.Parse()

	dbConn, err := sqlx.Connect("postgres", *dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", *port))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Listening on port", *port)
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		handler := server.NewHandler(dbConn, conn)
		if err != nil {
			log.Fatal(err)
		}
		go handler.HandleRequests()
	}

}
