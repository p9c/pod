package duoui

import (
	"gioui.org/layout"
	"gioui.org/unit"

	"image/color"
)

func DuoUInetwork(duo *DuoUI) layout.FlexChild {
	return duo.comp.content.l.Flex(duo.gc, 1, func() {
		duo.comp.network.i.Layout(duo.gc, func() {
			DuoUIdrawRect(duo.gc, duo.cs.Width.Max, duo.cs.Height.Max, color.RGBA{A: 0xff, R: 0x30, G: 0xcf, B: 0xcf}, 0, 0, 0, 0)
			// Overview <<<
			in := layout.UniformInset(unit.Dp(60))
			in.Layout(duo.gc, func() {
				drawRect(duo.gc, color.RGBA{A: 0xff, B: 0xff})

				duo.th.H5("network :").Layout(duo.gc)
			})
			// Overview >>>
		})
	})
}