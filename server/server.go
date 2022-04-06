package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/awfulbits/wikiofthings/database"
	"github.com/gorilla/mux"
)

var port = os.Getenv("PORT")

func Start() {
	if port == "" {
		port = "8080"
	}

	db, err := database.Start()
	if err != nil {
		log.Fatalf("Database Failed To Start:%+v", err)
	}
	log.Print("Database Started")

	router := mux.NewRouter()
	router.HandleFunc("/title/{title}", titleHandler(db))
	http.Handle("/", router)

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")

	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		db.Stop()
		log.Print("Database Stopped")
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}
