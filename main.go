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
	"strconv"
)

func main(){	
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
	chart.SetScaleX(5, &math32.Color{0.8, 0.8, 0.8})
	chart.SetFontSizeX(13)
	chart.SetRangeX(0.0, 7.0, 70.0)
	scene.Add(chart)

	createLine(scene, float32(height)-100, chart, one, "10*sin(x)*sin(x/10)")
	createLine(scene, float32(height)-80, chart, two, "sin(x)")
	createLine(scene, float32(height)-60, chart, three, "sqrt(x)")
	createLine(scene, float32(height)-40, chart, four, "sin(2.5*cos(x))")

	dd1 := gui.NewDropDown(100, gui.NewImageLabel("func"))
	dd1.SetPosition(480, float32(height)-100)
	dd1.Add(gui.NewImageLabel("10*sin(x)*sin(x/10)"))
	dd1.Add(gui.NewImageLabel("sin(x)"))
	dd1.Add(gui.NewImageLabel("sqrt(x)"))
	dd1.Add(gui.NewImageLabel("sin(2.5*cos(x))"))
	scene.Add(dd1)

	dd2 := gui.NewDropDown(100, gui.NewImageLabel("dots numb"))
	dd2.SetPosition(480, float32(height)-80)
	dd2.Add(gui.NewImageLabel("8"))
	dd2.Add(gui.NewImageLabel("20"))
	dd2.Add(gui.NewImageLabel("wrong y"))
	dd2.Add(gui.NewImageLabel("3"))
	scene.Add(dd2)

	ed1 := gui.NewEdit(100, "there can be your y")
	ed1.SetPosition(480, float32(height)-60)
	scene.Add(ed1)

	l1 := gui.NewLabel("")
	l1.SetPosition(590, float32(height)-100)
	scene.Add(l1)

	b1 := gui.NewButton("find y")
	b1.SetPosition(480, float32(height)-40)
	b1.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		var f func(float32) float32
		var dot int
		var wrong bool = false
		//fmt.Println(ed1.Text())
		switch dd1.Selected().Text(){
		case "10*sin(x)*sin(x/10)":
			f = one
		case "sin(x)":
			f = two
		case "sqrt(x)":
			f = three
		case "sin(x)sin(2.5*cos(x))":
			f = four
		}
		switch dd2.Selected().Text(){
		case "8":
			dot = 8
		case "20":
			dot = 20
		case "wrong y":
			dot = 20
			wrong = true
		case "3":
			dot = 3
		}
		xValues, yValues := getDots(math.Pi*10, dot, f)
		if wrong {
			yValues[10] = 14.2
		}
		x, err := strconv.ParseFloat(ed1.Text(),32)
		if err != nil {
			fmt.Println("oooooops") 
		}
		lagrange := InterpolateLagrangePolynomial(float32(x),math.Pi*10, dot, f, xValues, yValues)
		l1.SetText(fmt.Sprintf("%f", lagrange))
	})
	scene.Add(b1)

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

	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(scene, cam)
	})
}

func InterpolateLagrangePolynomial(x,max float32, steps int, f func(float32) float32, xValues,yValues []float32) float32{
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

func makeGraph(scene *core.Node, x,y float32, data []float32, color *math32.Color, chart *gui.Chart, xValues, yValues []float32, label string) {
	var dots []*gui.Image = make([]*gui.Image,0)
	if len(xValues) != 0 {
		for i := 0; i < len(xValues); i++ {
			dot, err := gui.NewImage("./dot.png")
			if err != nil {
				fmt.Println("ooops")
			}
			dot.SetPosition(49+xValues[i]*20.857,252-yValues[i]*21)
			dots = append(dots, dot)
		}
	}
	var g3 *gui.Graph
	cbG3 := gui.NewCheckBox(label)
	cbG3.SetPosition(x, y)
	cbG3.Subscribe(gui.OnChange, func(name string, ev interface{}) {
		if cbG3.Value() {
			g3 = chart.AddLineGraph(color, data)
			for i := 0; i < len(dots); i++ {
				scene.Add(dots[i])
			}
		} else {
			chart.RemoveGraph(g3)
			g3 = nil
			for i := 0; i < len(dots); i++ {
				scene.Remove(dots[i])
			}
		}
	})
	scene.Add(cbG3)
}

func createLine(scene *core.Node, height float32, chart *gui.Chart, f func(float32) float32, label string){
	data1 := make([]float32, 0)
	for x := 0.0; x < math.Pi*10; x += 0.1 {
		data1 = append(data1, f(float32(x)))
	}
	makeGraph(scene, 10, height, data1, &math32.Color{0, 0, 1}, chart, []float32{}, []float32{}, label)

	dlag := make([]float32, 0)
	xValues, yValues := getDots(math.Pi*10, 8, f)
	for i := 0.0; i < math.Pi*10; i+=0.1 {
		lagrange := InterpolateLagrangePolynomial(float32(i),math.Pi*10, 8, f, xValues, yValues)
		dlag = append(dlag, lagrange)
	}
	makeGraph(scene, 150, height, dlag, &math32.Color{1, 0, 0}, chart, xValues, yValues,"8 dots")

	xValues1, yValues1 := getDots(math.Pi*10, 20, f)
	dlagt := make([]float32, 0)
	for i := 0.0; i < math.Pi*10; i+=0.1 {
		lagrange := InterpolateLagrangePolynomial(float32(i),math.Pi*10, 20, f, xValues1, yValues1)
		dlagt = append(dlagt, lagrange)
	}
	makeGraph(scene, 220, height, dlagt, &math32.Color{0, 1, 0}, chart, xValues1, yValues1, "20 dots")

	xValues2, yValues2 := getDots(math.Pi*10, 20, f)
	yValues2[10] = 14.2
	dlagt2 := make([]float32, 0)
	for i := 0.0; i < math.Pi*10; i+=0.1 {
		lagrange := InterpolateLagrangePolynomial(float32(i),math.Pi*10, 20, f, xValues2, yValues2)
		dlagt2 = append(dlagt2, lagrange)
	}
	makeGraph(scene, 300, height, dlagt2, &math32.Color{0, 1, 1}, chart, xValues2, yValues2, "wrong y")

	xValues3, yValues3 := getDots(math.Pi*10, 3, f)
	dlagt3 := make([]float32, 0)
	for i := 0.0; i < math.Pi*10; i+=0.1 {
		lagrange := InterpolateLagrangePolynomial(float32(i),math.Pi*10, 3, f, xValues3, yValues3)
		dlagt3 = append(dlagt3, lagrange)
	}
	makeGraph(scene, 380, height, dlagt3, &math32.Color{0, 0, 0}, chart, xValues3, yValues3, "only 3 dots")
}