package db

import (
	"crypto/sha256"
)

// DB lists essential methods for the use of blog server
type DB interface {
	GetPostByID(id int) (*Post, error)
	GetPosts(pageSize, page int) ([]Post, error)
	GetPostsCount() (int, error)
	GetPostsByFTS(search string, pageSize, page int) ([]Post, error)
	GetPostsCountByFTS(search string) (int, error)
	UserLogin(user string, pass [sha256.Size]byte) (*User, error)
	InsertPost(title string, content string, tags []string) (int, error)
	DeletePost(pid int) (bool, error)
	UpdatePost(pid int, nTitle, nContent string, nTags []string) (bool, error)
	GetCommentsCount(pid int) (int, error)
	GetCommentsByPage(pid int, pageSize, page int) ([]Comment, error)
	InsertComment(pid int, content, authorEmail string) (int, error)
	DeleteComment(cid int) (bool, error)
	UpdateComment(cid int, nContent, nAE string) (bool, error)
	GetUser(userName string, passWord [sha256.Size]byte) (*User, error)
	InsertUser(userName string, passWord [sha256.Size]byte) (int, error)
	UpdateUser(uid int, nPW [sha256.Size]byte) (bool, error)
}
