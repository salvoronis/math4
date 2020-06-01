package main

func main(){
	createWindow()
}

func InterpolateLagrangePolynomial(x float32, steps int, f func(float32) float32, xValues,yValues []float32) float32{
	var lagrange_pol, basics_pol float32

	for i := 0; i < steps; i++ {
		basics_pol = 1
		for j := 0; j < steps; j++ {
			if j==i {
				continue
			}
			basics_pol *= (x - xValues[j])/(xValues[i] - xValues[j])
		}
		lagrange_pol += basics_pol*yValues[i]
	}
	return lagrange_pol
}

func getDots(max float32, steps int, f func(float32) float32) ([]float32, []float32){
	var step float32 = max/float32(steps)
	var xValues, yValues []float32 = make([]float32, steps), make([]float32, steps)
	var count int = 0
	for i := float32(0.0); i < max; i+=step {
		xValues[count] = i
		yValues[count] = f(i)
		count++
	}
	return xValues, yValues
}