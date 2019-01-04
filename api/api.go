package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"site/db"
	"strings"

	uuid "github.com/satori/go.uuid"
)

var signupRegex = regexp.MustCompile("^[a-zA-Z0-9_]+$")

type signupJSON struct {
	Username string
	Password string
}

// Signup endpoint (/api/signup POST)
func Signup(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var signup signupJSON

	if err := decoder.Decode(&signup); err != nil {
		http.Error(w, "Invalid signup", 500)
		return
	}

	if l := len(signup.Username); l < 3 || l > 16 {
		http.Error(w, "Your username must be within 3 and 16 characters long", 500)
		return
	}

	if l := len(signup.Password); l < 5 || l > 50 {
		http.Error(w, "Your password must be within 5 and 50 characters long", 500)
		return
	}

	if _, err := db.GetUserByUsername(signup.Username); err == nil {
		http.Error(w, "A user with that username already exists", 500)
		return
	}

	if !signupRegex.MatchString(signup.Username) {
		http.Error(w, "Usernames can only contain numbers, letters, and underscores", 500)
		return
	}

	err := db.CreateUser(signup.Username, signup.Password)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Could not create user", 500)
	} else {
		fmt.Println("New user", signup.Username)
		http.Redirect(w, r, "/api/login", http.StatusTemporaryRedirect)
	}
}

type loginJSON struct {
	Username string
	Password string
}

// Login endpoint (/api/login POST)
func Login(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var login loginJSON

	if err := decoder.Decode(&login); err != nil {
		http.Error(w, "Incorrect login", 500)
		return
	}

	if user, err := db.CheckLogin(login.Username, login.Password); err == nil {

		sessionID := uuid.Must(uuid.NewV4()).String()
		db.Sessions[sessionID] = user.ID

		sessionCookie := &http.Cookie{
			Name:  "session",
			Value: sessionID,
			Path:  "/",
		}
		http.SetCookie(w, sessionCookie)

		w.WriteHeader(http.StatusSeeOther)
		fmt.Fprintf(w, "/")
		return
	}
	http.Error(w, "Incorrect login", 500)
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
	userID, err := db.CheckAuth(r)
	if err != nil {
		http.Error(w, "You aren't logged in", 500)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var comment db.Comment

	if err := decoder.Decode(&comment); err != nil {
		http.Error(w, "Request was malformed", 500)
		return
	}

	comment.UserID = userID

	if l := len(comment.Body); l < 3 || l > 1000 {
		http.Error(w, "Comment must be within 3 and 1000 characters", 500)
		return
	}

	err = comment.Create()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error saving comment", 500)
		return
	}

	w.WriteHeader(http.StatusSeeOther)
	fmt.Fprintf(w, "/posts/%s", comment.PostID)
}

// CreateReply endpoint (/api/replies POST)
func CreateReply(w http.ResponseWriter, r *http.Request) {
	userID, err := db.CheckAuth(r)
	if err != nil {
		http.Error(w, "You aren't logged in", 500)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var reply db.Reply

	if err := decoder.Decode(&reply); err != nil {
		http.Error(w, "Request was malformed", 500)
		return
	}

	if l := len(reply.Body); l < 3 || l > 1000 {
		http.Error(w, "Reply must be within 3 and 1000 characters", 500)
		return
	}

	reply.UserID = userID

	err = reply.Create()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error saving reply", 500)
		return
	}

	comment, err := db.GetComment(reply.ParentID)
	if err == nil {
		w.WriteHeader(http.StatusSeeOther)
		fmt.Fprintf(w, "/posts/%s#%s", comment.PostID, reply.ID)
	} else {
		w.WriteHeader(http.StatusSeeOther)
		fmt.Fprintf(w, "/")
	}
}

// CreatePost endpoint (/api/post POST)
func CreatePost(w http.ResponseWriter, r *http.Request) {
	userID, err := db.CheckAuth(r)
	if err != nil {
		http.Error(w, "You aren't logged in", 500)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var post db.Post

	if err := decoder.Decode(&post); err != nil {
		http.Error(w, "Request was malformed", 500)
		return
	}

	post.UserID = userID
	post.Title = strings.TrimSpace(post.Title)

	if l := len(post.Title); l < 3 || l > 140 {
		http.Error(w, "Title must be within 3 and 140 characters", 500)
		return
	}

	if l := len(post.Body); l < 3 || l > 5000 {
		http.Error(w, "Post body must be within 3 and 5000 characters", 500)
		return
	}

	err = post.Create()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error creating post", 500)
		return
	}
	w.WriteHeader(http.StatusSeeOther)
	fmt.Fprintf(w, "/posts/%s", post.ID)
}
