package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gitlab.com/0xjonin/stratum/cmd/data"
	"gitlab.com/0xjonin/stratum/pkg/server"
)

func main() {
	port := flag.String("port", "8080", "port to listen on")

	dsn := flag.String("dsn", "postgresql://localhost:5432/luxor?user=uche&password=uchechukwu&sslmode=disable", "database connection string")

	flag.Parse()

	if err := data.RunMigration(*dsn); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

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
		handler.HandleRequests()
	}

}
