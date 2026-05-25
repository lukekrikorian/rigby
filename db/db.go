package db

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"site/config"

	uuid "github.com/satori/go.uuid"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"golang.org/x/crypto/bcrypt"
)

// DB is the database connection itself
var DB *sqlx.DB

// Sessions is user session tokens: token -> userID
var Sessions = make(map[string]string)

// Migration representation
type Migration struct {
	ID        string `db:"id"`
	Name      string
	Performed time.Time
}

// User representation
type User struct {
	ID       string `db:"id"`
	Username string
	Password string
	Created  time.Time
	Posts    []Post
}

// Post representation
type Post struct {
	ID        string `db:"id"`
	UserID    string `db:"userid"`
	Author    string
	Title     string
	Body      string
	GamerRage int `db:"gamerrage" json:"gamerrage"`
	VoteCount int `db:"votes"`
	Created   time.Time
	Comments  []Comment
	Votes     []Vote
}

// Comment representation
type Comment struct {
	ID      string `db:"id"`
	UserID  string `db:"userid"`
	Author  string
	Body    string
	PostID  string `db:"postid" json:"postid"`
	Created time.Time
	Replies []Reply
	Parent  *Post
}

// Reply representation
type Reply struct {
	ParentID string `db:"parentid"`
	ID       string `db:"id"`
	UserID   string `db:"userid"`
	Author   string
	Body     string
	Created  time.Time
}

// Vote representation
type Vote struct {
	UserID   string `db:"userid"`
	PostID   string `db:"postid"`
	Username string
}

// Init initializes and checks the DB connection for errors
func init() {
	DB, _ = sqlx.Connect("sqlite3", config.Config.DatabaseURL)

	if DB.Ping() != nil {
		log.Fatal("Error connecting to database")
	}

	Migrate()
	LoadSessions()
}

