package main

import (
	"database/sql"
	"log"
)

type sqlite struct {
	Path string
}

func (s sqlite) Init() {
	db, err := sql.Open("sqlite3", s.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
