package model

import (
	"database/sql"
	"fmt"
)

func tx(txFunc func(*sql.Tx)) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}

	defer func() {
		if errT := recover(); errT != nil {
			var ok bool
			err, ok = errT.(error)
			if !ok {
				err = fmt.Errorf("%s", errT)
			}
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	txFunc(tx)
	return nil
}
