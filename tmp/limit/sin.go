//go:build ignore
// +build ignore

package main

import (
	"math"
)

func main() {
	CalcAngle()
}

const Parse = math.Pi / 180

func AngleToRadian(angle float64) float64 {
	return angle * Parse
}

func RadianToAngle(radian float64) float64 {
	return radian / Parse
}

func CalcAngle() {
	var length float64 = 1
	var side1 = math.Tan(AngleToRadian(22)) * length
	var side2 = length / math.Tan(AngleToRadian(67))
	var a = math.Atan((length - side2) / (length - side1))
	println(180 - RadianToAngle(a) - (90 - 22)) // 68
}
