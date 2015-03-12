package model

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB = nil

func Init() {
	var err error
	db, err = sql.Open("sqlite3", "file:database.sqlite?cache=shared")
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(16)

	setupTables()
}

func Deinit() {
	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
}

func setupTables() {
	var sql string

	sql = `
		CREATE TABLE IF NOT EXISTS config (
			key VARCHAR(128) PRIMARY KEY NOT NULL,
			value BLOB
		);

		INSERT OR IGNORE INTO config
		(key, value)
		VALUES
		("dbVersion", 0);`
	if _, err := db.Exec(sql); err != nil {
		log.Fatal(err)
	}

	sql = `
		SELECT value
		FROM config
		WHERE key = "dbVersion";`
	var dbVersion int
	if err := db.QueryRow(sql).Scan(&dbVersion); err != nil {
		log.Fatal(err)
	}

	sql = ""
	switch dbVersion {
	case 0:
		sql += `
				CREATE TABLE IF NOT EXISTS comment (
					id INTEGER PRIMARY KEY ASC NOT NULL,
					site TEXT NOT NULL,
					text TEXT NOT NULL,
					ctime DATETIME NOT NULL,
					mtime DATETIME NOT NULL,
					website TEXT NOT NULL,
					author TEXT NOT NULL,
					email TEXT NOT NULL
				);

				CREATE INDEX comment_site_ctime
				ON comment (site, ctime);`
		//fallthrough
	}

	sql += `
		UPDATE config
		SET value = 1
		WHERE key = "dbVersion";`

	if _, err := db.Exec(sql); err != nil {
		log.Fatal(err)
	}
}
