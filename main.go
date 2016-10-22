//go:generate go-bindata -debug -prefix static/ -pkg server -o server/static_gen.go static/...

package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/k4jt/trinity/server"
	"github.com/k4jt/trinity/store"
	"github.com/pkg/browser"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var version string // set by the compiler

func getHttpHandler(c *cli.Context, db *store.Store) http.Handler {
	ctx := server.NewContext(c, db)
	httpHandler := server.NewRouter(ctx)
	httpHandler.PathPrefix("/").Handler(http.FileServer(&assetfs.AssetFS{
		Asset:     server.Asset,
		AssetDir:  server.AssetDir,
		AssetInfo: server.AssetInfo,
		Prefix:    "",
	}))

	return httpHandler
}

func run(c *cli.Context) error {
	log.SetLevel(log.Level(uint8(c.Int("log-level"))))
	log.WithField("version", version).Info("starting trinity db")

	db := store.Open()
	defer db.Close()

	httpHandler := getHttpHandler(c, db)
	go func() {
		log.WithField("bind", c.String("http-bind")).Info("starting frontend server (without tls)")
		log.Fatal(http.ListenAndServe(c.String("http-bind"), httpHandler))
	}()

	url := "http://" + c.String("http-bind")
	browser.OpenURL(url)

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	log.WithField("signal", <-sigChan).Info("signal received")
	log.Info("shutting down server")
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "frontend"
	app.Version = version
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "http-bind",
			Usage:  "ip:port to bind the http server (web-interface)",
			Value:  "0.0.0.0:4444",
			EnvVar: "HTTP_BIND",
		},
		cli.IntFlag{
			Name:   "log-level",
			Value:  4,
			Usage:  "debug=5, info=4, warning=3, error=2, fatal=1, panic=0",
			EnvVar: "LOG_LEVEL",
		},
	}
	app.Run(os.Args)
}
