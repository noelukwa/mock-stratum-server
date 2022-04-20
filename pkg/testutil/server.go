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

	db, err := NewFakeDB()
	if err != nil {
		fmt.Printf("Error creating db: %v\n", err)
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			for {

				conn, err := ln.Accept()
				if err != nil {
					fmt.Printf("Error accepting connection: %v\n", err)
					return err
				}
				// fmt.Println("Accepted connection from", conn.RemoteAddr())
				handler := server.NewHandler(db, conn)
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

	_, err = db.Exec(
		` CREATE TABLE IF NOT EXISTS requests (
			method VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			req_id VARCHAR(255) NOT NULL,
			id VARCHAR(255) NOT NULL,
			user_id VARCHAR(255)
		  );
		  CREATE TABLE IF NOT EXISTS subscriptions (
			  req_id VARCHAR(255) NOT NULL,
			  method VARCHAR(255) NOT NULL,
			  id VARCHAR(255) NOT NULL,
			  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			  user_agent VARCHAR(255),
			  extra_nonce VARCHAR NOT NULL,
			  FOREIGN KEY (req_id) REFERENCES requests(req_id) ON DELETE CASCADE
		  );`,
	)

	if err != nil {
		return nil, err
	}

	return db, nil
}
