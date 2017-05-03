package processor

import (
	"database/sql"
	"fmt"

	rpio "github.com/stianeikeland/go-rpio"
)

var pin = rpio.Pin(17) // corresponsds to physical pin 17

type LampProcessor struct{}

func (l LampProcessor) Process(_ int, _ *sql.DB) {
	fmt.Println("** Toggling pin")
	pin.High()
}
