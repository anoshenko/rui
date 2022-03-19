package main

import (
	"math"
	"strconv"

	"github.com/anoshenko/rui"
)

const canvasDemoText = `
GridLayout {
	width = 100%, height = 100%, cell-height = "auto, 1fr",
	content = [
		DropDownList {
			id = canvasType, current = 0, margin = 8px,
			items = ["Image", "Rectangles & ellipse", "Text style", "Text align", "Line style", "Transformations"]
		},
		CanvasView {
			id = canvas, row = 1, width = 100%, height = 100%,
		}
	]
}
`

var sampleImage rui.Image

func rectangleCanvasDemo(canvas rui.Canvas) {
	width := canvas.Width()
	height := canvas.Height()

	canvas.Save()

	canvas.SetSolidColorFillStyle(0xFF008000)
	canvas.SetSolidColorStrokeStyle(0xFFFF0000)
	w2 := width / 2
	h2 := height / 2
	canvas.FillRect(10, 10, w2-20, h2-20)
	canvas.StrokeRect(9.5, 9.5, w2-19, h2-19)

	canvas.SetLinearGradientFillStyle(w2+10, 10, 0xFFFF0000, width-20, 10, 0xFF0000FF, []rui.GradientPoint{
		{Offset: 0.3, Color: 0xFFFFFF00},
		{Offset: 0.5, Color: 0xFF00FF00},
		{Offset: 0.7, Color: 0xFF00FFFF},
	})
	canvas.SetLinearGradientStrokeStyle(10, 10, 0xFFFF00FF, 10, h2-20, 0xFF00FFFF, []rui.GradientPoint{
		{Offset: 0.5, Color: 0xFF00FF00},
	})
	canvas.SetLineWidth(5)
	canvas.FillAndStrokeRoundedRect(w2+7.5, 7.5, w2-15, h2-15, 20)

	canvas.SetRadialGradientFillStyle(w2/2-20, h2+h2/2-20, 10, 0xFFFF0000, w2/2+20, h2+h2/2+20, w2/2, 0xFF0000FF, []rui.GradientPoint{
		{Offset: 0.3, Color: 0xFFFFFF00},
		{Offset: 0.5, Color: 0xFF00FF00},
		{Offset: 0.7, Color: 0xFF00FFFF},
	})
	canvas.SetRadialGradientStrokeStyle(w2/2, h2+h2/2, h2/2, 0xFFFFFF00, w2/2, h2+h2/2, h2, 0xFF00FFFF, []rui.GradientPoint{
		{Offset: 0.5, Color: 0xFF00FF00},
	})
	canvas.SetLineWidth(7)
	canvas.FillAndStrokeRect(10, h2+10, w2-20, h2-20)

	//canvas.SetSolidColorFillStyle(0xFF00FFFF)
	canvas.SetImageFillStyle(sampleImage, rui.RepeatXY)
	canvas.SetSolidColorStrokeStyle(0xFF0000FF)
	canvas.SetLineWidth(4)
	canvas.FillAndStrokeEllipse(w2+w2/2, h2+h2/2, w2/2-10, h2/2-10, 0)

	canvas.Restore()
}

