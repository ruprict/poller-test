package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "bitbucket.org/phiggins/db2cli"
)

var counter = 0
var port *int

func init() {
	fmt.Println("Firing it up ***")
	port = flag.Int("port", 9090, "port for web server")
	connectToDb()

}

func Poll(c chan<- int, db *sql.DB) {
	tick := time.Tick(5 * time.Second)
	for _ = range tick {
		fmt.Println("*** Checking for orders: ", time.Now())
		rows, err := db.Query("SELECT * FROM bopis_orders where order_acknowledged='0';")
		if err != nil {
			fmt.Println("*** query error")
			log.Fatalln(err)
		}
		for rows.Next() {
			var id int
			var brand, store_id, order_id, customer string
			var order_date time.Time
			var order_ack bool
			err = rows.Scan(&id, &brand, &store_id, &order_id, &customer, &order_date, &order_ack)
			if err != nil {
				fmt.Println("Error scanning row: ", err)
			}
			c <- id
		}
		fmt.Println("*** Done")
		defer rows.Close()

	}
}

func connectToDb() *sql.DB {
	fmt.Println("*** connnecting to effing db ***")
	db, err := sql.Open("db2-cli", "DATABASE=geg_test; HOSTNAME=localhost; PORT=50000; PROTOCOL=TCPIP; UID=db2inst1; PWD=skookumpassword;")
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func startWebServer() {

	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Running on port ", *port)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "So far I've done this ", counter, " times.")
	})

	http.ListenAndServe(":"+strconv.Itoa(*port), nil)

}

func Processor(c <-chan int, db *sql.DB) {
	for str := range c {
		fmt.Println(str)
		// TODO: TURN ON LIGHT
		// Mark row as acknowledged
		db.Exec("update bopis_orders set order_acknowledged='1' where id=?;", str)
		counter++
	}
}

func main() {
	fmt.Println("** In main")
	db := connectToDb()
	defer db.Close()
	c := make(chan int)
	go Poll(c, db)
	go Processor(c, db)
	startWebServer()
}
