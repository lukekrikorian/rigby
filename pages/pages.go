package pages

import (
	"fmt"
	"net/http"
	"site/db"

	"github.com/gorilla/mux"
)

// Index page
func Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/dist/index.html")
}

// StaticPost handles static post pages like /post/{id}.txt
func StaticPost(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["post"]
	post, err := db.GetPost(postID)
	if err == nil {
		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		w.Write([]byte(post.Body))
	} else {
		fmt.Println(err)
		http.Error(w, "Error serving static file", 500)
	}
}

// Robots is the robots.txt page
func Robots(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/robots.txt")
}