func textCanvasDemo(canvas rui.Canvas) {

	canvas.Save()
	canvas.SetTextAlign(rui.LeftAlign)
	canvas.SetTextBaseline(rui.TopBaseline)

	if canvas.View().Session().DarkTheme() {
		canvas.SetSolidColorFillStyle(0xFFFFFFFF)
		canvas.SetSolidColorStrokeStyle(0xFFFFFFFF)
	} else {
		canvas.SetSolidColorFillStyle(0xFF000000)
		canvas.SetSolidColorStrokeStyle(0xFF000000)
	}
	canvas.FillText(10, 10, "Default font")
	canvas.StrokeText(300, 10, "Default font")

	canvas.SetSolidColorFillStyle(0xFF800000)
	canvas.SetSolidColorStrokeStyle(0xFF800080)
	canvas.SetFont("courier", rui.Pt(12))
	canvas.FillText(10, 30, "courier, 12pt")
	canvas.StrokeText(300, 30, "courier, 12pt")

	canvas.SetSolidColorFillStyle(0xFF008000)
	canvas.SetSolidColorStrokeStyle(0xFF008080)
	canvas.SetFontWithParams("Courier new, courier", rui.Pt(12), rui.FontParams{
		Italic: true,
	})
	canvas.FillText(10, 50, `Courier new, 12pt, italic`)
	canvas.StrokeText(300, 50, `Courier new, 12pt, italic`)

	canvas.SetSolidColorFillStyle(0xFF000080)
	canvas.SetLinearGradientStrokeStyle(10, 70, 0xFF00FF00, 10, 90, 0xFFFF00FF, nil)
	canvas.SetFontWithParams("sans-serif", rui.Pt(12), rui.FontParams{
		SmallCaps: true,
	})
	canvas.FillText(10, 70, "sans-serif, 12pt, small-caps")
	canvas.StrokeText(300, 70, "sans-serif, 12pt, small-caps")

	canvas.SetLinearGradientFillStyle(10, 90, 0xFFFF0000, 10, 110, 0xFF0000FF, nil)
	canvas.SetSolidColorStrokeStyle(0xFF800080)
	canvas.SetFontWithParams("serif", rui.Pt(12), rui.FontParams{
		Weight: 7,
	})
	canvas.FillText(10, 90, "serif, 12pt, weight: 7(bold)")
	canvas.StrokeText(300, 90, "serif, 12pt, weight: 7(bold)")

	widthSample := "Text width sample"
	w := canvas.TextWidth(widthSample, "sans-serif", rui.Px(20))
	canvas.SetFont("sans-serif", rui.Px(20))
	canvas.SetSolidColorFillStyle(rui.Blue)
	canvas.SetTextBaseline(rui.BottomBaseline)
	canvas.FillText(10, 150, widthSample)

	canvas.SetSolidColorStrokeStyle(rui.Black)
	canvas.SetLineWidth(1)
	canvas.DrawLine(10, 150, 10, 170)
	canvas.DrawLine(10+w, 150, 10+w, 170)
	canvas.DrawLine(10, 168, 10+w, 168)
	canvas.DrawLine(10, 168, 20, 165)
	canvas.DrawLine(10, 168, 20, 171)
	canvas.DrawLine(10+w, 168, w, 165)
	canvas.DrawLine(10+w, 168, w, 171)

	canvas.SetSolidColorFillStyle(rui.Black)
	canvas.SetFont("sans-serif", rui.Px(8))
	canvas.SetTextAlign(rui.CenterAlign)
	canvas.FillText(10+w/2, 167, strconv.FormatFloat(w, 'g', -1, 64))

	canvas.Restore()
}

func textAlignCanvasDemo(canvas rui.Canvas) {
	canvas.Save()
	canvas.SetFont("sans-serif", rui.Pt(10))
	if canvas.View().Session().DarkTheme() {
		canvas.SetSolidColorFillStyle(0xFFFFFFFF)
	} else {
		canvas.SetSolidColorFillStyle(0xFF000000)
	}
	canvas.SetSolidColorStrokeStyle(0xFF00FFFF)

	baseline := []string{"Alphabetic", "Top", "Middle", "Bottom", "Hanging", "Ideographic"}
	align := []string{"Left", "Right", "Center", "Start", "End"}
	center := []float64{20, 120, 70, 20, 120}
	for b, bText := range baseline {
		for a, aText := range align {
			canvas.SetTextAlign(a)
			canvas.SetTextBaseline(b)
			x := float64(a * 140)
			y := float64(b * 40)

			canvas.DrawLine(x+4, y+20, x+132, y+20)
			canvas.DrawLine(x+center[a], y+2, x+center[a], y+38)
			canvas.FillText(x+center[a], y+20, bText+","+aText)
		}
	}
	canvas.Restore()
}

