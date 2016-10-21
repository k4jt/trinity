package server

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/gorilla/mux"
	"github.com/k4jt/trinity/store"
	"net/http"
	"time"
)

type Context struct {
	RunContext *cli.Context
	DB         *store.Store
	Payload    interface{}
}

func NewContext(c *cli.Context, db *store.Store) *Context {
	return &Context{RunContext: c, DB: db}
}

type ctxHandler func(*Context, http.ResponseWriter, *http.Request)

func Logger(ctx *Context, inner ctxHandler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner(ctx, w, r)

		log.Printf(
			"%-6s %-10s\t%-10s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

func NewRouter(ctx *Context) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {

		handler := Logger(ctx, route.HandlerFunc, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

type Route struct {
	Name, Method, Pattern string
	HandlerFunc           ctxHandler
}

var routes = []Route{
	{
		"Index",
		"GET",
		"/",
		Index,
	},
	{
		"AddUser",
		"POST",
		"/",
		AddUser,
	},
	//{
	//	"EditUser",
	//	"PUT",
	//	"/{id}",
	//	EditUser,
	//},
	{
		"DeleteUser",
		"DELETE",
		"/{id}",
		DeleteUser,
	},
	{
		"Search",
		"GET",
		"/search",
		Search,
	},
	{
		"Search",
		"GET",
		"/search/{q}",
		Search,
	},
}
