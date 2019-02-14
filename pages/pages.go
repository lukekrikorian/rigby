package pages

import (
	"fmt"
	"net/http"
	"site/db"

	"github.com/gorilla/mux"
)

// Index page
func Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/app.html")
}

// StaticPost handles static post pages like /post/{id}.txt
func StaticPost(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["post"]
	post, err := db.GetPost(postID)
	if err == nil {
		postString := fmt.Sprintf("%s by %s\n\n%s", post.Title, post.Author, post.Body)
		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		w.Write([]byte(postString))
	} else {
		fmt.Println(err)
		http.Error(w, "Error serving static file", 500)
	}
}
