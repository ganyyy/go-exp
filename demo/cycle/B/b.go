package cycle_b

import (
	"fmt"

	cycal_a "ganyyy.com/go-exp/demo/cycle/A"
)

var B = 1

func Print() {
	fmt.Println("B", B, "A", cycal_a.A)
}
