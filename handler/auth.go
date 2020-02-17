package handler

import (
	"middleware/handler/db"
	"net/http"

	"github.com/google/uuid"
)

var userTable = make(map[string]*db.User, 1000)

func newUUID() *http.Cookie {
	return &http.Cookie{Name: "uuid", Value: uuid.New().String()}
}
