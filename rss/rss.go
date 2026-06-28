package rss

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"text/template"
	"time"

	"site/config"
	"site/db"
	"site/pages"

	md "github.com/russross/blackfriday/v2"
)

type Feed struct {
	Title        string
	URL          string
	FeedImageURL string
	Description  string
	Date         time.Time
	Items        []Item
}

type Item struct {
	ID     string
	Title  string
	Author string
	URL    string
	Body   string
	Date   time.Time
}

var templates = template.Must(template.New("temp.xml").Funcs(pages.FuncMap).ParseGlob("rss/*.xml"))

// RecentPostsRSS responds with an XML RSS feed of /browse/recent
func RecentPostsRSS(w http.ResponseWriter, r *http.Request) {
	var (
		posts []db.Post
		err   error
	)

	posts, err = db.GetRecentPosts()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not load recent posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/rss+xml; charset=UTF-8")

	items := []Item{}

	for _, post := range posts {
		body := md.Run([]byte(post.Body), md.WithExtensions(pages.MarkdownOptions))

		item := Item{
			ID:     post.ID,
			Title:  html.EscapeString(post.Title),
			Author: html.EscapeString(post.Author),
			URL:    fmt.Sprintf("%s/post/%s", config.Config.RSS.BaseURL, post.ID),
			Body:   string(body),
			Date:   post.Created,
		}

		items = append(items, item)
	}

	feed := Feed{
		Title:        "Rigby",
		URL:          fmt.Sprintf("%s/browse/recent", config.Config.RSS.BaseURL),
		FeedImageURL: fmt.Sprintf("%s/static/images/icon.png", config.Config.RSS.BaseURL),
		Description:  "posts as they are created",
		Date:         time.Now(),
		Items:        items,
	}

	err = templates.ExecuteTemplate(w, "feed", feed)
	if err != nil {
		log.Fatal(err)
	}
}
