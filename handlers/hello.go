package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Println("This is Hello Handler")
	//read the body
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Oops!", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Hello %s\n", bs) //writes string into client
}
