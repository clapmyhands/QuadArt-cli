package util

import (
	"fmt"
	"image/color"
	"math"
	"strconv"
	"strings"
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

// HexToColorWithAlpha converts a hex string (e.g., "#RRGGBBAA" or "RRGGBBAA") to a color.RGBA.
// It also handles the 8-digit format with an alpha channel.
// It returns an error if the hex string is invalid.
func HexToColorWithAlpha(hex string) (color.RGBA, error) {
	hex = strings.TrimPrefix(hex, "#")
	var size int
	if len(hex) == 6 {
		size = 3
	} else if len(hex) == 8 {
		size = 4
	} else {
		return color.RGBA{}, fmt.Errorf("invalid hex color format (must be 6 or 8 digits): %s", hex)
	}

	parse := func(s string) (uint8, error) {
		val, err := strconv.ParseUint(s, 16, 8)
		return uint8(val), err
	}

	var r, g, b, a uint8

	switch size {
	case 3:
		r, _ = parse(hex[0:2])
		g, _ = parse(hex[2:4])
		b, _ = parse(hex[4:6])
		a = 0xff // Default alpha
	case 4:
		r, _ = parse(hex[0:2])
		g, _ = parse(hex[2:4])
		b, _ = parse(hex[4:6])
		a, _ = parse(hex[6:8])
	}

	return color.RGBA{R: r, G: g, B: b, A: a}, nil
}
