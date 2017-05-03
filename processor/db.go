package processor

import "database/sql"

type DbProcessor struct{}

func (p DbProcessor) Process(str int, db *sql.DB) {
	db.Exec("update bopis_orders set order_acknowledged='1' where id=?;", str)
}
