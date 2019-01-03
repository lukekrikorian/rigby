package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"site/api"
	"site/config"
	"site/db"
	"site/pages"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {

	config.Init()

	url := fmt.Sprintf("%s:%s@/%s",
		config.Config.Database.Username,
		config.Config.Database.Password,
		config.Config.Database.Database)

	db.Init(url)

	r := mux.NewRouter()
	fileServer := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", removeDirectories(fileServer)))

	r.HandleFunc("/", pages.Index).Methods("GET")
	r.HandleFunc("/~{username}", pages.Profile).Methods("GET")
	r.HandleFunc("/posts/{post}.txt", pages.StaticPost).Methods("GET")
	r.HandleFunc("/posts/{post}", pages.Post).Methods("GET")
	r.HandleFunc("/comment/{post}", pages.Comment).Methods("GET")
	r.HandleFunc("/reply/{comment}", pages.Reply).Methods("GET")
	r.HandleFunc("/browse/recent", pages.Recent).Methods("GET")
	r.HandleFunc("/{page}", pages.Pages).Methods("GET")

	r.HandleFunc("/api/signup", api.Signup).Methods("POST")
	r.HandleFunc("/api/login", api.Login).Methods("POST")
	r.HandleFunc("/api/post", api.CreatePost).Methods("POST")
	r.HandleFunc("/api/comments", api.CreateComment).Methods("POST")
	r.HandleFunc("/api/replies", api.CreateReply).Methods("POST")

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Config.Server.Port),
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Second * 30,
		Handler:      r,
	}
	err := srv.ListenAndServe()
	log.Fatal(err)
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
