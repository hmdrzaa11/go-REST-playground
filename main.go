package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hmdrzaa11/micro-api/handlers"
)

func main() {
	l := log.New(os.Stdout, "REST-api: ", log.LstdFlags)
	ph := handlers.NewProducts(l)
	servMux := mux.NewRouter() //now we are to use the gorilla mux as our mux

	//CORS handler
	corsHandler := gorillaHandlers.CORS(gorillaHandlers.AllowedOrigins([]string{"http://localhost:3000"}))

	//create a sub route for GET methods
	getRouter := servMux.Methods(http.MethodGet).Subrouter() //we are separating our app in sub routes based on METHODS its going
	getRouter.HandleFunc("/", ph.GetProducts)                //to help us for better middleware implementation

	putRouter := servMux.Methods(http.MethodPut).Subrouter()
	putRouter.Use(ph.ValidateProductMiddleware)
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProducts) //here we specifying that id is of type number
	//we can add a custom middleware to this subrouter

	postRouter := servMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)
	postRouter.Use(ph.ValidateProductMiddleware) //also applied the middleware in here as well
	//to manage our resource better
	srv := &http.Server{
		Addr:         ":8000",
		Handler:      corsHandler(servMux), //wrap the serveMux with "CorsHandler"
		ErrorLog:     l,                    //sets the error logger for the server
		IdleTimeout:  time.Second * 120,
		ReadTimeout:  time.Second * 1,
		WriteTimeout: time.Second * 1,
	}
	go func() {
		err := srv.ListenAndServe() //as we now this line is going to block the main thread so we need to put it inside of a go routine
		if err != nil {
			l.Fatal(err)
		}
	}()
	//we used a "go routine" and this will make sure the code is "non-blocking" BUT also now its going to immediately exit
	//SOLUTION: use signals

	//1-create a channel of type signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt) //broadcast a message on the "sigChan" when the "os.Interrupt" event happens
	signal.Notify(sigChan, syscall.SIGTERM)

	log.Printf("listening on port :8000")
	//2-we need to now block the code so its not going to exit so we are going to "read" from the channel
	sig := <-sigChan
	fmt.Println("Received termination signal, graceful shutdown : ", sig)
	//graceful shutdown : server not going to accept new request and waits till all the previous task done and then shuts down
	tc, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	srv.Shutdown(tc)
}
