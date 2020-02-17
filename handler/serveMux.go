package handler

import (
	"context"
	"errors"
	"fmt"
	"middleware/handler/db"
	"net/http"

	"github.com/Jeffail/gabs/v2"
	uuidLib "github.com/google/uuid"
)

// New inits a http handler with functions of preprocessing, postprocessing and servemux
func New(d db.DB) http.Handler {
	var ServeMux = http.NewServeMux()

	ServeMux.Handle(`/post`, handlerWrapper(func(w http.ResponseWriter, r *http.Request) http.Handler {
		switch r.Method {
		case http.MethodGet:
			post, err := viewPost(d, r)
			if err != nil {
				return Err{fmt.Errorf("process request: %v", err)}
			}
			return JSONData{post}
		case http.MethodPost:
			usr, ok := r.Context().Value(db.BlogContext("user")).(*db.User)
			if !ok {
				return Err{errors.New("no user context: internal error")}
			}
			if usr == nil {
				return Err{errors.New("user unlogined")}
			}
			if usr.Privilege != 0 {
				return Err{errors.New("not enough privilege")}
			}
			var action string
			err := parseJSONReq(r, func(pJSON *gabs.Container) error {
				var ok bool
				action, ok = pJSON.Path("action").Data().(string)
				if !ok {
					return errors.New("parse action field in json is not string")
				}
				return nil
			})
			if err != nil {
				return Err{fmt.Errorf("parse json request: %v", err)}
			}
			pid, err := changePost(d, action, r)
			if err != nil {
				return Err{fmt.Errorf("change post : %v", err)}
			}
			return JSONData{pid}
		default:
			return Err{errors.New("request method is not POST/GET")}
		}
	}))

	ServeMux.Handle(`/posts`, handlerWrapper(func(w http.ResponseWriter, r *http.Request) http.Handler {
		switch r.Method {
		case http.MethodGet:
			pstsPage, err := viewPosts(d, r)
			if err != nil {
				return Err{fmt.Errorf("preocess request: %v", err)}
			}
			return JSONData{pstsPage}
		default:
			return Err{errors.New("request method is not GET")}
		}
	}))
	ServeMux.Handle(`/comments`, handlerWrapper(func(w http.ResponseWriter, r *http.Request) http.Handler {
		switch r.Method {
		case http.MethodGet:
			cmtsPage, err := viewComments(d, r)
			if err != nil {
				return Err{fmt.Errorf("process request: %v", err)}
			}
			return JSONData{cmtsPage}
		default:
			return Err{errors.New("request method is not GET")}
		}
	}))

	ServeMux.Handle(`/comment`, handlerWrapper(func(w http.ResponseWriter, r *http.Request) http.Handler {
		switch r.Method {
		case http.MethodPost:
			usr, ok := r.Context().Value(db.BlogContext("user")).(*db.User)
			if !ok {
				return Err{errors.New("no user context: internal error")}
			}
			if usr == nil {
				return Err{errors.New("user unlogined")}
			}
			if usr.Privilege != 0 {
				return Err{errors.New("not enough privilege")}
			}
			var (
				action string
			)
			err := parseJSONReq(r, func(pJSON *gabs.Container) error {
				var ok bool
				action, ok = pJSON.Path("action").Data().(string)
				if !ok {
					return errors.New("action field in json is not string")
				}
				return nil
			})
			if err != nil {
				return Err{fmt.Errorf("parse json in request: %v", err)}
			}
			cid, err := changeComment(d, action, r)
			if err != nil {
				return Err{fmt.Errorf("change comment: %v", err)}
			}
			return JSONData{cid}
		default:
			return Err{errors.New("request method is not POST")}
		}
	}))

	ServeMux.Handle(`/user`, handlerWrapper(func(w http.ResponseWriter, r *http.Request) http.Handler {
		switch r.Method {
		case http.MethodPost:
			var (
				action string
			)
			err := parseJSONReq(r, func(pJSON *gabs.Container) error {
				var ok bool
				action, ok = pJSON.Path("action").Data().(string)
				if !ok {
					return errors.New("action field in json is not string")
				}
				return nil
			})
			if err != nil {
				return Err{fmt.Errorf("parse json: %v", err)}
			}

			switch action {
			case "login":
				usr, err := viewUser(d, r)
				if err != nil {
					return Err{fmt.Errorf("login user: %v", err)}
				}

				uuid, ok := r.Context().Value(db.BlogContext("uuid")).(*UUID)
				if !ok {
					return Err{errors.New("no uuid context: internal error")}
				}

				// uuidVal := uuid.Val
				if !uuid.New {
					nUUID := newUUID()
					http.SetCookie(w, nUUID)
					uuid.Val = nUUID.Value
					uuid.New = true
				}

				userTable[uuid.Val] = &db.User{UID: usr.UID, UserName: usr.UserName, Privilege: usr.Privilege}

				return JSONData{usr}

			default:
				uid, err := changeUser(d, action, r)

				if err != nil {
					return Err{fmt.Errorf("change user info: %v", err)}
				}
				return JSONData{uid}
			}
		default:
			return Err{errors.New("request method is not POST")}
		}
	}))

	ServeMux.HandleFunc(`/ping`, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
		user, ok := r.Context().Value(db.BlogContext("user")).(*db.User)
		if ok && user != nil {
			w.Write([]byte(", " + user.UserName))
		}
		return
	})

	return postProcess(preProcess(ServeMux))
}

func preProcess(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uuid, err := r.Cookie("uuid")

		if err != nil {
			nUUID := newUUID()
			http.SetCookie(w, nUUID)
			uuid = nUUID
			ctx := context.WithValue(r.Context(), db.BlogContext("uuid"), &UUID{uuid.Value, true})
			ctx = context.WithValue(ctx, db.BlogContext("user"), (*db.User)(nil))

			h.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		_, err = uuidLib.Parse(uuid.Value)
		if err != nil {
			nUUID := newUUID()
			http.SetCookie(w, nUUID)
			uuid = nUUID
			ctx := context.WithValue(r.Context(), db.BlogContext("uuid"), &UUID{uuid.Value, true})

			ctx = context.WithValue(ctx, db.BlogContext("user"), (*db.User)(nil))
			h.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		user := userTable[uuid.Value]

		ctx := context.WithValue(r.Context(), db.BlogContext("uuid"), &UUID{uuid.Value, false})
		ctx = context.WithValue(ctx, db.BlogContext("user"), user)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func postProcess(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(`X-Content-Type-Options`, `nosniff`)
		w.Header().Set(`Cache-Control`, `no-cache, no-store, must-revalidate`)
		w.Header().Set(`Access-Control-Allow-Origin`, `https://www.redhand.vip`)
		w.Header().Set(`Connection`, `close`)
		h.ServeHTTP(w, r)
	})
}
