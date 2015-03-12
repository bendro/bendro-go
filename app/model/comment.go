package model

import (
	"database/sql"
	"time"
)

type Comment struct {
	Id         int64     `json:"id"`
	Site       string    `json:"site,omitempty"`
	ThreadId   int64     `json:"threadId"`
	ResponseTo *int64    `json:"responseTo"`
	Left       int64     `json:"left"`
	Right      int64     `json:"right"`
	Ctime      time.Time `json:"ctime,omitempty"`
	Mtime      time.Time `json:"mtime,omitempty"`
	Text       string    `json:"text,omitempty"`
	Website    string    `json:"website,omitempty"`
	Author     string    `json:"author,omitempty"`
	EMail      string    `json:"email,omitempty"`
}

func GetComments(site string) (comments []*Comment, err error) {
	err = tx(func(t *sql.Tx) {
		rows, err := t.Query(`
			SELECT
				c.id, c.site, c.thread_id, c.response_to, c.left, c.right, c.ctime,
				c.mtime, c.text, c.website, c.author, c.email
			FROM comment AS c
			INNER JOIN comment AS fc ON fc.id = c.thread_id
			WHERE c.site = ?
			ORDER BY fc.ctime, fc.id, c.left;`,
			site,
		)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			c := &Comment{}

			err := rows.Scan(
				&c.Id,
				&c.Site,
				&c.ThreadId,
				&c.ResponseTo,
				&c.Left,
				&c.Right,
				&c.Ctime,
				&c.Mtime,
				&c.Text,
				&c.Website,
				&c.Author,
				&c.EMail,
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

func CreateComment(c *Comment) (err error) {
	now := time.Now()
	c.Ctime = now
	c.Mtime = now

	err = tx(func(t *sql.Tx) {

		if c.ResponseTo != nil {
			row := t.QueryRow(`
				SELECT thread_id, right
				FROM comment
				WHERE id = ? AND site = ?`,
				c.ResponseTo,
				c.Site,
			)
			var threadId, right int64
			if err := row.Scan(&threadId, &right); err != nil {
				panic(err)
			}
			_, err := t.Exec(`
				UPDATE comment
				SET left = left + 2
				WHERE thread_id = ? AND left > ?;
				UPDATE comment
				SET right = right + 2
				WHERE thread_id = ? AND right >= ?;`,
				threadId,
				right,
				threadId,
				right,
			)
			if err != nil {
				panic(err)
			}
			c.Left = right
			c.Right = right + 1
			c.ThreadId = threadId
		} else {
			c.Left = 0
			c.Right = 1
		}

		res, err := t.Exec(`
			INSERT INTO comment
			(
				site, thread_id, response_to, left, right, ctime, mtime, text,
				website, author, email
			)
			VALUES
			(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`,
			&c.Site,
			&c.ThreadId,
			&c.ResponseTo,
			&c.Left,
			&c.Right,
			&c.Ctime,
			&c.Mtime,
			&c.Text,
			&c.Website,
			&c.Author,
			&c.EMail,
		)
		if err != nil {
			panic(err)
		}

		c.Id, err = res.LastInsertId()
		if err != nil {
			panic(err)
		}

		if c.ResponseTo == nil {
			c.ThreadId = c.Id

			_, err := t.Exec(`
				UPDATE comment
				SET thread_id = ?
				WHERE id = ?;`,
				c.ThreadId,
				c.Id,
			)
			if err != nil {
				panic(err)
			}
		}
	})

	return
}

func EditComment(c *Comment) (err error) {
	c.Mtime = time.Now()
	err = tx(func(t *sql.Tx) {
		_, err := t.Exec(`
			UPDATE comment
			SET
				mtime = ?, text = ?,
			WHERE id = ?;`,
			&c.Mtime,
			&c.Text,
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
