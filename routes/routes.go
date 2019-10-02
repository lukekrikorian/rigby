package routes

import (
	"fmt"
	"regexp"
	"site/db"
)

/*
	Routes are essentially just my way of generating page descriptions
	for crawlers that don't process JS and correct HTTP status codes.
	Static paths like /login and /signup only have a string path and
	description, but dynamic paths like post pages have dynamic
	descriptions and status codes to match.
*/

// Route represents a route, or URL
type Route struct {
	Path        string
	RegexpPath  *regexp.Regexp
	Matcher     func(matches []string) (Status int, Description string, Title string)
	Description string
	Title       string
}

// Routes is an exported list of routs
var Routes = []Route{
	Route{
		Path:        "/",
		Title:       "Rigby",
		Description: "Rigby is a micoblogging platform for cool kids",
	},
	Route{
		Path:        "/login",
		Title:       "Login - Rigby",
		Description: "Log into rigby",
	},
	Route{
		Path:        "/signup",
		Title:       "Signup - Rigby",
		Description: "Sign up for rigby",
	},
	Route{
		Path:        "/browse/popular",
		Title:       "Popular Posts - Rigby",
		Description: "Browse popular rigby posts",
	},
	Route{
		Path:        "/browse/recent",
		Title:       "Recent Posts - Rigby",
		Description: "Browse recent rigby posts",
	},
	Route{
		Path:        "/conversation",
		Title:       "Conversations - Rigby",
		Description: "Browse recent conversation",
	},
	Route{
		Path:        "/post",
		Title:       "New Post - Rigby",
		Description: "Create a rigby post",
	},
	Route{
		RegexpPath: regexp.MustCompile("/~([a-zA-Z0-9_]+)$"),
		Matcher: func(matches []string) (Status int, Description string, Title string) {
			username := matches[1]
			user, err := db.GetUserByUsername(username)
			if err == nil {
				Status = 200
				Description = fmt.Sprintf("The rigby profile for %s", user.Username)
				Title = fmt.Sprintf("~%s - Rigby", user.Username)
			}
			return
		},
	},
	Route{
		RegexpPath: regexp.MustCompile("/posts/([a-zA-Z0-9-]+$)"),
		Matcher: func(matches []string) (Status int, Description string, Title string) {
			postID := matches[1]
			post, err := db.GetPost(postID)
			if err == nil {
				Status = 200
				Description = fmt.Sprintf("%.300s...", post.Body)
				Title = fmt.Sprintf("%s by %s", post.Title, post.Author)
			}
			return
		},
	},
}
