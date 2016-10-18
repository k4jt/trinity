package server

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/k4jt/trinity/store"
	"html/template"
	"net/http"
	"path"
	"path/filepath"
)

func Index(ctx *Context, w http.ResponseWriter, r *http.Request) {
	users, _ := ctx.DB.GetAllUsers()
	ctx.Payload = struct {
		ActivePage string
		Users      []store.User
	}{
		ActivePage: "index",
		Users:      users,
	}

	render(ctx, w, "index")
}

func render(ctx *Context, w http.ResponseWriter, filenames ...string) {
	t := combineTemplates(filenames...)

	err := t.Execute(w, ctx.Payload)
	if err != nil {
		log.Println("template executing error: ", err)
	}
}

func combineTemplates(filenames ...string) *template.Template {
	b, err := Asset("index.html")
	if err != nil {
		log.Println("index.html getting asset error: ", err)
	}
	t := template.New("index.html")
	_, err = t.Parse(string(b))
	if err != nil {
		log.Println("template parsing error: ", err)
	}
	for _, url := range filenames {
		filename := path.Join("templates", fmt.Sprintf("%s.html", url))
		b, err = Asset(filename)
		if err != nil {
			log.Println(filename, " getting asset error: ", err)
			break
		}

		name := filepath.Base(filename)
		var tmpl *template.Template
		if t == nil {
			t = template.New(name)
		}
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name)
		}
		_, err = tmpl.Parse(string(b))
		if err != nil {
			log.Println(filename, " template parsing error: ", err)
			break
		}
	}
	return t
}
