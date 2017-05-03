package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/ruprict/poller/config"
	"github.com/ruprict/poller/poll"
	"github.com/ruprict/poller/processor"
	"github.com/ruprict/poller/web"
	rpio "github.com/stianeikeland/go-rpio"
)

var connString string
var pin = rpio.Pin(17) // corresponsds to physical pin 17
var conf *config.Config

func init() {
	fmt.Println("Firing it up ***")
	port := flag.Int("port", 9090, "port for web server")
	thisStore := flag.String("shipnode", "", "Shipnode value for this store")
	flag.Parse()
	if *thisStore == "" {
		log.Fatalln("you must supply a shipnode for this store")
	}
	connString = os.Getenv("DB_CONNECTION_STRING")
	if connString == "" {
		log.Fatalln("DB_CONNECTION_STRING must be set in the environment")
	}
	conf = config.New(port, thisStore, connString)

}

func connectToDb() *sql.DB {
	fmt.Println("*** connnecting to effing db ***")
	fmt.Println(connString)
	db, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Connected!")
	return db
}

func main() {
	fmt.Println("** In main")
	if err := rpio.Open(); err != nil {
		log.Println("Can't open rpio")
	} else {
		defer rpio.Close()

		pin.Output()
		//REset the pin
		pin.Low()

	}

	db := connectToDb()
	defer db.Close()
	c := make(chan int)
	go poll.Poll(c, conf, db)
	go processor.Start(c, db)
	web.Start(conf, db)

}
