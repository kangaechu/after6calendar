package main

import (
	"fmt"
	"github.com/kangaechu/after6calendar"
	"time"
)

func main() {

	//after6calendar.GetEventsJson()
	program := after6calendar.GetProgramSummary(time.Date(2019, 06, 14,
		18, 00, 00, 0, time.Local))
	fmt.Print(*program)
}
