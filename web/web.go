package web

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/ruprict/poller/config"
	rpio "github.com/stianeikeland/go-rpio"
)

var templates = template.Must(template.ParseGlob("web/html/*"))

func Start(config *config.Config, db *sql.DB) {

	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Running on port ", *config.Port)
		fmt.Fprintln(w, "Shipnode is ", *config.Shipnode)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			fmt.Println("** Executing tempalte")
			err := templates.ExecuteTemplate(w, "index.tmpl", nil)
			if err != nil {
				fmt.Println("** Template error: ", err)
			}
			return
		}

		fmt.Println("*** Going low")
		var pin = rpio.Pin(17) // corresponsds to physical pin 17
		pin.Low()

		http.Redirect(w, r, "/", 302)
	})

	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			err := templates.ExecuteTemplate(w, "new_order.tmpl", nil)
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

	http.ListenAndServe(":"+strconv.Itoa(*config.Port), nil)
}
