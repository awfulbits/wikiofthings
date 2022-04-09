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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db, err := database.Start("pages", "pageidindex")
	if err != nil {
		log.Fatalf("Database Failed To Start:%+v\n", err)
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
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")

	// Test page creation
	if err = testpage.RunTest(db); err != nil {
		log.Fatalf("cannot create test page: %s\n", err)
	}
	log.Print("Test page ran successfully, see /title/Hello_friend")

	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		db.Stop()
		log.Print("Database Stopped")
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v\n", err)
	}
	log.Print("Server Exited Properly")
}
