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
	Matcher     func(matches []string) (Status int, Description string)
	Description string
}

// Routes is an exported list of routs
var Routes = []Route{
	Route{
		Path:        "/",
		Description: "Rigby is a micoblogging platform for cool kids",
	},
	Route{
		Path:        "/login",
		Description: "Log into rigby",
	},
	Route{
		Path:        "/signup",
		Description: "Sign up for rigby",
	},
	Route{
		Path:        "/browse/popular",
		Description: "Browse popular rigby posts",
	},
	Route{
		Path:        "/browse/recent",
		Description: "Browse recent rigby posts",
	},
	Route{
		Path:        "/conversation",
		Description: "Browse recent conversation",
	},
	Route{
		Path:        "/post",
		Description: "Create a rigby post",
	},
	Route{
		RegexpPath: regexp.MustCompile("/~([a-zA-Z0-9_]+)$"),
		Matcher: func(matches []string) (Status int, Description string) {
			username := matches[1]
			user, err := db.GetUserByUsername(username)
			if err == nil {
				Status = 200
				Description = fmt.Sprintf("%s's rigby profile", user.Username)
			}
			return
		},
	},
	Route{
		RegexpPath: regexp.MustCompile("/posts/([a-zA-Z0-9-]+$)"),
		Matcher: func(matches []string) (Status int, Description string) {
			postID := matches[1]
			post, err := db.GetPost(postID)
			if err == nil {
				Status = 200
				Description = fmt.Sprintf("%.300s...", post.Body)
			}
			return
		},
	},
}
