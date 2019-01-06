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
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	config.Init()

	url := fmt.Sprintf("%s:%s@/%s?charset=utf8mb4",
		config.Config.Database.Username,
		config.Config.Database.Password,
		config.Config.Database.Database)

	db.Init(url)

	r := mux.NewRouter()
	fileServer := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", removeDirectories(fileServer)))

	r.HandleFunc("/logout", api.Logout).Methods("GET")

	r.HandleFunc("/", pages.Index).Methods("GET")
	r.HandleFunc("/~{username}", pages.Profile).Methods("GET")
	r.HandleFunc("/posts/{post}.txt", pages.StaticPost).Methods("GET")
	r.HandleFunc("/posts/{post}", pages.Post).Methods("GET")
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

	options := handlers.AllowedOrigins([]string{"localhost", config.Config.Server.Origin})

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Config.Server.Port),
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Second * 30,
		Handler:      handlers.CORS(options)(r),
	}

	var err error
	if config.Config.Server.Port == 443 {
		go func() {
			http.ListenAndServe(":80", http.HandlerFunc(redirectTLS))
		}()
		srv.ListenAndServeTLS(config.Config.HTTPS.Certificate, config.Config.HTTPS.Key)
	} else {
		srv.ListenAndServe()
	}
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

func redirectTLS(w http.ResponseWriter, r *http.Request) {
	redirectURL := fmt.Sprintf("https://%s%s", r.Host, r.RequestURI)
	http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
}
