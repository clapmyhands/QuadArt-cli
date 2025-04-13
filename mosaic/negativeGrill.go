package main

import (
	"image"
	"image/color"
	"math"
)

func TileWithNegativeGrill(img *image.RGBA, width, height, shiftRatio float64, negativeColor color.Color) []Grill {
	Py := func(x, y int, height, width, shiftRatio float64) float64 {
		if width < height {
			height := height * 1.15
			grillAdjust := math.Mod(shiftRatio*height*float64(x), height)
			return float64(y)*height + grillAdjust
		} else {
			return float64(y) * (height + math.Max(0.4*height, 4))
		}
	}

	Px := func(x, y int, height, width, shiftRatio float64) float64 {
		if width > height {
			width := width * 1.15
			grillAdjust := math.Mod(shiftRatio*width*float64(y), width)
			return float64(x)*width + grillAdjust
		} else {
			return float64(x) * (width + math.Max(0.4*width, 4))
		}
	}
	// prepare Grills
	tiles := make([]Grill, 0)
	var y, x int
	px, py := 0.0, 0.0
	for y = -1; Py(x, y, height, width, shiftRatio) <= float64(img.Bounds().Dy())+height; y++ {
		for x = -1; Px(x, y, height, width, shiftRatio) <= float64(img.Bounds().Dx())+width; x++ {
			py = Py(x, y, height, width, shiftRatio)
			px = Px(x, y, height, width, shiftRatio)

			tiles = append(tiles, Grill{
				x:        px,
				y:        py,
				color:    negativeColor,
				width:    width,
				height:   height,
				avgError: 0,
			})
		}
	}
	return tiles
}
