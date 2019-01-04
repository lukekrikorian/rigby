package pages

import (
	"fmt"
	"html/template"
	"net/http"
	"site/db"

	"github.com/gorilla/mux"
)

var templates = template.Must(template.ParseGlob("templates/*"))

// NotFound is the 404 page page
var NotFound = templates.Lookup("404")

// Profile page
func Profile(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	user, err := db.GetUserByUsername(username)
	if err != nil {
		fmt.Println(err)
		NotFound.Execute(w, nil)
		return
	}

	err = user.GetPosts()
	if err != nil {
		fmt.Println(err)
		NotFound.Execute(w, nil)
		return
	}

	templates.ExecuteTemplate(w, "profile", user)
}

// Index page
func Index(w http.ResponseWriter, r *http.Request) {
	if userID, err := db.CheckAuth(r); err == nil {
		if userInfo, err := db.GetUserByID(userID); err == nil {
			templates.ExecuteTemplate(w, "/", userInfo)
			return
		}
	}
	templates.ExecuteTemplate(w, "/", nil)
}

// Post page
func Post(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["post"]
	post, err := db.GetPost(postID)
	if err == nil {
		post.GetComments()

		for i := range post.Comments {
			post.Comments[i].GetReplies()
		}

		templates.ExecuteTemplate(w, "posts", post)
	} else {
		fmt.Println(err)
		NotFound.Execute(w, nil)
	}
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

// Comment creation page
func Comment(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["post"]
	post, err := db.GetPost(postID)
	if err == nil {
		templates.ExecuteTemplate(w, "comment", post)
	} else {
		NotFound.Execute(w, nil)
	}
}

// Reply creation page
func Reply(w http.ResponseWriter, r *http.Request) {
	parentID := mux.Vars(r)["comment"]
	parent, err := db.GetComment(parentID)
	if err == nil {
		templates.ExecuteTemplate(w, "reply", parent)
	} else {
		NotFound.Execute(w, nil)
	}
}

// Pages handers all the "static" pages
func Pages(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["page"]
	if page := templates.Lookup(path); page != nil {
		page.Execute(w, nil)
	} else {
		NotFound.Execute(w, nil)
	}
}

// Recent posts page
func Recent(w http.ResponseWriter, r *http.Request) {
	posts, err := db.GetRecentPosts()
	if err != nil {
		NotFound.Execute(w, nil)
		return
	}
	templates.ExecuteTemplate(w, "browse", posts)
}
