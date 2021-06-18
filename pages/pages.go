package pages

import (
	"fmt"
	"html/template"
	"net/http"
	"site/db"
	"time"

	"github.com/gorilla/mux"
	md "github.com/russross/blackfriday/v2"
)

var templates = template.Must(template.New("temp.html").Funcs(funcMap).ParseGlob("templates/[^.]*"))

var options = md.NoIntraEmphasis | md.FencedCode | md.Autolink | md.Strikethrough | md.HeadingIDs

var funcMap = template.FuncMap{
	"prettyTime": func(t time.Time) string {
		return t.Format("2006-01-02")
	},
	"markdown": func(s string) template.HTML {
		body := md.Run([]byte(s), md.WithExtensions(options))
		return template.HTML(string(body))
	},
}

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

		templates.ExecuteTemplate(w, "post", post)
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

// Browse posts, either popular or recent
func Browse(w http.ResponseWriter, r *http.Request) {
	page := mux.Vars(r)["page"]
	var (
		posts []db.Post
		err   error
	)
	switch page {
	case "recent":
		posts, err = db.GetRecentPosts()
	case "popular":
		posts, err = db.GetPopularPosts()
	default:
		NotFound.Execute(w, nil)
		return
	}
	if err != nil {
		fmt.Println(err)
		NotFound.Execute(w, nil)
		return
	}
	templates.ExecuteTemplate(w, "postsList", posts)
}

// Conversation page
func Conversation(w http.ResponseWriter, r *http.Request) {
	comments, err := db.GetRecentComments()
	if err != nil {
		fmt.Println(err)
		NotFound.Execute(w, nil)
		return
	}

	for i := range comments {
		if err := comments[i].GetParent(); err != nil {
			fmt.Println(err)
		}
		if err = comments[i].GetReplies(); err != nil {
			fmt.Println(err)
		}
	}

	templates.ExecuteTemplate(w, "conversation", comments)
}
