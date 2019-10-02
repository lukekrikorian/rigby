package pages

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"site/db"
	"site/routes"
	"strings"
	text "text/template"

	"github.com/gorilla/mux"
)

var (
	indexTemplate   = template.Must(template.ParseFiles("static/template.html")).Lookup("index")
	sitemapTemplate = text.Must(text.ParseFiles("static/sitemap.xml")).Lookup("sitemap")
)

type indexInject struct {
	Description string
	Script      string
	Styles      string
	Title       string
}

// Index page
func Index(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path

	var (
		code        int
		description string
		styles      string
		script      string
		title       string = "Rigby"
	)

	for _, route := range routes.Routes {
		if route.Path == p {
			code = 200
			description = route.Description
			title = route.Title
			break
		} else if route.RegexpPath != nil {
			matches := route.RegexpPath.FindStringSubmatch(p)
			if matches != nil {
				code, description, title = route.Matcher(matches)
				break
			}
		}
	}

	if code == 0 {
		code, description, title = 404, "Not found!", "404 Error"
	}

	files, _ := ioutil.ReadDir("static/dist")
	for _, file := range files {
		if n := file.Name(); strings.HasPrefix(n, "main") {
			script = n
		} else {
			styles = n
		}
	}

	w.WriteHeader(code)
	indexTemplate.Execute(w, indexInject{
		Description: description,
		Script:      script,
		Styles:      styles,
		Title:       title,
	})
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

// Sitemap is the sitemap.xml page
func Sitemap(w http.ResponseWriter, r *http.Request) {
	var paths []string

	for _, route := range routes.Routes {
		if route.Path != "" {
			paths = append(paths, route.Path)
		}
	}

	w.Header().Set("Content-Type", "application/xml; charset=UTF-8")
	sitemapTemplate.Execute(w, paths)
}
