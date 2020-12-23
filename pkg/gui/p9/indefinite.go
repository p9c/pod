package p9

import (
	"image"
	"image/color"
	"math"
	"time"

	"gioui.org/f32"
	l "gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type Indefinite struct {
	th    *Theme
	color color.NRGBA
	scale float32
}

// Indefinite creates an indefinite loading animation icon
func (th *Theme) Indefinite() *Indefinite {
	return &Indefinite{
		th:    th,
		color: th.Colors.Get("Primary"),
	}
}

// Scale sets the size of the spinner
func (lo *Indefinite) Scale(scale float32) *Indefinite {
	lo.scale = scale
	return lo
}

// Color sets the color of the spinner
func (lo *Indefinite) Color(color string) *Indefinite {
	lo.color = lo.th.Colors.Get(color)
	return lo
}

// Fn renders the loader
func (lo *Indefinite) Fn(gtx l.Context) l.Dimensions {
	diam := gtx.Constraints.Min.X
	if minY := gtx.Constraints.Min.Y; minY > diam {
		diam = minY
	}
	if diam == 0 {
		diam = gtx.Px(lo.th.TextSize.Scale(lo.scale))
	}
	sz := gtx.Constraints.Constrain(image.Pt(diam, diam))
	radius := float64(sz.X) * .5
	defer op.Push(gtx.Ops).Pop()
	op.Offset(f32.Pt(float32(radius), float32(radius))).Add(gtx.Ops)

	dt := (time.Duration(gtx.Now.UnixNano()) % (time.Second)).Seconds()
	startAngle := dt * math.Pi * 2
	endAngle := startAngle + math.Pi*1.5

	clipLoader(gtx.Ops, startAngle, endAngle, radius)
	paint.ColorOp{
		Color: lo.color,
	}.Add(gtx.Ops)
	op.Offset(f32.Pt(-float32(radius), -float32(radius))).Add(gtx.Ops)
	paint.PaintOp{
		// Rect: f32.Rectangle{Max: l.FPt(sz)},
	}.Add(gtx.Ops)
	op.InvalidateOp{}.Add(gtx.Ops)
	return l.Dimensions{
		Size: sz,
	}
}

func clipLoader(ops *op.Ops, startAngle, endAngle, radius float64) {
	const thickness = .25
	
	var (
		width = float32(radius * thickness)
		delta = float32(endAngle - startAngle)
		
		vy, vx = math.Sincos(startAngle)
		
		pen    = f32.Pt(float32(vx), float32(vy)).Mul(float32(radius))
		center = f32.Pt(0, 0).Sub(pen)
		
		style = clip.StrokeStyle{
			Cap: clip.FlatCap,
		}
		
		p clip.Path
	)
	
	p.Begin(ops)
	p.Move(pen)
	p.Arc(center, center, delta)
	p.Stroke(width, style).Add(ops)
}
