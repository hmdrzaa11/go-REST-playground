package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hmdrzaa11/micro-api/handlers"
)

func main() {
	l := log.New(os.Stdout, "REST-api", log.LstdFlags)
	hh := handlers.NewHello(l)
	gh := handlers.NewGoodbye(l)
	servMux := http.NewServeMux() //now we are going to create a new servMux then register all of our handlers into it

	servMux.Handle("/", hh) //pass the type Hello and because of "ServeHTTP" it qualifies as "Handler"
	servMux.Handle("/goodbye", gh)
	//http.ListenAndServe(":8000", servMux) this is going to give us a basic server its better to control the timeouts
	//to manage our resource better
	srv := &http.Server{
		Addr:         ":8000",
		Handler:      servMux,
		IdleTimeout:  time.Second * 120,
		ReadTimeout:  time.Second * 1,
		WriteTimeout: time.Second * 1,
	}
	srv.ListenAndServe()
}
