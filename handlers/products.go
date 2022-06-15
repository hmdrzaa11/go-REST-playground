package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hmdrzaa11/micro-api/data"
)

type Products struct {
	l *log.Logger
}

//its doing the job of "contractor" fn in other languages
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Get /")
	lp := data.GetProducts()
	err := lp.ToJson(rw)
	if err != nil {
		http.Error(rw, "failed to marshal", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST /")
	//create the product obj
	prod := &data.Product{}
	err := prod.FromJson(r.Body) //call the method an this will fill the struct
	if err != nil {
		http.Error(rw, "failed to convert into json", http.StatusBadRequest)
	}
	data.AddProducts(prod)

}

func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "invalid id", http.StatusBadRequest)
		return
	}
	p.l.Printf("Handle PUT /%d\n", id)
	//create the product obj
	prod := &data.Product{}
	err = prod.FromJson(r.Body) //call the method an this will fill the struct
	if err != nil {
		http.Error(rw, "failed to convert into json", http.StatusBadRequest)
	}

	err = data.UpdateProducts(id, prod)
	if err != nil {
		http.Error(rw, "product not found", http.StatusNotFound)
	}
}
