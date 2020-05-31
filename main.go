package main

import(
	"fmt"
	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/window"
	"time"
	"math"
)

func main(){
	dlag := make([]float32, 0)
	for i := 0.0; i < math.Pi*10; i+=0.1 {
		lagrange := InterpolateLagrangePolynomial(float32(i), 32, one)
		dlag = append(dlag, lagrange)
		//fmt.Println(lagrange)	
	}
	
	a := app.App()
	scene := core.NewNode()

	gui.Manager().Set(scene)

	chart := gui.NewChart(0, 0)
	chart.SetMargins(10, 10, 10, 10)
	chart.SetBorders(2, 2, 2, 2)
	chart.SetBordersColor(math32.NewColor("green"))
	chart.SetColor(math32.NewColor("white"))
	chart.SetTitle("Chart Title", 16)
	chart.SetPosition(0, 0)
	width, height := a.GetSize()
	chart.SetSize(float32(width), float32(height)-100)
	chart.SetScaleY(5, &math32.Color{0.8, 0.8, 0.8})
	chart.SetFontSizeY(13)
	//chart.SetRangeY(-10.0, 10.0)
	chart.SetScaleX(5, &math32.Color{0.8, 0.8, 0.8})
	chart.SetFontSizeX(13)
	chart.SetRangeX(0.0, 7.0, 70.0)
	scene.Add(chart)

	var g1 *gui.Graph
	data1 := make([]float32, 0)
	var x float32
	for x = 0; x < math.Pi*10; x += 0.1 {
		data1 = append(data1, 10*math32.Sin(x)*math32.Sin(x/10))
	}
	//chart.AddLineGraph(&math32.Color{0, 0, 1}, data1)
	cbG1 := gui.NewCheckBox("Graph1")
	cbG1.SetPosition(10, float32(height)-100)
	cbG1.Subscribe(gui.OnChange, func(name string, ev interface{}) {
		if cbG1.Value() {
			g1 = chart.AddLineGraph(&math32.Color{0, 0, 1}, data1)
		} else {
			chart.RemoveGraph(g1)
			g1 = nil
		}
	})
	scene.Add(cbG1)
	cbG1.SetValue(true)

	var g2 *gui.Graph
	/*data2 := make([]float32, 0)
	var x2 float32
	for x2 = 0; x2 < 2*math.Pi*10; x2 += 0.1 {
		data2 = append(data2, -2+5*math32.Cos(x2/3))
	}*/
	cbG2 := gui.NewCheckBox("Graph2")
	cbG2.SetPosition(90, float32(height)-100)
	cbG2.Subscribe(gui.OnChange, func(name string, ev interface{}) {
		if cbG2.Value() {
			g2 = chart.AddLineGraph(&math32.Color{1, 0, 0}, dlag)
		} else {
			chart.RemoveGraph(g2)
			g2 = nil
		}
	})
	scene.Add(cbG2)

	cam := camera.New(1)
	cam.SetPosition(0, 0, 3)
	scene.Add(cam)

	onResize := func(evname string, ev interface{}) {
		width, height := a.GetSize()
		a.Gls().Viewport(0, 0, int32(width), int32(height))
		cam.SetAspect(float32(width) / float32(height))
	}
	a.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	a.Gls().ClearColor(0.2, 0.7, 0.9, 10.0)

	// Run the application
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(scene, cam)
	})
}

func InterpolateLagrangePolynomial(x float32, size int, f func(float32) float32) float32{
	var xValues, yValues []float32 = make([]float32, size), make([]float32, size)
	var lagrange_pol, basics_pol float32
	for i := 0; i < size; i++ {
		xValues[i] = float32(i)
		yValues[i] = f(float32(i))
		fmt.Println(xValues[i])
	}

	for i := 0; i < size; i++ {
		basics_pol = 1
		for j := 0; j < size; j++ {
			if j==i {
				continue
			}
			basics_pol *= (x - xValues[j])/(xValues[i] - xValues[j])
		}
		lagrange_pol += basics_pol*yValues[i]
	}
	return lagrange_pol
}

func one(x float32) float32{
	return 10*math32.Sin(x)*math32.Sin(x/10)
}


