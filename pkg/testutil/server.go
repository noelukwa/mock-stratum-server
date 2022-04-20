package testutil

import (
	"context"
	"fmt"
	"net"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"gitlab.com/0xjonin/stratum/pkg/server"
)

func NewTCP(ctx context.Context, ln net.Listener) error {
	fmt.Println("Listening on port", ln.Addr().(*net.TCPAddr).Port)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			for {
				fmt.Println("Waiting for connection")
				conn, err := ln.Accept()
				if err != nil {
					fmt.Printf("Error accepting connection: %v\n", err)
					return err
				}
				fmt.Println("Accepted connection from", conn.RemoteAddr())
				handler := server.NewHandler(nil, conn)
				handler.HandleRequests()
			}

		}
	}

}

func NewFakeDB() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	return db, nil
}
