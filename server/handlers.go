package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/k4jt/trinity/store"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"time"
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

func decodeParams(content []byte) map[string]interface{} {
	params := map[string]interface{}{}
	if err := json.NewDecoder(bytes.NewReader(content)).Decode(&params); err != nil {
		log.Println("Error while decoding request params. But we try to handle it.")

		data := strings.Split(string(content), "&")
		for _, pair := range data {
			values := strings.Split(pair, "=")
			params[values[0]] = values[1]
		}

	}

	return params
}

func AddUser(ctx *Context, w http.ResponseWriter, r *http.Request) {

	content, _ := ioutil.ReadAll(r.Body)
	log.Println("Raw POST content: ", string(content))
	params := decodeParams(content)
	log.Info("params: ", params)

	t, _ := time.Parse("02-01-2006", params["birth_day"].(string))
	name, _ := url.QueryUnescape(params["name"].(string))
	family, _ := url.QueryUnescape(params["family"].(string))
	address, _ := url.QueryUnescape(params["address"].(string))
	phone, _ := url.QueryUnescape(params["phone"].(string))

	u := store.User{
		Name:         name,
		Family:       family,
		BirthDayFull: t,
		Address:      address,
		Phone:        []string{phone},
	}

	ctx.DB.CreateUser(&u)

	Index(ctx, w, r)
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
