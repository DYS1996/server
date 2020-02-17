package main

import (
	"log"
	"middleware/handler"
	"middleware/pgsql"
	"net/http"
	"os"
	"path/filepath"
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

	if err := srv.ListenAndServeTLS(filepath.Join(dir, "./ssl.crt"), filepath.Join(dir, "./ssl.key")); err != nil {
		log.Fatalf("setup server: %v\n", err)
	}
}
