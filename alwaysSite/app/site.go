package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
)

type obj map[string]interface{}

func main() {
	templates := NewTplModule()

	mux := http.NewServeMux()
	mux.HandleFunc("/", templates.IndexHtml)
	mux.HandleFunc("/hello", helloHandler)

	// mux.HandleFunc("/api/v1/rick", rickHandler)
	mux.Handle("/api/v1/rick", http.RedirectHandler("https://youtu.be/dQw4w9WgXcQ?si=ex-7DZgrTus1Vu8K", http.StatusSeeOther))

	mux.HandleFunc("/api/v1/register", templates.Register)
	mux.HandleFunc("/api/v1/rickRolled", templates.RickRolledHandler)
	mux.HandleFunc("/api/v1/say", sayHandler)
	mux.HandleFunc("/api/v1/calculate", calculateHandler)
	mux.HandleFunc("/api/v1/searchBooks", searchBooks)
	mux.HandleFunc("/api/v1/translateText", translateText)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8030"
	}

	ctx := context.Background()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, addrFromCtx, l.Addr().String())
			return ctx
		},
	}

	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error listening for server: %s\n", err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("helloHandler called")
	io.WriteString(w, "Hello, World!\n")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("rootHandler called")
	io.WriteString(w, "You are on main page!\n")
}

const addrFromCtx = "addrFromCtx"

// func rickHandler(w http.ResponseWriter, r *http.Request) {
// 	io.WriteString(w, "https://youtu.be/dQw4w9WgXcQ?si=ex-7DZgrTus1Vu8K")
// }

func sayHandler(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("name") {
		w.Header().Set("x-missing-field", "name")

		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Who the f are you?")

		return
	}

	name := r.URL.Query().Get("name")

	str := fmt.Sprintf("- Say my name.\n - %s.\n - You goddamn right!", name)

	io.WriteString(w, str)
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf("Method POST is required; You are using %s method", r.Method))

		return
	}

	var body struct {
		FirstNum  float64 `json:"first_num"`
		SecondNum float64 `json:"second_num"`
		Action    string  `json:"action"`
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

	switch strings.ToLower(body.Action) {
	case "addition":
		WriteJSON(w, obj{
			"message": body.Action,
			"result":  body.FirstNum + body.SecondNum,
		})

		return
	case "subtraction":
		WriteJSON(w, obj{
			"message": body.Action,
			"result":  body.FirstNum - body.SecondNum,
		})

		return
	case "multiplication":
		WriteJSON(w, obj{
			"message": body.Action,
			"result":  body.FirstNum * body.SecondNum,
		})

		return
	case "division":
		WriteJSON(w, obj{
			"message": body.Action,
			"result":  body.FirstNum / body.SecondNum,
		})

		return
	}

	w.WriteHeader(http.StatusBadRequest)
	WriteJSON(w, obj{
		"message": body.Action,
		"result":  "Unknown action",
	})
}

func WriteJSON(w http.ResponseWriter, obj interface{}) error {
	defer w.Header().Set("Content-Type", "application/json; charset=utf-8")
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = w.Write(jsonBytes)
	return err
}

type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

type BookResponse struct {
	Status string `json:"status"`
	Books  []Book `json:"books"`
}

func searchBooks(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	author := r.URL.Query().Get("author")

	books := []Book{
		{Title: "Go Programming", Author: "John Doe"},
		{Title: "Advanced Go Programming", Author: "John Doe"},
		{Title: "Dorohedoro", Author: "Q Hayashida"},
		{Title: "Captain Blood: His Odyssey", Author: "Rafael Sabatini"},
		{Title: "Harry Potter", Author: "Joanne Rowling"},
	}

	var filteredBooks []Book
	for _, book := range books {
		if (title == "" || book.Title == title) && (author == "" || book.Author == author) {
			filteredBooks = append(filteredBooks, book)
		}
	}

	if len(filteredBooks) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		resp := BookResponse{
			Status: "no books with this filter",
		}
		json.NewEncoder(w).Encode(resp)

		return
	}

	resp := BookResponse{
		Status: "success",
		Books:  filteredBooks,
	}
	json.NewEncoder(w).Encode(resp)
}

type TranslateRequest struct {
	Text           string `json:"text"`
	SourceLanguage string `json:"sourceLanguage"`
	TargetLanguage string `json:"targetLanguage"`
}

type TranslateResponse struct {
	Status         string `json:"status"`
	TranslatedText string `json:"translatedText"`
}

func translateText(w http.ResponseWriter, r *http.Request) {
	var req TranslateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := TranslateResponse{
		Status:         "success",
		TranslatedText: fmt.Sprintf("Imagine that \"%s\" translated from \"%s\" to \"%s\". There is almost none of free-to-use api translator...", req.Text, req.SourceLanguage, req.TargetLanguage),
	}
	json.NewEncoder(w).Encode(resp)
}