func lineStyleCanvasDemo(canvas rui.Canvas) {
	canvas.Save()

	canvas.SetSolidColorStrokeStyle(0xFF00FFFF)
	canvas.SetLineWidth(1)
	canvas.DrawLine(20, 30, 20, 90)
	canvas.DrawLine(170, 30, 170, 90)

	canvas.SetSolidColorStrokeStyle(0xFF0000FF)
	canvas.SetFont("courier", rui.Pt(12))
	canvas.SetTextBaseline(rui.MiddleBaseline)
	canvas.FillText(80, 15, "SetLineCap(...)")

	canvas.SetFont("courier", rui.Pt(10))
	for i, cap := range []string{"ButtCap", "RoundCap", "SquareCap"} {
		canvas.SetSolidColorStrokeStyle(0xFF00FFFF)
		canvas.SetLineWidth(1)
		y := float64(40 + 20*i)
		canvas.DrawLine(10, y, 180, y)
		if canvas.View().Session().DarkTheme() {
			canvas.SetSolidColorStrokeStyle(0xFFFFFFFF)
		} else {
			canvas.SetSolidColorStrokeStyle(0xFF000000)
		}

		canvas.SetLineWidth(10)
		canvas.SetLineCap(i)
		canvas.DrawLine(20, y, 170, y)
		canvas.FillText(200, y, cap)
	}

	if canvas.View().Session().DarkTheme() {
		canvas.SetSolidColorStrokeStyle(0xFFFFFFFF)
		canvas.SetSolidColorFillStyle(0xFF00FFFF)
	} else {
		canvas.SetSolidColorStrokeStyle(0xFF000000)
		canvas.SetSolidColorFillStyle(0xFF0000FF)
	}
	canvas.SetFont("courier", rui.Pt(12))
	canvas.FillText(80, 115, "SetLineJoin(...)")

	canvas.SetLineWidth(10)
	canvas.SetLineCap(rui.ButtCap)

	canvas.SetFont("courier", rui.Pt(10))
	for i, join := range []string{"MiterJoin", "RoundJoin", "BevelJoin"} {
		y := float64(140 + 40*i)
		path := rui.NewPath()
		path.MoveTo(20, y)
		path.LineTo(50, y+40)
		path.LineTo(80, y)
		path.LineTo(110, y+40)
		path.LineTo(140, y)
		path.LineTo(170, y+40)
		path.LineTo(200, y)
		canvas.SetLineJoin(i)
		canvas.StrokePath(path)
		canvas.FillText(210, y+20, join)
	}
	canvas.SetFont("courier", rui.Pt(12))
	canvas.FillText(20, 300, "SetLineDash([]float64{16, 8, 4, 8}, ...)")

	canvas.SetFont("courier", rui.Pt(10))
	canvas.SetLineDash([]float64{16, 8, 4, 8}, 0)
	canvas.SetLineWidth(4)

	canvas.SetLineCap(rui.ButtCap)
	canvas.DrawLine(20, 330, 200, 330)
	canvas.FillText(220, 330, "SetLineCap(ButtCap)")

	canvas.SetLineDash([]float64{16, 8, 4, 8}, 4)
	canvas.DrawLine(20, 360, 200, 360)
	canvas.FillText(220, 360, "offset = 4")

	canvas.SetLineDash([]float64{16, 8, 4, 8}, 0)
	canvas.SetLineCap(rui.RoundCap)
	canvas.SetShadow(4, 4, 2, 0xFF808080)
	canvas.DrawLine(20, 390, 200, 390)
	canvas.ResetShadow()
	canvas.FillText(220, 390, "SetLineCap(RoundCap)")

	canvas.Restore()
}

