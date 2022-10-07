package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type sqlite struct {
	Path string
}

type Database interface {
	GetShortUrl(id int) (string, error)
	AddUrl(url string, shorturl string) (int, error)
}

func (s sqlite) Init() {
	db, err := sql.Open("sqlite3", s.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlInit := `create table if not exists urls (id integer not null primary key, url text, shorturl text);`
	_, err = db.Exec(sqlInit)
	if err != nil {
		log.Fatal(err)
	}
}

func (s sqlite) AddUrl(url string, shortUrl string) (int, error) {
	db, err := sql.Open("sqlite3", s.Path)
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	stmt, err := tx.Prepare("insert into urls(url) values(?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(url)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	tx.Commit()

	return id, nil
}

func (s sqlite) GetShortUrl(url string) (string, error) {
	db, err := sql.Open("sqlite3", s.Path)
	stmt, err := db.Prepare("select shorturl from urls where url = ?")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	var shorturl string
	err = stmt.QueryRow(url).Scan(&shorturl)
	if err != nil {
		return "", err
	}
	return url, nil
}
