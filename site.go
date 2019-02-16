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

	r.HandleFunc("/posts/{post}.txt", pages.StaticPost).Methods("GET")

	a := r.PathPrefix("/api/").Subrouter()
	a.HandleFunc("/logout", api.Logout).Methods("POST")
	a.HandleFunc("/signup", api.Signup).Methods("POST")
	a.HandleFunc("/login", api.Login).Methods("POST")
	a.HandleFunc("/post", api.CreatePost).Methods("POST")
	a.HandleFunc("/comments", api.CreateComment).Methods("POST")
	a.HandleFunc("/replies", api.CreateReply).Methods("POST")
	a.HandleFunc("/vote/{post}", api.Vote).Methods("POST")
	a.HandleFunc("/posts/{post}", api.Post).Methods("GET")
	a.HandleFunc("/browse/{page}", api.Browse).Methods("GET")
	a.HandleFunc("/conversation", api.Conversation).Methods("GET")
	a.HandleFunc("/users/{user}", api.User).Methods("GET")
	a.HandleFunc("/isLoggedIn", api.IsLoggedIn).Methods("GET")

	r.PathPrefix("/").HandlerFunc(pages.Index)

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
			http.Error(w, "Not found", 404)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func redirectTLS(w http.ResponseWriter, r *http.Request) {
	redirectURL := fmt.Sprintf("https://%s%s", r.Host, r.RequestURI)
	http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
}
