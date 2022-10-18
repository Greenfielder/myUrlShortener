package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type urlline struct {
	id       int
	url      string
	shorturl string
}

var tpl = template.Must(template.ParseFiles("index.html"))
var tpl2 = template.Must(template.ParseFiles("done.html"))
var Done string
var HostName string = "http://localhost:8080/"

func main() {
	initDb()

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

	fs := http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/")))
	router.PathPrefix("/assets/").Handler(fs)
	router.HandleFunc("/", mainPage)
	router.HandleFunc("/geturl", getAll)
	router.HandleFunc("/add", addUrl)
	router.HandleFunc("/{url}", shortUrtToOrig)

	go func() {
		fmt.Println("Server starting on port: 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println(err)
		}
	}()

	<-stop
	fmt.Println(" Server was stoped")
}

// func MakeShortUrl convert url to hash
func MakeShortUrl(url string) string {
	tempUrl := md5.Sum([]byte(url))
	hashUrl := hex.EncodeToString(tempUrl[:])
	shortHashUrl := hashUrl[:len(hashUrl)-25]
	return shortHashUrl
}

// func initDb is initialize database
func initDb() {
	db, err := sql.Open("sqlite3", "urldata.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlInit := `create table if not exists urlslist (id integer not null primary key, url text, shorturl text);`
	_, err = db.Exec(sqlInit)
	if err != nil {
		log.Fatal(err)
	}
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

// func addUrl, add url to database and return his short variant
func addUrl(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	params := u.Query()
	searchKey := params.Get("q")
	Done = HostName + MakeShortUrl(searchKey)
	fmt.Sprintln(searchKey, Done)

	db, err := sql.Open("sqlite3", "urldata.db")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	result, err := db.Exec("insert into urlslist (url, shorturl) values ($1, $2)", searchKey, MakeShortUrl(searchKey))
	if err != nil {
		log.Println(err)
	}
	fmt.Println(result)

	tpl2.ExecuteTemplate(w, "done.html", struct{ Test string }{Test: Done})
}

// func getShortUrl get shorl url, like: http://localhost:8080/b10e94d And redirect to https://mail.ru
func shortUrtToOrig(w http.ResponseWriter, r *http.Request) {
	newUrl := mux.Vars(r)
	adress := fmt.Sprint(newUrl["url"])

	db, err := sql.Open("sqlite3", "urldata.db")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	row := db.QueryRow("select * from urlslist where shorturl = $1", adress)
	uls := urlline{}
	err = row.Scan(&uls.id, &uls.url, &uls.shorturl)
	if err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, uls.url, http.StatusSeeOther)
}

//func getAll returns all urls
func getAll(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "urldata.db")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	rows, err := db.Query("select * from urlslist")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	uls := []urlline{}

	for rows.Next() {
		ul := urlline{}
		err := rows.Scan(&ul.id, &ul.url, &ul.shorturl)
		if err != nil {
			log.Println(err)
			continue
		}
		uls = append(uls, ul)
	}

	for _, ul := range uls {
		fmt.Println(ul.id, ul.url, ul.shorturl)
	}
}
