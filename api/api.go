package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"site/db"
	"site/pages"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/gorilla/schema"
	uuid "github.com/satori/go.uuid"
)

var signupRegex = regexp.MustCompile("^[a-zA-Z0-9_]+$")
var decoder = schema.NewDecoder()

func fail(w http.ResponseWriter, m string) {
	w.WriteHeader(500)
	pages.Failure.Execute(w, m)
}

func auth(w http.ResponseWriter, r *http.Request) (UserID string, Error error) {
	if UserID, Error = db.CheckAuth(r); Error != nil {
		fail(w, "You aren't logged in.")
	}
	return
}

// Signup endpoint (/api/signup POST)
func Signup(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username, password := r.FormValue("username"), r.FormValue("password")

	if l := len(username); l < 3 || l > 16 {
		fail(w, "Your username must be within 3 and 16 characters long")
		return
	}

	if l := len(password); l < 5 || l > 50 {
		fail(w, "Your password must be within 5 and 50 characters long")
		return
	}

	if _, err := db.GetUserByUsername(username); err == nil {
		fail(w, "A user with that username already exists")
		return
	}

	if !signupRegex.MatchString(username) {
		fail(w, "Usernames can only contain numbers, letters, and underscores")
		return
	}

	err := db.CreateUser(username, password)
	if err != nil {
		fail(w, "Couldn't create user")
		return
	}

	http.Redirect(w, r, "/api/login", 307)
}

// Login endpoint (/api/login POST)
func Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	user, err := db.CheckLogin(r.FormValue("username"), r.FormValue("password"))
	if err != nil {
		fail(w, "Incorrect username or password")
		return
	}

	sessionID := uuid.NewV4().String()
	db.Sessions[sessionID] = user.ID

	sessionCookie := &http.Cookie{
		Name:    "session",
		Value:   sessionID,
		Path:    "/",
		Expires: time.Now().Add(time.Hour * 24 * 30),
	}
	http.SetCookie(w, sessionCookie)

	http.Redirect(w, r, "/", 303)
}

// Logout endpoint (/logout GET)
func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")

	if cookie == nil || err != nil {
		http.Error(w, "Error signing out", 500)
		return
	}

	if _, ok := db.Sessions[cookie.Value]; ok {
		delete(db.Sessions, cookie.Value)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// CreateComment endpoint (/api/comments POST)
func CreateComment(w http.ResponseWriter, r *http.Request) {
	userID, err := auth(w, r)
	if err != nil {
		return
	}

	r.ParseForm()
	var comment db.Comment

	if err := decoder.Decode(&comment, r.Form); err != nil {
		fail(w, "Request was malformed")
		return
	}

	comment.UserID = userID

	if l := len(comment.Body); l < 3 || l > 1000 {
		fail(w, "Comment must be within 3 and 1000 characters")
		return
	}

	err = comment.Create()
	if err != nil {
		fmt.Println(err)
		fail(w, "Error saving comment")
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%s#%s", comment.PostID, comment.ID), 303)
}

// CreateReply endpoint (/api/replies POST)
func CreateReply(w http.ResponseWriter, r *http.Request) {
	userID, err := auth(w, r)
	if err != nil {
		return
	}

	r.ParseForm()
	var reply db.Reply

	if err := decoder.Decode(&reply, r.PostForm); err != nil {
		fail(w, "Request was malformed")
		return
	}

	if l := len(reply.Body); l < 3 || l > 1000 {
		fail(w, "Reply must be within 3 and 1000 characters")
		return
	}

	reply.UserID = userID

	err = reply.Create()
	if err != nil {
		fmt.Println(err)
		fail(w, "Error saving reply")
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s#%s", r.Referer(), reply.ID), 303)
}

// CreatePost endpoint (/api/post POST)
func CreatePost(w http.ResponseWriter, r *http.Request) {
	userID, err := auth(w, r)
	if err != nil {
		return
	}

	r.ParseForm()
	var post db.Post

	if err := decoder.Decode(&post, r.Form); err != nil {
		fmt.Println(err)
		fail(w, "Request was malformed")
		return
	}

	post.UserID = userID
	post.Title = strings.TrimSpace(post.Title)
	post.Votes = 0

	if l := len(post.Title); l < 3 || l > 140 {
		fail(w, "Title must be within 3 and 140 characters")
		return
	}

	if l := len(post.Body); l < 3 || l > 10000 {
		fail(w, "Body must be within 3 and 10,000 characters")
		return
	}

	err = post.Create()
	if err != nil {
		fmt.Println(err)
		fail(w, "Error creating post")
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%s", post.ID), 303)
}

// Vote endpoint (/api/vote/{post} POST)
func Vote(w http.ResponseWriter, r *http.Request) {
	userID, err := auth(w, r)
	if err != nil {
		return
	}

	postID := mux.Vars(r)["post"]

	vote := db.Vote{
		UserID: userID,
		PostID: postID,
	}

	if vote.Exists() {
		fail(w, "You've already voted on that")
		return
	}

	if _, err := vote.Create(); err != nil {
		fmt.Println(err)
		fail(w, "Error saving vote")
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%s", vote.PostID), 303)
}

// Post endpoint (/api/post/{post} GET)
func Post(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["post"]
	post, err := db.GetPost(postID)
	if err != nil {
		http.Error(w, "Post not found", 404)
		return
	}

	post.GetComments()
	post.GetCommentReplies()

	bytes, err := json.MarshalIndent(post, "", "\t")
	if err != nil {
		http.Error(w, "Error marshalling JSON", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

// Browse endpoint (/api/browse/(popular|recent) GET)
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
		http.Error(w, "Invalid page option, try \"popular\" or \"recent\"", 400)
		return
	}

	if err != nil {
		http.Error(w, "Error retrieving posts", 500)
		return
	}

	bytes, err := json.MarshalIndent(posts, "", "\t")
	if err != nil {
		http.Error(w, "Error marshalling JSON", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

// Conversation endpoint (/api/conversation GET)
func Conversation(w http.ResponseWriter, r *http.Request) {
	comments, err := db.GetRecentComments()
	if err != nil {
		http.Error(w, "Error retrieving comments", 500)
		return
	}

	for i := range comments {
		comments[i].GetReplies()
	}

	bytes, err := json.MarshalIndent(comments, "", "\t")
	if err != nil {
		http.Error(w, "Error marshalling JSON", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}