func transformCanvasDemo(canvas rui.Canvas) {
	drawFigure := func() {
		w := int(canvas.Width() / 2)
		h := int(canvas.Height() / 2)
		nx := (w/2)/20 + 1
		ny := (h/2)/20 + 1
		x0 := float64((w - nx*20) / 2)
		y0 := float64((h - ny*20) / 2)
		x1 := x0 + float64((nx-1)*20)
		y1 := y0 + float64((ny-1)*20)

		canvas.SetFont("serif", rui.Pt(10))
		if canvas.View().Session().DarkTheme() {
			canvas.SetSolidColorStrokeStyle(0xFFFFFFFF)
			canvas.SetSolidColorFillStyle(0xFFFFFFFF)
		} else {
			canvas.SetSolidColorStrokeStyle(0xFF000000)
			canvas.SetSolidColorFillStyle(0xFF000000)
		}
		canvas.SetTextAlign(rui.CenterAlign)
		canvas.SetTextBaseline(rui.BottomBaseline)
		for i := 0; i < nx; i++ {
			x := x0 + float64(i*20)
			canvas.DrawLine(x, y0, x, y1)
			canvas.FillText(x, y0-4, strconv.Itoa(i))
		}

		canvas.SetTextAlign(rui.RightAlign)
		canvas.SetTextBaseline(rui.MiddleBaseline)
		for i := 0; i < ny; i++ {
			y := y0 + float64(i*20)
			canvas.DrawLine(x0, y, x1, y)
			canvas.FillText(x0-4, y, strconv.Itoa(i))
		}
	}

	canvas.SetFont("courier", rui.Pt(14))
	if canvas.View().Session().DarkTheme() {
		canvas.SetSolidColorFillStyle(0xFFFFFFFF)
	} else {
		canvas.SetSolidColorFillStyle(0xFF000000)
	}
	canvas.SetTextAlign(rui.CenterAlign)
	canvas.SetTextBaseline(rui.TopBaseline)

	canvas.FillText(canvas.Width()/4, 8, "Original")
	canvas.FillText(canvas.Width()*3/4, 8, "SetScale(1.2, 0.8)")
	canvas.FillText(canvas.Width()/4, canvas.Height()/2+8, "SetRotation(math.Pi / 6)")
	canvas.FillText(canvas.Width()*3/4, canvas.Height()/2+8, "SetTransformation(0.8, 1.2, 0.2, 0.4, ...)")

	drawFigure()

	canvas.SetScale(1.2, 0.8)
	canvas.SetTranslation(canvas.Width()/2.4, 0)
	drawFigure()

	canvas.ResetTransformation()
	canvas.SetTranslation(canvas.Width()/8, canvas.Height()/2-canvas.Height()/8)
	canvas.SetRotation(math.Pi / 6)
	drawFigure()

	canvas.ResetTransformation()
	//canvas.SetTranslation(canvas.Width()/2, canvas.Height()/2)
	canvas.SetTransformation(0.8, 1.2, 0.2, 0.4, canvas.Width()/(2*0.8)-canvas.Width()/8, canvas.Height()/(2*1.2)-canvas.Height()/16)
	drawFigure()
}

var image rui.Image

func imageCanvasDemo(canvas rui.Canvas) {
	if image != nil {
		canvas.DrawImage(50, 20, image)
	} else {
		image = rui.LoadImage("tile00.svg", func(img rui.Image) {
			if img.LoadingStatus() == rui.ImageReady {
				canvas.View().Redraw()
			}
		}, canvas.View().Session())
	}
}

func createCanvasDemo(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, canvasDemoText)
	if view == nil {
		return nil
	}

	rui.Set(view, "canvas", rui.DrawFunction, imageCanvasDemo)

	rui.Set(view, "canvasType", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		drawFuncs := []func(rui.Canvas){
			imageCanvasDemo,
			rectangleCanvasDemo,
			textCanvasDemo,
			textAlignCanvasDemo,
			lineStyleCanvasDemo,
			transformCanvasDemo,
		}
		if number >= 0 && number < len(drawFuncs) {
			rui.Set(view, "canvas", rui.DrawFunction, drawFuncs[number])
		}
	})

	sampleImage = rui.LoadImage("image_sample.png", nil, session)

	return view
}
