package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"site/api"
	"site/config"
	"site/db"
	"site/pages"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	fileServer := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", removeDirectories(fileServer)))

	r.HandleFunc("/logout", api.Logout).Methods("GET")

	r.HandleFunc("/", pages.Index).Methods("GET")
	r.HandleFunc("/~{username}", pages.Profile).Methods("GET")
	r.HandleFunc("/post/{post}.txt", pages.StaticPost).Methods("GET")
	r.HandleFunc("/post/{post}", pages.Post).Methods("GET")
	r.HandleFunc("/comment/{post}", pages.Comment).Methods("GET")
	r.HandleFunc("/reply/{comment}", pages.Reply).Methods("GET")
	r.HandleFunc("/browse/{page}", pages.Browse).Methods("GET")
	r.HandleFunc("/conversation", pages.Conversation).Methods("GET")
	r.HandleFunc("/{page}", pages.Pages).Methods("GET")

	r.HandleFunc("/api/signup", api.Signup).Methods("POST")
	r.HandleFunc("/api/login", api.Login).Methods("POST")
	r.HandleFunc("/api/post", api.CreatePost).Methods("POST")
	r.HandleFunc("/api/comments", api.CreateComment).Methods("POST")
	r.HandleFunc("/api/replies", api.CreateReply).Methods("POST")
	r.HandleFunc("/api/vote/{post}", api.Vote).Methods("POST")
	r.HandleFunc("/api/posts/{post}", api.Post).Methods("GET")
	r.HandleFunc("/api/browse/{page}", api.Browse).Methods("GET")
	r.HandleFunc("/api/conversation", api.Conversation).Methods("GET")

	options := handlers.AllowedOrigins([]string{"localhost", config.Config.Server.Origin})

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Config.Server.Port),
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Second * 30,
		Handler:      handlers.CORS(options)(r),
	}

	go func() {
		err := srv.ListenAndServe()
		log.Fatal(err)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	db.SaveSessions()
	srv.Shutdown(context.Background())
}

func removeDirectories(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			pages.NotFound.Execute(w, nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}
