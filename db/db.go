package db

import (
	"errors"
	"log"
	"net/http"

	"github.com/satori/go.uuid"

	// Needed for sqlx
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"golang.org/x/crypto/bcrypt"
)

// DB is the database connection itself
var DB *sqlx.DB

// Sessions is user session tokens: token -> userID
var Sessions = make(map[string]string)

// User representation
type User struct {
	ID       string
	Username string
	Password string
	Created  string
	Posts    []Post
}

// Post representation
type Post struct {
	ID        string `db:"id"`
	GamerRage bool   `db:"gamerRage"`
	UserID    string `db:"userID"`
	Author    string
	Title     string
	Body      string
	Created   string
	Votes     uint
	Comments  []Comment
}

// Comment representation
type Comment struct {
	ID      string `db:"id"`
	UserID  string `db:"userID"`
	Author  string
	Body    string
	PostID  string `db:"postID" json:"postid"`
	Created string
	Replies []Reply
	Parent  *Post
}

// Reply representation
type Reply struct {
	ParentID string `db:"parentID" json:"parentid"`
	ID       string `db:"id"`
	UserID   string `db:"userID"`
	Author   string
	Body     string
	Created  string
}

// Vote representation
type Vote struct {
	UserID string `db:"userID"`
	PostID string `db:"postID"`
}

// Init initializes and checks the DB connection for errors
func Init(url string) {
	DB, _ = sqlx.Open("mysql", url)
	if DB.Ping() != nil {
		log.Fatal("Error connecting to database")
	}
}

// CheckAuth validates a request session
func CheckAuth(r *http.Request) (UserID string, Error error) {
	sessionCookie, Error := r.Cookie("session")
	if sessionCookie == nil || Error != nil {
		return "", errors.New("Error retrieving session cookie")
	}
	if UserID, ok := Sessions[sessionCookie.Value]; ok {
		return UserID, nil
	}
	return "", errors.New("User cookie incorrect")
}

// CheckLogin validates a login, comparing the two hashes
func CheckLogin(username, password string) (User User, Error error) {
	var hash string
	if Error = DB.Get(&hash, "SELECT password FROM users WHERE username = ?", username); Error == nil {
		if isCorrect := CheckPasswordHash(password, hash); !isCorrect {
			Error = errors.New("Incorrect login")
			return
		}
		Error = DB.Get(&User, "SELECT * FROM users WHERE username = ?", username)
	}
	return
}

// CreateUser saves a new user profile
func CreateUser(username, password string) (Error error) {
	id := uuid.Must(uuid.NewV4())
	if password, Error = HashPassword(password); Error == nil {
		_, Error = DB.Exec("INSERT INTO users VALUES (?, ?, ?, NOW())", id, username, password)
	}
	return
}

// GetUserByID returns a user based on their ID
func GetUserByID(id string) (User User, Error error) {
	Error = DB.Get(&User, "SELECT * FROM users WHERE id = ?", id)
	return
}

// GetUserByUsername returns a user based on their username
func GetUserByUsername(username string) (User User, Error error) {
	Error = DB.Get(&User, "SELECT * FROM users WHERE username = ?", username)
	return
}

// Create saves a post
func (p *Post) Create() (Error error) {
	p.ID = uuid.Must(uuid.NewV4()).String()
	user, Error := GetUserByID(p.UserID)
	if Error == nil {
		p.Author = user.Username
		q := "INSERT INTO posts VALUES (:id, :userID, :author, :title, :body, :gamerRage, :votes, NOW())"
		_, Error = DB.NamedExec(q, p)
	}
	return
}

// Create saves a comment
func (c *Comment) Create() (Error error) {
	c.ID = uuid.Must(uuid.NewV4()).String()

	user, _ := GetUserByID(c.UserID)
	c.Author = user.Username

	q := "INSERT INTO comments VALUES (:id, :postID, :userID, :author, :body, NOW())"
	_, Error = DB.NamedExec(q, c)
	return
}

