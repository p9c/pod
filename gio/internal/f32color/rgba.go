// SPDX-License-Identifier: Unlicense OR MIT

package f32color

import (
	"image/color"
	"math"
)

// RGBA is a 32 bit floating point linear premultiplied color space.
type RGBA struct {
	R, G, B, A float32
}

// Array returns rgba values in a [4]float32 array.
func (rgba RGBA) Array() [4]float32 {
	return [4]float32{rgba.R, rgba.G, rgba.B, rgba.A}
}

// Float32 returns r, g, b, a values.
func (col RGBA) Float32() (r, g, b, a float32) {
	return col.R, col.G, col.B, col.A
}

// SRGBA converts from linear to sRGB color space.
func (col RGBA) SRGB() color.NRGBA {
	if col.A == 0 {
		return color.NRGBA{}
	}
	return color.NRGBA{
		R: uint8(linearTosRGB(col.R/col.A)*255 + .5),
		G: uint8(linearTosRGB(col.G/col.A)*255 + .5),
		B: uint8(linearTosRGB(col.B/col.A)*255 + .5),
		A: uint8(col.A*255 + .5),
	}
}

// Opaque returns the color without alpha component.
func (col RGBA) Opaque() RGBA {
	col.A = 1.0
	return col
}

// LinearFromSRGB converts from SRGBA to RGBA.
func LinearFromSRGB(col color.NRGBA) RGBA {
	af := float32(col.A) / 0xFF
	return RGBA{
		R: sRGBToLinear(float32(col.R)/0xff) * af,
		G: sRGBToLinear(float32(col.G)/0xff) * af,
		B: sRGBToLinear(float32(col.B)/0xff) * af,
		A: af,
	}
}

// NRGBAToRGBA converts from non-premultiplied sRGB color to premultiplied sRGB color.
//
// Each component in the result is `sRGBToLinear(c * alpha)`, where `c`
// is the linear color.
func NRGBAToRGBA(col color.NRGBA) color.RGBA {
	if col.A == 0xFF {
		return color.RGBA(col)
	}
	c := LinearFromSRGB(col)
	return color.RGBA{
		R: uint8(linearTosRGB(c.R)*255 + .5),
		G: uint8(linearTosRGB(c.G)*255 + .5),
		B: uint8(linearTosRGB(c.B)*255 + .5),
		A: col.A,
	}
}

// RGBAToNRGBA converts from premultiplied sRGB color to non-premultiplied sRGB color.
func RGBAToNRGBA(col color.RGBA) color.NRGBA {
	if col.A == 0xFF {
		return color.NRGBA(col)
	}

	linear := RGBA{
		R: sRGBToLinear(float32(col.R) / 0xff),
		G: sRGBToLinear(float32(col.G) / 0xff),
		B: sRGBToLinear(float32(col.B) / 0xff),
		A: float32(col.A) / 0xff,
	}

	return linear.SRGB()
}

// linearTosRGB transforms color value from linear to sRGB.
func linearTosRGB(c float32) float32 {
	// Formula from EXT_sRGB.
	switch {
	case c <= 0:
		return 0
	case 0 < c && c < 0.0031308:
		return 12.92 * c
	case 0.0031308 <= c && c < 1:
		return 1.055*float32(math.Pow(float64(c), 0.41666)) - 0.055
	}

	return 1
}

// sRGBToLinear transforms color value from sRGB to linear.
func sRGBToLinear(c float32) float32 {
	// Formula from EXT_sRGB.
	if c <= 0.04045 {
		return c / 12.92
	} else {
		return float32(math.Pow(float64((c+0.055)/1.055), 2.4))
	}
}

// MulAlpha applies the alpha to the color.
func MulAlpha(c color.NRGBA, alpha uint8) color.NRGBA {
	c.A = uint8(uint32(c.A) * uint32(alpha) / 0xFF)
	return c
}
