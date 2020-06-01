package main

import(
	"github.com/g3n/engine/math32"
)

func one(x float32) float32{
	return 10*math32.Sin(x)*math32.Sin(x/10)
}

func two(x float32) float32{
	return math32.Sin(x)
}

func three(x float32) float32{
	return math32.Sqrt(x)
}

func four(x float32) float32{
	return math32.Sin(2.5*math32.Cos(x))
}