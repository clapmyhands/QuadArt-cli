package util

import (
	"image/color"
	"math"
)

func SetAlpha(c color.Color, alpha float64) color.NRGBA {
	r, g, b, _ := c.RGBA() // Get the original color's RGBA values

	alpha255 := uint8(math.Round(alpha * 255))

	// Convert to 8-bit values (0-255)
	return color.NRGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: alpha255,
	}
}
