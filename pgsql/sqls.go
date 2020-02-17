package pgsql

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"middleware/handler/db"
	"time"

	"github.com/lib/pq"
)

// GetPostByID use pid to filter posts, then returns it
func (pg *PGSQL) GetPostByID(pid int) (*db.Post, error) {
	var (
		id    int
		c, t  string
		cDate time.Time
		mDate pq.NullTime

		tgs pq.StringArray
		p   *db.Post
	)
	err := pg.instance.QueryRow(`SELECT * FROM public.getPostByID($1)`, pid).Scan(&id, &t, &cDate, &mDate, &c, &tgs)
	if err != nil {
		return nil, fmt.Errorf("select from getPostByID(): %v", err)
	}
	if mDate.Valid {
		cD := db.Jstime(cDate)
		mD := db.Jstime(mDate.Time)

		p = &db.Post{PostID: id, Title: t, CDate: &cD, MDate: &mD, Content: c, Tags: []string(tgs)}
	} else {
		cD := db.Jstime(cDate)
		p = &db.Post{PostID: pid, Title: t, CDate: &cD, MDate: nil, Content: c, Tags: []string(tgs)}
	}
	return p, nil
}

// GetCommentsCount return total number of comments bound to a certain post
func (pg *PGSQL) GetCommentsCount(pid int) (int, error) {
	var (
		count int
	)
	err := pg.instance.QueryRow(`SELECT public.getCommentsCount($1)`, pid).Scan(&count)
	if err != nil {
		return -1, fmt.Errorf("select from getCommentsCount(): %v", err)
	}
	return count, nil
}

// GetCommentsByPage accepts pageSize and page, then returns comments bound to a post
func (pg *PGSQL) GetCommentsByPage(pid int, pageSize int, page int) ([]db.Comment, error) {
	cmts := []db.Comment{}

	rs, err := pg.instance.Query(`SELECT * FROM public.getCommentsByPage($1, $2, $3)`, pid, pageSize, page)
	if err != nil {
		return nil, fmt.Errorf("select from getPostsByPage(): %v", err)
	}
	defer rs.Close()
	for rs.Next() {
		var (
			id, cid int
			e       string
			cDate   time.Time
			c       string
		)
		err := rs.Scan(&id, &cid, &e, &cDate, &c)
		if err != nil {
			return nil, fmt.Errorf("parse query result: %v", err)
		}
		cD := db.Jstime(cDate)
		cmts = append(cmts, db.Comment{PostID: id, CommentID: cid, Email: e, CDate: &cD, Content: c})
	}
	if rs.Err() != nil {
		return nil, fmt.Errorf("perform query: %v", err)
	}
	if len(cmts) == 0 {
		return nil, errors.New("no comments found")
	}
	return cmts, nil

}

// GetPostsCount returns count of total posts in db
func (pg *PGSQL) GetPostsCount() (int, error) {
	var (
		count int
	)
	err := pg.instance.QueryRow(`SELECT public.getPostsCount()`).Scan(&count)
	if err != nil {
		return -1, fmt.Errorf("select from getPostsCount(): %v", err)
	}
	return count, nil
}

// GetPosts use page and pageSize to select posts, then return them
func (pg *PGSQL) GetPosts(pageSize, page int) ([]db.Post, error) {
	posts := []db.Post{}
	rs, err := pg.instance.Query(`SELECT * FROM public.getPostsByPage($1, $2)`, pageSize, page)
	if err != nil {
		return nil, fmt.Errorf("select from getPostsByPage(): %v", err)
	}
	defer rs.Close()

	for rs.Next() {
		var (
			id    int
			c, t  string
			cDate time.Time
			mDate pq.NullTime
			tgs   pq.StringArray
		)
		err := rs.Scan(&id, &t, &cDate, &mDate, &c, &tgs)
		if err != nil {
			return nil, fmt.Errorf("parse query result: %v", err)
		}
		if mDate.Valid {
			cD := db.Jstime(cDate)
			mD := db.Jstime(mDate.Time)
			posts = append(posts, db.Post{PostID: id, Title: t, CDate: &cD, MDate: &mD, Content: c, Tags: []string(tgs)})
		} else {
			cD := db.Jstime(cDate)
			posts = append(posts, db.Post{PostID: id, Title: t, CDate: &cD, MDate: nil, Content: c, Tags: []string(tgs)})
		}

	}
	if rs.Err() != nil {
		return nil, fmt.Errorf("perform query: %v", err)
	}
	if len(posts) == 0 {
		return nil, errors.New("no posts found")
	}
	return posts, nil
}

// GetPostsCountByFTS uses search string to filter posts, then returns count of found posts
func (pg *PGSQL) GetPostsCountByFTS(search string) (int, error) {
	var (
		count int
	)

	err := pg.instance.QueryRow(`SELECT public.getPostsCountByFTS($1)`, search).Scan(&count)
	if err != nil {
		return -1, fmt.Errorf("select from getPostsCountByFTS(): %v", err)
	}
	return count, nil
}

