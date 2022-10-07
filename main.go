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

	_ "github.com/mattn/go-sqlite3"

	"github.com/gorilla/mux"
)

func someTempFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello! This is Main page.")
}

func getShortUrl(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "We showing short URL")
}

func addUrl(w http.ResponseWriter, r *http.Request) {
	newUrl := mux.Vars(r)
	temp := fmt.Sprint(newUrl["url"])
	Database.AddUrl(temp, MakeShortUrl(temp))
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

	path := "./sqlite-database.db"
	db := sqlite{Path: path}
	db.Init()

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

	router.HandleFunc("/", someTempFunc)
	router.HandleFunc("/geturl", getShortUrl)
	router.HandleFunc("/add/{url}", addUrl)

	go func() {
		fmt.Println("Server starting on port: 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println(err)
		}
	}()

	// log.Println("Creating sqlite-database.db...")
	// file, err := os.Create("sqlite-database.db") // Create SQLite file
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// file.Close()
	// log.Println("sqlite-database.db created")

	// sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.db") // Open the created SQLite File
	// defer sqliteDatabase.Close()                                     // Defer Closing the database

	<-stop
	fmt.Println(" Server was stoped")
}
