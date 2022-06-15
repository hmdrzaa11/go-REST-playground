package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/hmdrzaa11/micro-api/data"
)

type Products struct {
	l *log.Logger
}

//its doing the job of "contractor" fn in other languages
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	bs, err := json.Marshal(lp)
	if err != nil {
		http.Error(w, "unable to marshal", http.StatusInternalServerError)
		return
	}
	w.Write(bs)
}
