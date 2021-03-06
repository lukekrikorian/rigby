package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"site/config"
	"site/db"
	"strings"
	"time"

	"github.com/gorilla/mux"

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
		fmt.Fprintf(w, "Success")
	}
}

type loginJSON struct {
	Username string
	Password string
}

// Login endpoint (/api/login POST)
func Login(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var login db.Login

	if err := decoder.Decode(&login); err != nil {
		http.Error(w, "Incorrect login", 500)
		return
	}

	if user, err := login.Check(); err == nil {

		sessionID := uuid.Must(uuid.NewV4()).String()
		db.Sessions[sessionID] = user.ID

		sessionCookie := &http.Cookie{
			Name:     "session",
			Value:    sessionID,
			Path:     "/",
			HttpOnly: true,
		}

		if login.SaveSession {
			sessionCookie.Expires = time.Now().Add(time.Hour * 24 * 7)
		}

		if config.Config.Server.Port == 443 {
			sessionCookie.Secure = true
		}

		http.SetCookie(w, sessionCookie)

		fmt.Fprintf(w, "Success")
		return
	}
	http.Error(w, "Incorrect login", 500)
}

// IsLoggedIn endpoint (/isLoggedIn GET)
func IsLoggedIn(w http.ResponseWriter, r *http.Request) {
	if _, err := db.CheckAuth(r); err != nil {
		http.Error(w, "Not logged in", 500)
		return
	}
	fmt.Fprintf(w, "Logged in")
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

	_, err = db.GetPost(comment.PostID)
	if err != nil {
		http.Error(w, "Invalid post", 404)
		return
	}

	err = comment.Create()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error saving comment", 500)
		return
	}

	fmt.Fprintf(w, "Success")
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

	_, err = db.GetComment(reply.ParentID)
	if err != nil {
		http.Error(w, "Invalid comment", 404)
	} else {
		err = reply.Create()
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error saving reply", 500)
			return
		}
		fmt.Fprintf(w, "Success")
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
	post.Votes = 0

	if l := len(post.Title); l < 3 || l > 140 {
		http.Error(w, "Title must be within 3 and 140 characters", 500)
		return
	}

	if l := len(post.Body); l < db.MinimumPostLength || l > 10000 {
		http.Error(w, fmt.Sprintf("Post body must be within %d and 10,000 characters", db.MinimumPostLength), 500)
		return
	}

	err = post.Create()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error creating post", 500)
		return
	}
	fmt.Fprintf(w, "/posts/%s", post.ID)
}

// Vote endpoint (/api/vote/{post} POST)
func Vote(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["post"]
	userID, err := db.CheckAuth(r)
	if err != nil {
		http.Error(w, "You aren't logged in", 500)
		return
	}
	vote := db.Vote{
		UserID: userID,
		PostID: postID,
	}

	if vote.Exists() {
		http.Error(w, "You've already voted on that", 500)
		return
	}

	votes, err := vote.Create()
	if err == nil {
		fmt.Fprintf(w, "%d", votes)
	} else {
		fmt.Println(err)
		http.Error(w, "Error saving vote", 500)
	}
}

// Post endpoint (/api/posts/{post} GET)
func Post(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["post"]
	post, err := db.GetPost(postID)
	if err != nil {
		http.Error(w, "Post not found", 404)
		return
	}

	post.GetComments()
	post.GetCommentReplies()

	sendAsJSON(post, w, r)
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

	sendAsJSON(posts, w, r)
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

	sendAsJSON(comments, w, r)
}

// User endpoint (/api/users/{user} GET)
func User(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["user"]

	var (
		user db.User
		err  error
	)

	if username == "me" {
		userID, err := db.CheckAuth(r)
		if err != nil {
			http.Error(w, "You aren't logged in", 500)
			return
		}
		user, err = db.GetUserByID(userID)
	} else {
		user, err = db.GetUserByUsername(username)
	}

	if err != nil {
		http.Error(w, "Error retrieving user", 500)
		return
	}

	user.GetPosts()

	sendAsJSON(user, w, r)
}

func sendAsJSON(JSON interface{}, w http.ResponseWriter, r *http.Request) {
	bytes, err := json.Marshal(JSON)
	if err != nil {
		http.Error(w, "Error formatting JSON", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}
