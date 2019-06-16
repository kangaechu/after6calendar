package main

import (
	"fmt"
	"github.com/kangaechu/after6renamer"
	"time"
)

func main() {

	//after6renamer.GetEventsJson()
	program := after6renamer.GetProgramSummary(time.Date(2019, 06, 14,
		18, 00, 00, 0, time.Local))
	fmt.Print(*program)
}
