package handler

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"math"
	"middleware/handler/db"
	"net/http"
	"strconv"

	"github.com/Jeffail/gabs/v2"
)

func viewPost(d db.DB, r *http.Request) (*db.Post, error) {
	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, fmt.Errorf("convert id to int: %v", err)
	}
	if id < 0 {
		return nil, errors.New("id value cannot be less than 0")
	}
	post, err := d.GetPostByID(id)
	if err != nil {
		return nil, fmt.Errorf("get post by id: %v", err)
	}

	return post, nil

}

func viewPosts(d db.DB, r *http.Request) (*db.PostsPage, error) {
	filterStr := r.FormValue("keyword")

	pageStr := r.FormValue("page")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return nil, fmt.Errorf("convert page to int: %v", err)
	}
	if page <= 0 {
		return nil, errors.New("page value cannot be less than 1")
	}

	pageSizeStr := r.FormValue("pageSize")

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		return nil, fmt.Errorf("convert pageSize to int: %v", err)
	}
	if pageSize <= 0 {
		return nil, errors.New("pageSize value cannot be less than 1")
	}

	search := filterStr

	var (
		posts []db.Post
		count int
	)
	if search == "" {
		count, err = d.GetPostsCount()
	} else {
		count, err = d.GetPostsCountByFTS(search)
	}
	if err != nil {
		return nil, fmt.Errorf("get count of posts: %v", err)
	}
	maxPage := int(math.Ceil(float64(count) / float64(pageSize)))
	if page > maxPage && maxPage != 0 {
		return nil, errors.New("page number bigger than maxPage")
	}

	if search == "" {
		posts, err = d.GetPosts(pageSize, page)
	} else {
		posts, err = d.GetPostsByFTS(search, pageSize, page)
	}
	if err != nil {
		return nil, fmt.Errorf("get posts: %v", err)
	}

	return &db.PostsPage{Posts: posts, MaxPage: maxPage}, nil
}

func viewComments(d db.DB, r *http.Request) (*db.CommentsPage, error) {
	pidStr := r.FormValue("pid")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return nil, fmt.Errorf("convert pid to int: %v", err)
	}
	if pid <= 0 {
		return nil, errors.New("pid cannot be less than 1")
	}

	pageStr := r.FormValue("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return nil, fmt.Errorf("convert page to int: %v", err)
	}
	if page <= 0 {
		return nil, errors.New("page cannot be less than 1")
	}

	pageSizeStr := r.FormValue("pageSize")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		return nil, fmt.Errorf("convert pageSize to int: %v", err)
	}
	if pageSize <= 0 {
		return nil, errors.New("pageSize cannot be less than 1")
	}

	cnt, err := d.GetCommentsCount(pid)
	if err != nil {
		return nil, fmt.Errorf("get counts of comments: %v", err)
	}

	maxPage := int(math.Ceil(float64(cnt) / float64(pageSize)))
	if page > maxPage && maxPage != 0 {
		return nil, errors.New("page number bigger than max page")
	}

	cmts, err := d.GetCommentsByPage(pid, pageSize, page)
	if err != nil {
		return nil, fmt.Errorf("get comments: %v", err)
	}
	return &db.CommentsPage{Comments: cmts, MaxPage: maxPage}, nil

}

func changeComment(d db.DB, action string, r *http.Request) (int, error) {
	switch action {
	case "insert":
		var (
			pid            int
			content, email string
		)
		err := parseJSONReq(r, func(pJSON *gabs.Container) error {
			var ok bool
			pid, ok = pJSON.Path("pid").Data().(int)
			if !ok {
				return errors.New("pid field in json is not int")
			}
			content, ok = pJSON.Path("content").Data().(string)
			if !ok {
				return errors.New("content field in json is not string")
			}
			email, ok = pJSON.Path("email").Data().(string)
			if !ok {
				return errors.New("email field in json is not string")
			}
			return nil
		})
		if err != nil {
			return -1, fmt.Errorf("parse json in request: %v", err)
		}
		cid, err := d.InsertComment(pid, content, email)
		if err != nil {
			return -1, fmt.Errorf("insert comment: %v", err)
		}
		return cid, nil
	case "delete":
		var (
			cid int
		)
		err := parseJSONReq(r, func(pJSON *gabs.Container) error {
			var ok bool
			cid, ok = pJSON.Path("commentID").Data().(int)
			if !ok {
				return errors.New("commentID field in json is not string")
			}
			return nil
		})
		if err != nil {
			return -1, fmt.Errorf("parse json in request: %v", err)
		}
		performed, err := d.DeleteComment(cid)
		if err != nil {
			return -1, fmt.Errorf("delete comment: %v", err)
		}
		if !performed {
			return -1, errors.New("no matched comment found")
		}
		return -1, nil
	case "update":
		var (
			cid           int
			nContent, nAE string
		)
		err := parseJSONReq(r, func(pJSON *gabs.Container) error {
			var ok bool
			cid, ok = pJSON.Path("commentID").Data().(int)
			if !ok {
				return errors.New("commentID field in json is not int")
			}
			nContent, ok = pJSON.Path("newContent").Data().(string)
			if !ok {
				return errors.New("newContent field in json is not string")
			}
			nAE, ok = pJSON.Path("newEmail").Data().(string)
			if !ok {
				return errors.New("newEmail field in json is not string")
			}
			return nil
		})
		if err != nil {
			return -1, fmt.Errorf("parse json in request: %v", err)
		}
		performed, err := d.UpdateComment(cid, nContent, nAE)
		if err != nil {
			return -1, fmt.Errorf("update comment: %v", err)
		}
		if !performed {
			return -1, errors.New("no matched comment found")
		}
		return -1, nil
	default:
		return -1, errors.New("action is unknown")
	}

}