func Migrate() {
	migrationFiles, err := os.ReadDir("db/migrations")
	if err != nil {
		log.Fatal("Failed to load migrations from db/migrations")
	}

	for _, file := range migrationFiles {
		fileName := file.Name()
		nameParts := strings.Split(fileName, ".")
		id, name := nameParts[0], nameParts[1]

		var migration Migration
		err = DB.Get(&migration, "SELECT * FROM migrations WHERE id = $1 AND performed IS NOT NULL", id)
		if err == nil {
			continue
		}

		fileContents, err := os.ReadFile("db/migrations/" + fileName)
		if err != nil {
			log.Fatal("Failed to read migration", name)
		}

		log.Println("Applying migration", name)

		tx := DB.MustBegin()
		tx.MustExec(string(fileContents))
		tx.MustExec("INSERT INTO migrations (id, name) VALUES ($1, $2)", id, name)
		err = tx.Commit()
		if err != nil {
			log.Fatal("Error committing transaction for migration", name)
		}
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
	if Error = DB.Get(&hash, "SELECT password FROM users WHERE username = $1", username); Error == nil {
		if isCorrect := CheckPasswordHash(password, hash); !isCorrect {
			Error = errors.New("Incorrect login")
			return
		}
		Error = DB.Get(&User, "SELECT * FROM users WHERE username = $1", username)
	}
	return
}

// CreateUser saves a new user profile
func CreateUser(username, password string) (Error error) {
	id := uuid.NewV4()
	if password, Error = HashPassword(password); Error == nil {
		_, Error = DB.Exec("INSERT INTO users (id, username, password) VALUES ($1, $2, $3)", id, username, password)
	}
	return
}

// GetUserByID returns a user based on their ID
func GetUserByID(id string) (User User, Error error) {
	Error = DB.Get(&User, "SELECT * FROM users WHERE id = $1", id)
	return
}

// GetUserByUsername returns a user based on their username
func GetUserByUsername(username string) (User User, Error error) {
	Error = DB.Get(&User, "SELECT * FROM users WHERE username = $1", username)
	return
}

// Create saves a post
func (p *Post) Create() (Error error) {
	p.ID = uuid.NewV4().String()
	user, Error := GetUserByID(p.UserID)
	if Error == nil {
		p.Author = user.Username
		q := "INSERT INTO posts (id, userid, author, title, body, gamerrage, votes) VALUES (:id, :userid, :author, :title, :body, :gamerrage, :votes)"
		_, Error = DB.NamedExec(q, p)
	}
	return
}

// Create saves a comment
func (c *Comment) Create() (Error error) {
	c.ID = uuid.NewV4().String()

	user, _ := GetUserByID(c.UserID)
	c.Author = user.Username

	q := "INSERT INTO comments (id, postid, userid, author, body) VALUES (:id, :postid, :userid, :author, :body)"
	_, Error = DB.NamedExec(q, c)
	return
}

// Create saves a reply
func (r *Reply) Create() (Error error) {
	r.ID = uuid.NewV4().String()

	user, _ := GetUserByID(r.UserID)
	r.Author = user.Username

	q := "INSERT INTO replies (id, parentid, userid, author, body) VALUES (:id, :parentid, :userid, :author, :body)"
	_, Error = DB.NamedExec(q, r)
	return
}

// Create saves a vote
func (v *Vote) Create() (CurrentVotes int32, Error error) {
	Error = DB.Get(&CurrentVotes, "SELECT votes FROM posts WHERE id = $1", v.PostID)
	if Error == nil {
		insertQuery := "INSERT INTO votes (userid, postid) VALUES (:userid, :postid)"
		if _, Error = DB.NamedExec(insertQuery, v); Error == nil {
			CurrentVotes++
			updateQuery := "UPDATE posts SET votes = $1 WHERE id = $2"
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
	q := "SELECT userid FROM votes WHERE userid = $1 AND postid = $2 LIMIT 1"

	if err := DB.Get(&userID, q, v.UserID, v.PostID); err == nil {
		return true
	}
	return false
}

// LoadSessions loads the user sessions
func LoadSessions() {
	rows, _ := DB.Query("SELECT * FROM sessions")

	for rows.Next() {
		var token, id string
		rows.Scan(&token, &id)
		Sessions[token] = id
	}

	log.Println("Loaded user sessions")
}

// SaveSessions saves the server sessions
func SaveSessions() {
	DB.Exec("DELETE FROM sessions")
	for token, id := range Sessions {
		DB.Exec("INSERT INTO sessions (token, userid) VALUES ($1, $2)", token, id)
	}

	log.Println("Saved user sessions")
}

// GetPost returns a post based on the post's ID
func GetPost(ID string) (Post Post, Error error) {
	Error = DB.Get(&Post, "SELECT * FROM posts WHERE id = $1", ID)
	return
}

// GetComment returns a top-level comment based on the comment's ID
func GetComment(ID string) (Comment Comment, Error error) {
	Error = DB.Get(&Comment, "SELECT * FROM comments WHERE id = $1", ID)
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
	Error = DB.Select(&u.Posts, "SELECT * FROM posts WHERE userid = $1 ORDER BY created DESC", u.ID)
	return
}

// GetComments populates a list of comments on a post
func (p *Post) GetComments() (Error error) {
	Error = DB.Select(&p.Comments, "SELECT * FROM comments WHERE postid = $1 ORDER BY created ASC", p.ID)
	return
}

// GetVotes populates a list of votes on a post
func (p *Post) GetVotes() (Error error) {
	Error = DB.Select(&p.Votes, "SELECT votes.*, users.username FROM votes JOIN users ON userid = id WHERE postid = $1", p.ID)
	return
}

func (p Post) GetVoterList() string {
	if err := p.GetVotes(); err != nil {
		return ""
	}

	usernames := []string{}

	for _, vote := range p.Votes {
		usernames = append(usernames, vote.Username)
	}

	return strings.Join(usernames, ", ")
}

// GetCommentReplies populates the replies for each comment on a post
func (p *Post) GetCommentReplies() {
	for i, comment := range p.Comments {
		DB.Select(&p.Comments[i].Replies, "SELECT * FROM replies WHERE parentid = $1 ORDER BY created ASC", comment.ID)
	}
}

// GetReplies populates a list of replies to a comment
func (c *Comment) GetReplies() (Error error) {
	var replies []Reply
	Error = DB.Select(&replies, "SELECT * FROM replies WHERE parentid = $1 ORDER BY created ASC", c.ID)
	c.Replies = replies
	return
}

// GetParent populates a comment's parent post data
func (c *Comment) GetParent() (Error error) {
	var parent Post
	Error = DB.Get(&parent, "SELECT * FROM posts WHERE id = $1", c.PostID)
	c.Parent = &parent
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
