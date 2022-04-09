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
	"github.com/awfulbits/wikiofthings/testpage"
	"github.com/gorilla/mux"
)

func Start() {
	port := os.Getenv("WOT_PORT")
	if port == "" {
		port = "8080"
	}

	db, err := database.Start("pages", "pageidindex")
	if err != nil {
		log.Panicf("database failed to start:%+v\n", err)
	}
	log.Print("Database Started")

	router := mux.NewRouter()
	router.HandleFunc("/title/{title}", titleHandler(db))
	// router.HandleFunc("/edit/", editHandler(db))
	// router.HandleFunc("/edit/{pageId}", editHandler(db))
	http.Handle("/", router)

	srv := &http.Server{
		Handler:      router,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panicf("listen: %s\n", err)
		}
	}()
	log.Printf("Server Started On Port %v", port)

	// Test page creation
	log.Print("Running Test Page")
	if err = testpage.RunTest(db); err != nil {
		log.Panicf("cannot create test page: %s\n", err)
	}
	log.Print("Test Page Ran Successfully, See /title/Hello_friend")

	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		db.Stop()
		log.Print("Database Stopped")
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Panicf("server shutdown failed:%+v\n", err)
	}
	log.Print("Server Exited Properly")
}
