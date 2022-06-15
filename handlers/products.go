package handlers

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

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

	//PUT /
	if r.Method == http.MethodPut {
		path := r.URL.Path //gives us the "/<SOMETHING>" but we need to take out the "id"
		//use a regex
		regex := regexp.MustCompile(`/([0-9]+)`) //it it's an invalid regex its going to panic
		m := regex.FindAllStringSubmatch(path, -1)
		fmt.Println("Matched: ", m)
		if len(m) != 1 {
			//means we match a lot more OR nil and in this case is invalid
			http.Error(rw, "invalid URI", http.StatusBadRequest)
			return
		}

		if len(m[0]) != 2 {
			http.Error(rw, "invalid URI", http.StatusBadRequest)
			return
		}
		idString := m[0][1] //because the returned result of regex includes something like this " [[/12 12]]" we want the second index
		id, err := strconv.Atoi(idString)

		if err != nil {
			http.Error(rw, "invalid URI", http.StatusBadRequest)
			return
		}
		p.updateProducts(id, rw, r)
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

func (p *Products) updateProducts(id int, rw http.ResponseWriter, r *http.Request) {
	fmt.Println("Handle PUT /")
	//create the product obj
	prod := &data.Product{}
	err := prod.FromJson(r.Body) //call the method an this will fill the struct
	if err != nil {
		http.Error(rw, "failed to convert into json", http.StatusBadRequest)
	}

	err = data.UpdateProducts(id, prod)
	if err != nil {
		http.Error(rw, "product not found", http.StatusNotFound)
	}
}
