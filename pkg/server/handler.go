package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"

	"github.com/jmoiron/sqlx"
	"gitlab.com/0xjonin/stratum/pkg/server/pg"
)

type Request struct {
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
	Id     string        `json:"id"`
}

type Response struct {
	Id     string        `json:"id"`
	Result interface{}   `json:"result"`
	Error  []interface{} `json:"error"`
}

type Handler struct {
	store *pg.Store
	conn  net.Conn
}

func NewHandler(db *sqlx.DB, conn net.Conn) *Handler {
	store := pg.NewStore(db)
	return &Handler{
		store: store,
		conn:  conn,
	}
}

func (h *Handler) HandleRequests() {

	reader := bufio.NewReader(h.conn)
	for {
		msg, err := reader.ReadBytes('\n')
		if err != nil {
			return
		}
		var req Request
		err = json.Unmarshal(msg, &req)
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		switch req.Method {
		case "mining.subscribe":
			h.handleSubscribe(req)
		case "mining.authorize":
			h.handleAuthorize(req)
		case "mining.notify":
			h.handleNotify(req)
		default:
			fmt.Println("error:", err)
		}
	}
}

func (h *Handler) handleSubscribe(req Request) error {

	return nil
}

func (h *Handler) handleAuthorize(req Request) error {

	if len(req.Params) < 1 || len(req.Params) > 2 {
		return fmt.Errorf("invalid number of params")
	}
	var res Response

	res.Id = req.Id
	res.Result = true
	res.Error = nil

	err := h.respond(res)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) handleNotify(req Request) error {

	return nil
}

func (h *Handler) respond(res Response) error {
	msg, err := json.Marshal(res)
	if err != nil {
		return err
	}
	msg = append(msg, '\n')
	_, err = h.conn.Write(msg)
	return err
}
