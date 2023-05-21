package cycle_a

import (
	"fmt"

	cycal_b "ganyyy.com/go-exp/demo/cycle/B"
)

var A = 1

func Print() {
	fmt.Println("A", A, "B", cycal_b.B)
}