// Create saves a reply
func (r *Reply) Create() (Error error) {
	r.ID = uuid.Must(uuid.NewV4()).String()

	user, _ := GetUserByID(r.UserID)
	r.Author = user.Username

	q := "INSERT INTO replies VALUES (:id, :parentID, :userID, :author, :body, NOW())"
	_, Error = DB.NamedExec(q, r)
	return
}

// Create saves a vote
func (v *Vote) Create() (CurrentVotes int32, Error error) {
	Error = DB.Get(&CurrentVotes, "SELECT votes FROM posts WHERE id = ?", v.PostID)
	if Error == nil {
		insertQuery := "INSERT INTO votes VALUES (:userID, :postID)"
		if _, Error = DB.NamedExec(insertQuery, v); Error == nil {
			CurrentVotes++
			updateQuery := "UPDATE posts SET votes = ? WHERE id = ?"
			if _, Error = DB.Exec(updateQuery, CurrentVotes, v.PostID); Error == nil {
				return CurrentVotes, Error
			}
		}
	}
	return 0, Error
}

// Exists checks if a post has already been saved
func (v *Vote) Exists() (Exists bool) {
	var userID string
	q := "SELECT userID FROM votes WHERE userID = ? AND postID = ? LIMIT 1"

	if err := DB.Get(&userID, q, v.UserID, v.PostID); err == nil {
		return true
	}
	return false
}

// GetPost returns a post based on the post's ID
func GetPost(ID string) (Post Post, Error error) {
	Error = DB.Get(&Post, "SELECT * FROM posts WHERE id = ?", ID)
	return
}

// GetComment returns a top-level comment based on the comment's ID
func GetComment(ID string) (Comment Comment, Error error) {
	Error = DB.Get(&Comment, "SELECT * FROM comments WHERE id = ?", ID)
	return
}

// GetRecentPosts returns a list of recent posts
func GetRecentPosts() (Posts []Post, Error error) {
	Error = DB.Select(&Posts, "SELECT * FROM posts ORDER BY created DESC LIMIT 50")
	return
}

// GetPopularPosts returns a list of popular posts
func GetPopularPosts() (Posts []Post, Error error) {
	Error = DB.Select(&Posts, "SELECT * FROM posts ORDER BY votes DESC LIMIT 50")
	return
}

// GetRecentComments returns a list of recent comments
func GetRecentComments() (Comments []Comment, Error error) {
	Error = DB.Select(&Comments, "SELECT * FROM comments ORDER BY created DESC LIMIT 50")
	return
}

// GetPosts populates a list of the users posts
func (u *User) GetPosts() (Error error) {
	Error = DB.Select(&u.Posts, "SELECT * FROM posts WHERE userID = ? ORDER BY created DESC", u.ID)
	return
}

// GetComments populates a list of comments on a post
func (p *Post) GetComments() (Error error) {
	Error = DB.Select(&p.Comments, "SELECT * FROM comments WHERE postID = ? ORDER BY created ASC", p.ID)
	return
}

// GetCommentReplies populates the replies for each comment on a post
func (p *Post) GetCommentReplies() {
	for i, comment := range p.Comments {
		DB.Select(&p.Comments[i].Replies, "SELECT * FROM replies WHERE parentID = ? ORDER BY created ASC", comment.ID)
	}
}

// GetReplies populates a list of replies to a comment
func (c *Comment) GetReplies() (Error error) {
	Error = DB.Select(&c.Replies, "SELECT * FROM replies WHERE parentID = ? ORDER BY created ASC", c.ID)
	return
}

// GetParent populates a comment's parent post data
func (c *Comment) GetParent() (Error error) {
	Error = DB.Get(&c.Parent, "SELECT * FROM posts WHERE id = ?", c.PostID)
	return
}

// HashPassword hashes a password using bycrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes), err
}

// CheckPasswordHash checks a password and hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