func changePost(d db.DB, action string, r *http.Request) (int, error) {
	switch action {
	case "insert":
		var (
			title, content string
			tags           []string
		)
		err := parseJSONReq(r, func(pJSON *gabs.Container) error {
			var ok bool
			title, ok = pJSON.Path("title").Data().(string)
			if !ok {
				return errors.New("title field in json is not string")
			}
			content, ok = pJSON.Path("content").Data().(string)
			if !ok {
				return errors.New("content field in json is not string")
			}
			tags, ok = pJSON.Path("tags").Data().([]string)
			if !ok {
				return errors.New("tags field in json is not string array")
			}
			return nil
		})
		if err != nil {
			return -1, fmt.Errorf("parse json in request: %v", err)
		}
		pid, err := d.InsertPost(title, content, tags)
		if err != nil {
			return -1, fmt.Errorf("insert post: %v", err)
		}
		return pid, nil
	case "delete":
		var (
			pid int
		)
		err := parseJSONReq(r, func(pJSON *gabs.Container) error {
			var ok bool
			pid, ok = pJSON.Path("pid").Data().(int)
			if !ok {
				return errors.New("pid field in json is not int")
			}
			return nil
		})
		if err != nil {
			return -1, fmt.Errorf("parse json in request: %v", err)
		}
		performed, err := d.DeletePost(pid)
		if err != nil {
			return -1, fmt.Errorf("delete post: %v", err)
		}
		if !performed {
			return -1, errors.New("no matched post found in db")
		}
		return -1, nil
	case "update":
		var (
			pid              int
			nTitle, nContent string
			nTags            []string
		)
		err := parseJSONReq(r, func(pJSON *gabs.Container) error {
			var ok bool
			pid, ok = pJSON.Path("pid").Data().(int)
			if !ok {
				return errors.New("pid field in json is not int")
			}
			nTitle, ok = pJSON.Path("newTitle").Data().(string)
			if !ok {
				return errors.New("newTitle field in json is not string")
			}
			nContent, ok = pJSON.Path("newContent").Data().(string)
			if !ok {
				return errors.New("newContent field in json is not string")
			}
			nTags, ok = pJSON.Path("newTags").Data().([]string)
			if !ok {
				return errors.New("newTags field in json is not string array")
			}
			return nil
		})
		if err != nil {
			return -1, fmt.Errorf("parse json in request: %v", err)
		}
		performed, err := d.UpdatePost(pid, nTitle, nContent, nTags)
		if err != nil {
			return -1, fmt.Errorf("update post: %v", err)
		}
		if !performed {
			return -1, fmt.Errorf("no matched post found in db")
		}
		return -1, nil
	default:
		return -1, errors.New("unknown action")
	}

}

func viewUser(d db.DB, r *http.Request) (*db.User, error) {
	var (
		passWord string
		userName string
	)
	err := parseJSONReq(r, func(pJSON *gabs.Container) error {
		var ok bool
		passWord, ok = pJSON.Path("passWord").Data().(string)
		if !ok {
			return errors.New("cannot parse passWord field in json as string")
		}
		userName, ok = pJSON.Path("userName").Data().(string)
		if !ok {
			return errors.New("cannot parse userName field in json as string")
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("cannot parse json in request: %v", err)
	}

	passHash := sha256.Sum256([]byte(passWord))
	usr, err := d.UserLogin(userName, passHash)
	if err != nil {
		return nil, fmt.Errorf("cannot login user: %v", err)
	}
	return usr, nil
}

func changeUser(d db.DB, action string, r *http.Request) (int, error) {
	switch action {
	case "register":
		var (
			uN, pW string
		)
		err := parseJSONReq(r, func(pJSON *gabs.Container) error {
			var ok bool
			uN, ok = pJSON.Path("userName").Data().(string)
			if !ok {
				return errors.New("cannot parse userName field as string in json")
			}

			pW, ok = pJSON.Path("passWord").Data().(string)
			if !ok {
				return errors.New("cannot parse passWord field as string in json")
			}

			return nil
		})

		if err != nil {
			return -1, fmt.Errorf("cannot parse json in request: %v", err)
		}

		uid, err := d.InsertUser(uN, sha256.Sum256([]byte(pW)))
		if err != nil {
			return -1, fmt.Errorf("cannot insert user: %v", err)
		}
		return uid, nil
	case "update":
		var (
			uid int
			nPW string
		)
		err := parseJSONReq(r, func(pJSON *gabs.Container) error {
			var ok bool
			uid, ok = pJSON.Path("uid").Data().(int)
			if !ok {
				return errors.New("cannot parse uid field in json as int")
			}

			nPW, ok = pJSON.Path("newPassWord").Data().(string)
			if !ok {
				return errors.New("cannot parse newPassWord field in json as string")
			}

			return nil
		})

		if err != nil {
			return -1, fmt.Errorf("cannot parse json in request: %v", err)
		}

		performed, err := d.UpdateUser(uid, sha256.Sum256([]byte(nPW)))
		if err != nil {
			return -1, fmt.Errorf("cannot update user: %v", err)
		}
		if !performed {
			return -1, errors.New("no matched user found")
		}
		return -1, nil
	default:
		return -1, fmt.Errorf("cannot perform action %s on resource: unknown action", action)
	}
}
