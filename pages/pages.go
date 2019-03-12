package pages

import (
	"fmt"
	"html/template"
	"net/http"
	"site/db"
	"site/routes"

	"github.com/gorilla/mux"
)

var indexTemplate = template.Must(template.ParseFiles("static/dist/index.html")).Lookup("index")

// Index page
func Index(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path

	var (
		code        int
		description string
	)

	for _, route := range routes.Routes {
		if route.Path == p {
			code = 200
			description = route.Description
			break
		} else if route.RegexpPath != nil {
			matches := route.RegexpPath.FindStringSubmatch(p)
			if matches != nil {
				code, description = route.Matcher(matches)
			}
		}
	}

	if code == 0 {
		code, description = 404, "Not found!"
	}

	w.WriteHeader(code)
	indexTemplate.Execute(w, description)
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
