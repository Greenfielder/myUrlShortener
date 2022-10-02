package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func someTempFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Main page")
}

func saveUrl(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Saving url")
}

func getShortUrl(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "We showing short URL")
}

func main() {
	router := mux.NewRouter()

	http.Handle("/", router)
	router.HandleFunc("/", someTempFunc)
	router.HandleFunc("/save", saveUrl)
	router.HandleFunc("/geturl", getShortUrl)

	fmt.Println("Server starting on port: 8080")
	http.ListenAndServe(":8080", nil)
}
