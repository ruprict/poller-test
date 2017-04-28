package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "bitbucket.org/phiggins/db2cli"
)

var counter = 0
var port *int
var connString string

func init() {
	fmt.Println("Firing it up ***")
	port = flag.Int("port", 9090, "port for web server")
	connString = os.Getenv("DB2_CONNECTION_STRING")
	if connString == "" {
		log.Fatalln("DB2_CONNECTION_STRING must be set in the environment")
	}

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
	fmt.Println(connString)
	db, err := sql.Open("db2-cli", connString)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Connected!")
	return db
}

func startWebServer(db *sql.DB) {

	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Running on port ", *port)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "So far I've done this ", counter, " times.")
	})

	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			fmt.Println("** Executing tempalte")
			t, err := template.ParseFiles("./html/new_order.tmpl")
			if err != nil {
				fmt.Println("** Template error: ", err)
			}
			err = t.Execute(w, nil)
			if err != nil {
				fmt.Println("** Template error: ", err)
			}
			return
		}
		r.ParseForm()
		brand := r.FormValue("brand")
		store_id := r.FormValue("store_id")
		order_id := r.FormValue("order_id")
		customer := r.FormValue("customer")
		date := time.Now()

		result, err := db.Exec("insert into bopis_orders (brand, store_id, order_id, customer_name, order_date) values (?,?,?,?,?)", brand, store_id, order_id, customer, date)
		if err != nil {
			fmt.Println("**** insert error: ", err)
		}
		fmt.Println(result)
		http.Redirect(w, r, "/", 302)
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
	startWebServer(db)
}