// GetPostsByFTS uses search string, page, pageSize to filter posts, then returns them
func (pg *PGSQL) GetPostsByFTS(search string, pageSize, page int) (posts []db.Post, err error) {
	posts = []db.Post{}

	var rs *sql.Rows
	rs, err = pg.instance.Query(`SELECT * FROM public.getPostsByFTS($1, $2, $3)`, search, pageSize, page)
	if err != nil {
		return nil, fmt.Errorf("select from getPostsByFTS(): %v", err)
	}
	defer rs.Close()

	for rs.Next() {
		var (
			pid   int
			t, c  string
			cDate time.Time
			mDate pq.NullTime
			tgs   pq.StringArray
		)
		err := rs.Scan(&pid, &t, &cDate, &mDate, &c, &tgs)
		if err != nil {
			return nil, fmt.Errorf("parse query result: %v", err)
		}
		if mDate.Valid {
			cD := db.Jstime(cDate)
			mD := db.Jstime(mDate.Time)
			posts = append(posts, db.Post{PostID: pid, Title: t, CDate: &cD, MDate: &mD, Content: c, Tags: tgs})
		} else {
			cD := db.Jstime(cDate)
			posts = append(posts, db.Post{PostID: pid, Title: t, CDate: &cD, MDate: nil, Content: c, Tags: tgs})

		}

	}
	if rs.Err() != nil {
		return nil, fmt.Errorf("perform query correctly: %v", err)
	}
	if len(posts) == 0 {
		return nil, errors.New("no posts found")
	}
	return posts, nil
}

// UserLogin uses username, password to login user, then returns true if user existing
func (pg *PGSQL) UserLogin(userName string, pass [sha256.Size]byte) (*db.User, error) {
	var (
		id  int
		unm string
		pri int
	)

	err := pg.instance.QueryRow(`SELECT * FROM public.userLogin($1, $2)`, userName, pass[:]).Scan(&id, &unm, &pri)
	if err != nil {
		return nil, fmt.Errorf("select from userLogin(): %v", err)
	}
	return &db.User{UID: id, UserName: unm, Privilege: pri}, nil
}

// InsertPost inserts post and return its pid
func (pg *PGSQL) InsertPost(t string, c string, tags []string) (int, error) {
	var (
		pid int
	)
	err := pg.instance.QueryRow(`SELECT * FROM public.insertPost($1, $2, $3)`, t, c, pq.StringArray(tags)).Scan(&pid)
	if err != nil {
		return -1, fmt.Errorf("select from insertPost(): %v", err)
	}
	return pid, nil
}

// DeletePost delete post, returns true if deletion performed while false if not found
func (pg *PGSQL) DeletePost(pid int) (bool, error) {
	var (
		performed bool
	)
	err := pg.instance.QueryRow(`SELECT * FROM public.deletePost($1)`, pid).Scan(&performed)
	if err != nil {
		return false, fmt.Errorf("select from deletePost(): %v", err)
	}
	return performed, nil
}

// UpdatePost update existing post, returns true if update performed while false if not found
func (pg *PGSQL) UpdatePost(pid int, nTitle, nContent string, nTags []string) (bool, error) {
	var (
		performed bool
	)

	err := pg.instance.QueryRow(`SELECT * FROM public.updatePost($1, $2, $3, $4)`, pid, nTitle, nContent, pq.StringArray(nTags)).Scan(&performed)

	if err != nil {
		return false, fmt.Errorf("select from updatePost(): %v", err)

	}
	return performed, nil
}

// InsertComment insert comment, returns cid of the inserted comment
func (pg *PGSQL) InsertComment(pid int, c, authorEmail string) (int, error) {
	var (
		cid int
	)
	err := pg.instance.QueryRow(`SELECT * FROM public.insertComment($1, $2, $3)`, pid, c, authorEmail).Scan(&cid)
	if err != nil {
		return -1, fmt.Errorf("select from insertComment(): %v", err)
	}
	return cid, nil
}

// DeleteComment deletes comment of cid, returns true if performed while false in case of not found
func (pg *PGSQL) DeleteComment(cid int) (bool, error) {
	var (
		performed bool
	)
	err := pg.instance.QueryRow(`SELECT * FROM public.deleteComment($1)`, cid).Scan(&performed)
	if err != nil {
		return false, fmt.Errorf("select from deleteComment(): %v", err)
	}
	return performed, nil
}

// UpdateComment updates comment of cid, returns true if performed while false in case of not found
func (pg *PGSQL) UpdateComment(cid int, nConent, nAE string) (bool, error) {
	var (
		performed bool
	)
	err := pg.instance.QueryRow(`SELECT * FROM public.updateComment($1, $2, $3)`, cid, nConent, nAE).Scan(&performed)
	if err != nil {
		return false, fmt.Errorf("select from updateComment(): %v", err)
	}
	return performed, nil
}

// GetUser accepts userName and passWord to check if certain user exists for purpose of login
func (pg *PGSQL) GetUser(userName string, passWord [sha256.Size]byte) (*db.User, error) {
	var (
		id  int
		uN  string
		pri int
	)
	err := pg.instance.QueryRow(`SELECT * FROM public.getUser($1, $2)`, userName, passWord[:]).Scan(&id, &uN, &pri)
	if err != nil {
		return nil, fmt.Errorf("select from getUser(): %v", err)
	}
	return &db.User{UID: id, UserName: uN, Privilege: pri}, nil
}

// InsertUser inserts new user record into database and returns its uid
func (pg *PGSQL) InsertUser(userName string, passWord [sha256.Size]byte) (int, error) {
	var (
		uid int
	)
	err := pg.instance.QueryRow(`SELECT * FROM public.insertUser($1, $2, $3)`, userName, passWord[:]).Scan(&uid)
	if err != nil {
		return -1, fmt.Errorf("select from insertUser(): %v", err)
	}
	return uid, nil
}

// UpdateUser update passWord of existing user based on uid
func (pg *PGSQL) UpdateUser(uid int, nPW [sha256.Size]byte) (bool, error) {
	var (
		performed bool
	)
	err := pg.instance.QueryRow(`SELECT * FROM public.updateUser($1,$2,$3)`, uid, nPW).Scan(&performed)
	if err != nil {
		return false, fmt.Errorf("select from updateUser(): %v", err)
	}
	return performed, nil
}
