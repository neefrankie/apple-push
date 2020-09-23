package main

import (
	"flag"
	"fmt"
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/render"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/neefrankie/apple-push/internal/repo"
	"github.com/neefrankie/apple-push/pkg/config"
	"github.com/neefrankie/apple-push/pkg/db"
	"github.com/neefrankie/apple-push/pkg/message"
	"log"
	"net/http"
	"os"
)

var (
	version    string
	build      string
	production bool
)

func init() {
	flag.BoolVar(&production, "production", false, "Connect to production MySQL database if present. Default to localhost.")
	var v = flag.Bool("v", false, "print current version")

	flag.Parse()

	if *v {
		fmt.Printf("%s\nBuild at %s\n", version, build)
		os.Exit(0)
	}

	config.MustSetupViper()
}

func main() {
	logger := config.MustGetLogger(production)

	myDB := db.MustNewDB(config.MustDBConn(production))

	rp := repo.New(myDB, logger)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/message", func(writer http.ResponseWriter, request *http.Request) {
		var m message.Message
		if err := gorest.ParseJSON(request.Body, &m); err != nil {
			_ = render.New(writer).BadRequest(err.Error())
			return
		}

		go rp.Push(&m)

		_ = render.New(writer).NoContent()
	})

	log.Fatal(http.ListenAndServe("9002", r))
}
