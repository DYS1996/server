package main

import (
	"log"
	"middleware/handler"
	"middleware/pgsql"
	"net/http"
	"os"
	"path/filepath"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf("resolve file abs path: %v", err)
	}

	db, err := pgsql.New(&pgsql.PGConfig{User: "alncl", Pass: "l;n1223", DBName: "blog"})
	if err != nil {
		log.Fatalf("setup db server: %v\n", err)
	}
	defer db.Close()

	// srvConfig := &SrvConfig{Host: "172.31.41.201", Port: 8443}
	srv := http.Server{
		Addr:    ":8443",
		Handler: handler.New(db),
	}

	crt := autocert.NewListener("api.redhand.vip")

	if err := srv.Listen(crt); err != nil && err != http.ErrServerClosed {
		log.Fatalf("setup server: %v\n", err)
	}
}
