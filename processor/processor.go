package processor

import (
	"database/sql"
	"fmt"
)

type Processor interface {
	Process(str int, db *sql.DB)
}

var processors []Processor

func init() {

	processors = make([]Processor, 2)
	processors[0] = LampProcessor{}
	processors[1] = DbProcessor{}
}

func Start(c <-chan int, db *sql.DB) {
	for id := range c {
		fmt.Println(id)
		handle(id, db)
	}
}

func handle(id int, db *sql.DB) {
	for _, p := range processors {
		p.Process(id, db)
	}
}
