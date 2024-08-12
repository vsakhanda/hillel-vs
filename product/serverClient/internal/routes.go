// using net and http without gin-gonic

package routes

import (
	"fmt"
	"io"
	"net/http"
)

const addrFromCtx = "addrFromCtx"

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", IndexHandler)
	mux.HandleFunc("/api/data", ApiDataHandler)
	mux.HandleFunc("/hello", helloHendler)
	mux.HandleFunc("/api/v1/rick", rickHandler)
	mux.HandleFunc("/api/v1/say", sayHandler)
	mux.HandleFunc("/api/v1/searchBooks", searchBooksHandler)

	return mux
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintln(w, "Welcome to homepage")
	if err != nil {
		return
	}
}

func ApiDataHandler(w http.ResponseWriter, r *http.Request) {
	data := "Some data from the API"
	fmt.Fprintln(w, data)
}
func helloHendler(w http.ResponseWriter, r *http.Request) {
	data := "Hello new world! with new IDE and productivity"
	fmt.Fprintln(w, data)
}
func rickHandler(w http.ResponseWriter, r *http.Request) {
	data := "Hello Rick! Where is morty?"
	fmt.Fprintln(w, data)
}

func sayHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Printf("%s: got /hello request\n", ctx.Value(addrFromCtx))

	myName := r.PostFormValue("myName")
	if myName == "" {
		w.Header().Set("x-missing-field", "myName")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err := io.WriteString(w, fmt.Sprintf("Hello, %s!\n", myName))
	if err != nil {
		return
	}
}

func searchBooksHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method == http.MethodPost {
		fmt.Println("Used method POST")
	} else {
		fmt.Println("Used method", r.Method)
	}

	title := r.URL.Query().Get("title")
	author := r.URL.Query().Get("author")

	fmt.Printf("%s: got / request. title=%s. author=%s.\n",
		ctx.Value(addrFromCtx),
		title, author)
	io.WriteString(w, fmt.Sprintf("This is your liblrary. \n Your search includes Book '%s'  of author '%s'", title, author))
}
