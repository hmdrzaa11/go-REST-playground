package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// "HandleFunc" is going to take your "fn" and creates a "Handler" from it <Convert it>  and attach that with the "path"
	//to the "DefaultServeMux"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//read the body
		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Oops!", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "Hello %s\n", bs) //writes string into client
	})

	http.ListenAndServe("127.0.0.1:8000", nil) //if you do not pass a handler its going to use "defaultServeMux"

}
