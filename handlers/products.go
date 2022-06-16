package handlers

import (
	"context"
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
	//we can now get the prod from the request
	prod := r.Context().Value(KeyProduct{}).(*data.Product) //because returns "any or interface{}" as result we need to cast it
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
	//now we get product from the request context
	prod := r.Context().Value(KeyProduct{}).(*data.Product) //here we are doing the conversion of the type to "*data.Product"

	err = data.UpdateProducts(id, prod)
	if err != nil {
		http.Error(rw, "product not found", http.StatusNotFound)
	}
}

//its better to use a struct as a key to your context with value instead of "string"
type KeyProduct struct{}

func (p *Products) ValidateProductMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) { //this converts the normal fn to "HandlerFunc"
		prod := &data.Product{}
		err := prod.FromJson(r.Body)
		if err != nil {
			http.Error(rw, "failed to convert into json", http.StatusBadRequest)
			return
		}

		//if everything went well we are going to add the "prod" into request context
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod) //here we adding to the "r.Context" SO DO NOT use a "context.Background"
		//because you will loos all of your request data
		//now we need a copy of request with the context we just created
		req := r.WithContext(ctx) //here we make a copy of the request we can use "r.Clone()" to get a deep copy
		//then we call next and pass the new request
		next.ServeHTTP(rw, req)
	})
}
