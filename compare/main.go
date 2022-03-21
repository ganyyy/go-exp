package main

func Check1(v int) int {
	if v == 10 {
		return 100
	}
	if v == 200 {
		return 10
	} else {
		return 0
	}
}

func Check2(v int) int {
	if v == 10 {
		return 100
	} else if v == 200 {
		return 10
	} else {
		return 0
	}
}
