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

	//2-we need to now block the code so its not going to exit so we are going to "read" from the channel
	sig := <-sigChan
	fmt.Println("Received termination signal, graceful shutdown : ", sig)
	//graceful shutdown : server not going to accept new request and waits till all the previous task done and then shuts down
	tc, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	srv.Shutdown(tc)
}
