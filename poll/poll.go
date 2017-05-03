package poll

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/ruprict/poller/config"
)

func Poll(c chan<- int, conf *config.Config, db *sql.DB) {
	tick := time.Tick(5 * time.Second)
	for _ = range tick {
		fmt.Println("*** Checking for orders: ", time.Now())
		rows, err := db.Query("SELECT * FROM bopis_orders where order_acknowledged='0' and shipnode=?;", *conf.Shipnode)
		if err != nil {
			fmt.Println("*** query error")
			log.Fatalln(err)
		}
		for rows.Next() {
			var id int
			var brand, shipnode, order_id, customer string
			var order_date time.Time
			var order_ack bool
			err = rows.Scan(&id, &brand, &shipnode, &order_id, &customer, &order_date, &order_ack)
			if err != nil {
				fmt.Println("Error scanning row: ", err)
			}
			c <- id
		}
		fmt.Println("*** Done")
		defer rows.Close()

	}
}
