package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"site/api"
	"site/config"
	"site/db"
	"site/pages"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/shurcooL/httpgzip"
)

func main() {

	config.Init()

	url := fmt.Sprintf("%s:%s@/%s?charset=utf8mb4",
		config.Config.Database.Username,
		config.Config.Database.Password,
		config.Config.Database.Database)

	db.Init(url)

	r := mux.NewRouter()

	fileServer := httpgzip.FileServer(http.Dir("./static"), httpgzip.FileServerOptions{
		IndexHTML: false,
	})
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", cacheFiles(fileServer)))

	r.HandleFunc("/posts/{post}.txt", pages.StaticPost).Methods("GET")
	r.HandleFunc("/robots.txt", pages.Robots).Methods("GET")
	r.HandleFunc("/sitemap.xml", pages.Sitemap).Methods("GET")

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
	if len(config.Config.HTTPS.Certificate) > 0 {
		srv.ListenAndServeTLS(config.Config.HTTPS.Certificate, config.Config.HTTPS.Key)
	} else {
		srv.ListenAndServe()
	}
	log.Fatal(err)
}

func cacheFiles(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		next.ServeHTTP(w, r)
	})
}

func redirectTLS(w http.ResponseWriter, r *http.Request) {
	redirectURL := fmt.Sprintf("https://%s%s", r.Host, r.RequestURI)
	http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
}
