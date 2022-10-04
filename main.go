package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

func addUrl(w http.ResponseWriter, r *http.Request) {
	newUrl := mux.Vars(r)
	temp := fmt.Sprint(newUrl["url"])
	fmt.Println(temp, MakeShortUrl(temp))
}

func MakeShortUrl(url string) string {
	tempUrl := md5.Sum([]byte(url))
	hashUrl := hex.EncodeToString(tempUrl[:])
	shortHashUrl := hashUrl[:len(hashUrl)-25]
	return shortHashUrl
}

func main() {
	router := mux.NewRouter()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	defer server.Close()

	http.Handle("/", router)
	router.HandleFunc("/", someTempFunc)
	router.HandleFunc("/save", saveUrl)
	router.HandleFunc("/geturl", getShortUrl)
	router.HandleFunc("/add/{url}", addUrl)

	// fmt.Println("Server starting on port: 8080")
	// http.ListenAndServe(":8080", nil)

	go func() {
		fmt.Println("Server starting on port: 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println(err)
		}
	}()

	<-stop

	fmt.Println(" Server was stoped")

}
