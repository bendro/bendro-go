package model

import (
	"database/sql"
	"time"
)

type Comment struct {
	Id      int64     `json:"id"`
	Text    string    `json:"text,omitempty"`
	Ctime   time.Time `json:"ctime,omitempty"`
	Mtime   time.Time `json:"mtime,omitempty"`
	Website string    `json:"website,omitempty"`
	Author  string    `json:"author,omitempty"`
	EMail   string    `json:"email,omitempty"`
	Site    string    `json:"site,omitempty"`
}

func GetComments(quary string, quaryArgs ...interface{}) (comments []*Comment, err error) {
	err = tx(func(t *sql.Tx) {
		rows, err := t.Query(`
			SELECT id, text, ctime, mtime, website, author, email, site
			FROM comment
			`+quary+";",
			quaryArgs...,
		)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			c := &Comment{}

			err := rows.Scan(
				&c.Id,
				&c.Text,
				&c.Ctime,
				&c.Mtime,
				&c.Website,
				&c.Author,
				&c.EMail,
				&c.Site,
			)
			if err != nil {
				panic(err)
			}

			comments = append(comments, c)
		}
		if err := rows.Err(); err != nil {
			panic(err)
		}
	})

	return
}

func CreateComment(comment *Comment) (err error) {
	now := time.Now()
	comment.Ctime = now
	comment.Mtime = now

	err = tx(func(t *sql.Tx) {
		res, err := t.Exec(`
			INSERT INTO comment
			(text, ctime, mtime, website, author, email, site)
			VALUES
			(?, ?, ?, ?, ?, ?, ?);`,

			comment.Text,
			comment.Ctime,
			comment.Mtime,
			comment.Website,
			comment.Author,
			comment.EMail,
			comment.Site,
		)
		if err != nil {
			panic(err)
		}

		comment.Id, err = res.LastInsertId()
		if err != nil {
			panic(err)
		}
	})

	return
}

func EditComment(comment *Comment) (err error) {
	err = tx(func(t *sql.Tx) {
		_, err := t.Exec(`
			UPDATE comment
			SET
				text = ?, ctime = ?, mtime = ?, website = ?, author = ?, email = ?,
				site = ?,
			WHERE id = ?;`,

			comment.Text,
			comment.Ctime,
			comment.Mtime,
			comment.Website,
			comment.Author,
			comment.EMail,
			comment.Site,
			comment.Id,
		)

		if err != nil {
			panic(err)
		}
	})

	return
}

func DeleteComment(id int) (err error) {
	err = tx(func(t *sql.Tx) {
		_, err := t.Exec(`
		DELETE FROM comment
		WHERE id = ?;`,
			id,
		)

		if err != nil {
			panic(err)
		}
	})

	return
}
