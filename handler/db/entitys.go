package db

import (
	"fmt"
	"strconv"
	"time"
)

// Post model in blog
type Post struct {
	PostID  int      `json:"pid"`
	Title   string   `json:"title"`
	CDate   *Jstime  `json:"cDate"`
	MDate   *Jstime  `json:"mDate"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

// Comment contains info about a comment of a post in blog
type Comment struct {
	PostID    int     `json:"pid"`
	CommentID int     `json: "cid`
	Email     string  `json: "email"`
	CDate     *Jstime `json:"cDate"`
	Content   string  `json: "content"`
}

// User contains info that depicts a user
type User struct {
	UID       int    `json:"uid"`
	UserName  string `json:"userName"`
	Privilege int    `json:"privilege"`
}

// PostsPage packs posts and maxpage together for convenience
type PostsPage struct {
	Posts   []Post `json:"posts"`
	MaxPage int    `json:"maxPage"`
}

// CommentsPage that packs comments and maxPage number of these comments
type CommentsPage struct {
	Comments []Comment `json: "comments"`
	MaxPage  int       `json:"maxPage"`
}

// Response is our generic json response schema
type Response struct {
	Err  *JsError    `json:"err"`
	Data interface{} `json:"data"`
}

// JsError helps to Marhshal error data into json
type JsError struct {
	Err error
}

// MarshalJSON marshals string of error into json field
func (jE *JsError) MarshalJSON() ([]byte, error) {
	if jE.Err == nil {
		return []byte("null"), nil
	}
	return []byte(strconv.Quote(jE.Err.Error())), nil
}

// BlogContext is our context type to be bound to http request
type BlogContext string

// Jstime for better time.Time marshal in json message
type Jstime time.Time

// MarshalJSON marshals jstime to Mon Jan _2, 2006 format
func (jT *Jstime) MarshalJSON() ([]byte, error) {
	ts := time.Time(*jT).Format(time.RFC3339)
	return []byte(strconv.Quote(ts)), nil
}

// UnmarshalJSON unmarshals Mon Jan _2, 2006 format to jstime
func (jT *Jstime) UnmarshalJSON(b []byte) error {
	s, err := strconv.Unquote(string(b))
	if err != nil {
		// return fmt.Errorf("Cannot unquote time field in json, it may not be a valid string: %v", err)
		return fmt.Errorf("unquote time field in json: %v", err)
	}
	ts, err := time.Parse("Mon Jan _2, 2006", s)
	if err != nil {
		return fmt.Errorf("parse time: %v", err)
	}
	*jT = Jstime(ts)
	return nil
}
