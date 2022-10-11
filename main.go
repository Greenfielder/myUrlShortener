package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type urlline struct {
	id       int
	url      string
	shorturl string
}

func MakeShortUrl(url string) string {
	tempUrl := md5.Sum([]byte(url))
	hashUrl := hex.EncodeToString(tempUrl[:])
	shortHashUrl := hashUrl[:len(hashUrl)-25]
	return shortHashUrl
}

func main() {
	router := mux.NewRouter()

	// db, err := sql.Open("sqlite3", "urldata.db")
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()

	// temp := "mail.ru"
	// sqlInit := `create table if not exists urlslist (id integer not null primary key, url text, shorturl text);`
	// _, err = db.Exec(sqlInit)
	// if err != nil {
	// 	log.Fatal(err)
	// }

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
	router.HandleFunc("/geturl", getAll)
	router.HandleFunc("/add/{url}", addUrl)
	router.HandleFunc("/get/{url}", getUrl)
	router.HandleFunc("/{url}", getShortUrl)

	go func() {
		fmt.Println("Server starting on port: 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println(err)
		}
	}()

	<-stop
	fmt.Println(" Server was stoped")
}

func someTempFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello! This is Main page.")
	ip := strings.Split(r.RemoteAddr, " :")
	xtype := fmt.Sprintf("%T", ip)
	fmt.Println(xtype)
}

//func getAll returns all urls
func getAll(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "urldata.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("select * from urlslist")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	uls := []urlline{}

	for rows.Next() {
		ul := urlline{}
		err := rows.Scan(&ul.id, &ul.url, &ul.shorturl)
		if err != nil {
			fmt.Println(err)
			continue
		}
		uls = append(uls, ul)
	}

	for _, ul := range uls {
		fmt.Println(ul.id, ul.url, ul.shorturl)
	}

}

// func addUrl Add url, like: http://localhost:8080/add/mail.ru
func addUrl(w http.ResponseWriter, r *http.Request) {
	newUrl := mux.Vars(r)
	temp := fmt.Sprint(newUrl["url"])
	fmt.Println(temp, MakeShortUrl(temp))

	db, err := sql.Open("sqlite3", "urldata.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	result, err := db.Exec("insert into urlslist (url, shorturl) values ($1, $2)", temp, MakeShortUrl(temp))
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}

// func getUrl get url, like: http://localhost:8080/get/mail.ru
func getUrl(w http.ResponseWriter, r *http.Request) {
	newUrl := mux.Vars(r)
	temp := fmt.Sprint(newUrl["url"])
	fmt.Println(temp, MakeShortUrl(temp))

	db, err := sql.Open("sqlite3", "urldata.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	row := db.QueryRow("select * from urlslist where url = $1", temp)
	uls := urlline{}
	err = row.Scan(&uls.id, &uls.url, &uls.shorturl)
	if err != nil {
		panic(err)
	}
	fmt.Println(uls.shorturl)
}

// func getShortUrl get shorl url, like: http://localhost:8080/b10e94d And redirect to https://mail.ru
func getShortUrl(w http.ResponseWriter, r *http.Request) {
	newUrl := mux.Vars(r)
	temp := fmt.Sprint(newUrl["url"])
	fmt.Println(temp)

	db, err := sql.Open("sqlite3", "urldata.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	row := db.QueryRow("select * from urlslist where shorturl = $1", temp)
	uls := urlline{}
	err = row.Scan(&uls.id, &uls.url, &uls.shorturl)
	if err != nil {
		panic(err)
	}
	fmt.Println(uls.url)
	final := "https://" + uls.url
	http.Redirect(w, r, final, http.StatusSeeOther)
}
