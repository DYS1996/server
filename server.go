package main

import (
        "log"
        "middleware/handler"
        "middleware/pgsql"
        "net/http"
        "golang.org/x/crypto/acme/autocert"
)

func main() {

        db, err := pgsql.New(&pgsql.PGConfig{User: "blogdbu", Pass: "123s;,nl", DBName: "blog"})
        if err != nil {
                log.Fatalf("setup db server: %v\n", err)
        }
        defer db.Close()

        // srvConfig := &SrvConfig{Host: "172.31.41.201", Port: 8443}
        srv := http.Server{
                Addr:    ":443",
                Handler: handler.New(db),
        }

        crt := autocert.NewListener("api.redhand.vip")

        if err := srv.Serve(crt); err != nil && err != http.ErrServerClosed {
                log.Fatalf("setup server: %v\n", err)
        }
}