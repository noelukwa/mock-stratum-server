package server

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"gitlab.com/0xjonin/stratum/pkg/server/pg"
)

type Request struct {
	Sid    string        //unique id tied to every request
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
	Id     uint32        `json:"id"`
}

type Response struct {
	Id     uint32        `json:"id"`
	Result interface{}   `json:"result"`
	Error  []interface{} `json:"error"`
}

type Notification struct {
	Id     *string       `json:"id"`
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
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

	go h.startFakeJob()

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
		req.Sid = uuid.New().String()
		switch req.Method {
		case "mining.subscribe":
			h.handleSubscribe(req)
		case "mining.authorize":
			h.handleAuthorize(req)
		case "mining.notify":
			h.handleNotify(req)
		default:
			log.Printf("unknown method: %s", req.Method)
		}
	}
}

func (h *Handler) handleSubscribe(req Request) error {

	res := Response{
		Id: req.Id,
	}
	agent, ok := req.Params[0].(string)
	if !ok {
		return errors.New("invalid agent")
	}

	var err error

	subIdOne := randHex()
	subIdTwo := randHex()
	extraNonceOne := randHex()
	extraNonceTwo := 4
	res.Result = fmt.Sprintf(`[[["mining.set_difficulty", "%s"], ["mining.notify", "%s"]], "%s", %d]`, subIdOne, subIdTwo, extraNonceOne, extraNonceTwo)
	res.Error = nil

	err = h.store.SaveSubRequest(req.Method, agent, fmt.Sprintf("%d", req.Id), req.Sid, extraNonceOne)
	if err != nil {
		log.Println("sub save error: ", err)
		return err
	}

	err = h.respond(res)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) handleAuthorize(req Request) error {

	if len(req.Params) < 1 || len(req.Params) > 2 {
		return fmt.Errorf("invalid number of params")
	}
	var res Response
	var err error

	res.Id = req.Id
	res.Result = true
	res.Error = nil

	err = h.store.SaveAuthRequest(req.Method, req.Sid, fmt.Sprintf("%d", req.Id), req.Params[1].(string))
	if err != nil {
		log.Println("error saving auth request:", err)
		return err
	}
	err = h.respond(res)
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

func randHex() string {
	buf := make([]byte, 8)
	n := uint64(rand.Uint32())<<32 + uint64(rand.Uint32())
	binary.LittleEndian.PutUint64(buf, n)

	return hex.EncodeToString(buf)
}

func (h *Handler) startFakeJob() {

	defer func() {
		log.Println("fake job stopped")
	}()

	log.Println("starting fake job")
	for {
		time.Sleep(time.Second * 3)
		var notis Notification
		notis.Method = "mining.notify"
		notis.Id = nil
		notis.Params = []interface{}{
			randomUint16(),
			randomHash(),
			randomHash(),
			randomHash(),
			[]string{},
			"20000000",
			randHex(),
			randHex(),
			false,
		}
		msg, err := json.Marshal(notis)
		if err != nil {
			log.Println("error marshalling notification:", err)
			continue
		}
		msg = append(msg, '\n')
		h.conn.Write(msg)
	}

}

func randomUint16() uint16 {
	return uint16(rand.Uint32() & 0xffff)
}

func randomHash() string {
	buf := make([]byte, 32)
	rand.Read(buf)
	return hex.EncodeToString(buf)
}
