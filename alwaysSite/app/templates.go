package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"
)

type (
	User struct {
		Name     string
		Password string
		IsSafe   bool
	}

	IndexPageData struct {
		PageTitle     string
		Version       int64
		Curious       int
		CuriousExists bool
		Users         []User
		UsersExists   bool
	}
)

type tplModule struct {
	Tpls    map[string]*template.Template
	Curious int
	Users   []User
}

func NewTplModule() *tplModule {
	module := &tplModule{
		Tpls:  make(map[string]*template.Template),
		Users: []User{},
	}

	module.Tpls["index.html"] = template.Must(template.ParseFiles("front/templates/index.html"))

	return module
}

func (tm *tplModule) IndexHtml(w http.ResponseWriter, r *http.Request) {
	data := IndexPageData{
		PageTitle:     "Index Page",
		Version:       time.Now().Unix(),
		UsersExists:   len(tm.Users) > 0,
		Users:         tm.Users,
		CuriousExists: tm.Curious > 0,
		Curious:       tm.Curious,
	}
	tm.Tpls["index.html"].Execute(w, data)
}

func (tm *tplModule) RickRolledHandler(w http.ResponseWriter, r *http.Request) {
	tm.Curious++
}

func (tm *tplModule) Register(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read body: %s\n", err)
	}

	if err := json.Unmarshal(bytes, &body); err != nil {
		fmt.Printf("json.Unmarshal: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())

		return
	}

	tm.Users = append(tm.Users, User{
		Name:     body.Name,
		Password: body.Password,
		IsSafe:   len(body.Password) > 8,
	})

	io.WriteString(w, "Successfully added")
}
