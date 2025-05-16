package errgroup

import(
	"log"
	"time"
	"errors"
	"database/sql"
	
	"github.com/RileySun/Scaled/utils"
)

var DB *sql.DB

type Post struct {
	date time.Time
	title string
	contents string
}

type Comment struct {
	date time.Time 
	contents string
}

func init() {
	creds := utils.LoadCredentials("../../utils/.env")
	var err error
	DB, err = utils.NewDB(creds.Host, creds.Port, creds.User, creds.Pass, creds.Database)
	if err != nil {
		log.Fatal("Could not connect to DB")
	}
}

func getUsername(userID int) (string, error) {
	row := DB.QueryRow("SELECT `name` FROM Users WHERE `id` = ?;", userID)
	
	var name string
	scanErr := row.Scan(&name)
	if scanErr != nil {
		return "", errors.New("DB Error - getUsername")
	}
	
	return name, nil
}

func getLastPost(userID int) (*Post, error) {
	row := DB.QueryRow("SELECT `date`, `title`, `contents` FROM Posts WHERE `userID` = ? ORDER BY date DESC LIMIT 1;", userID)
	
	var post Post
	scanErr := row.Scan(&post.date, &post.title, &post.contents)
	if scanErr != nil {
		return &post, errors.New("DB Error - getLastPost")
	}
	
	return &post, nil
}

func getLastComment(userID int) (*Comment, error) {
	row := DB.QueryRow("SELECT `date`, `contents` FROM Comments WHERE `userID` = ? ORDER BY date DESC LIMIT 1;", userID)
	
	var comment Comment
	scanErr := row.Scan(&comment.date, &comment.contents)
	if scanErr != nil {
		return &comment, errors.New("DB Error- getLastComment")
	}
	
	return &comment, nil
}