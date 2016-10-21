package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/k4jt/trinity/store"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strconv"
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

func DeleteUser(ctx *Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	id := vars["id"]
	log.Info("id:", id)
	i32, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := ctx.DB.DeleteUser(uint64(i32)); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusFound)
}

func Search(ctx *Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	q := vars["q"]
	log.Info("q:", q)

	q = strings.TrimSpace(q)
	q = strings.ToLower(q)

	users, _ := ctx.DB.GetAllUsers()

	if len(q) > 0 {
		filtered := []store.User{}
		for _, u := range users {
			if u.Contains(q) {
				filtered = append(filtered, u)
			}
		}
		users = filtered
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Error("Error json encode users:", err)
	}

	//content, _ := ioutil.ReadAll(r.Body)
	//log.Println("Raw POST content: ", string(content))
	//params := decodeParams(content)
	//log.Info("params: ", params)

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
