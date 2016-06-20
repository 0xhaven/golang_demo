package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// TODO: Move to subpackge and export/comment approriate members.

type request interface{}

type writeReq struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type readReq struct {
	Name  string `json:"name"`
	valCh chan<- string
}

type server struct {
	reqChan chan request
	values  map[string]string
}

func (s *server) serve() {
	for req := range s.reqChan {
		switch req := req.(type) {
		case *writeReq:
			s.values[req.Name] = req.Value
		case *readReq:
			req.valCh <- s.values[req.Name]
		}
	}
}

func startServer() (*writeHandler, *readHandler) {
	s := &server{make(chan request), make(map[string]string)}
	go s.serve() // TODO: add ability to stop server/close reqChan.
	return &writeHandler{s.reqChan}, &readHandler{s.reqChan}
}

type writeHandler struct {
	reqChan chan request
}

func (h *writeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	req := new(writeReq)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.reqChan <- req
}

type readHandler struct {
	reqChan chan request
}

func (h *readHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	req := new(readReq)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	valCh := make(chan string)
	req.valCh = valCh
	h.reqChan <- req
	fmt.Fprint(w, <-valCh)
}
