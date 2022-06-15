package handlers

import (
	"fmt"
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

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	//GET /
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	//POST /
	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	//catch all
	rw.WriteHeader(http.StatusMethodNotAllowed)

}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJson(rw)
	if err != nil {
		http.Error(rw, "failed to marshal", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("Handle POST /")
	//create the product obj
	prod := &data.Product{}
	err := prod.FromJson(r.Body) //call the method an this will fill the struct
	if err != nil {
		http.Error(rw, "failed to convert into json", http.StatusBadRequest)
	}
	data.AddProducts(prod)

}
