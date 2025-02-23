package quadtree

import (
	"image"
	"image/color"
	"math"
)

func calcAvgColor(img image.Image) color.RGBA64 {
	var (
		area                          = uint64(img.Bounds().Dx() * img.Bounds().Dy())
		cumR, cumG, cumB, cumA uint64 = 0, 0, 0, 0
		minBound                      = img.Bounds().Min
		maxBound                      = img.Bounds().Max
	)

	for y := minBound.Y; y < maxBound.Y; y++ {
		for x := minBound.X; x < maxBound.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			cumR += uint64(r)
			cumG += uint64(g)
			cumB += uint64(b)
			cumA += uint64(a)
		}
	}

	return color.RGBA64{
		R: uint16(cumR / area),
		G: uint16(cumG / area),
		B: uint16(cumB / area),
		A: uint16(cumA / area),
	}
}

func calcImgToColorMSE(img image.Image, c color.RGBA64) float64 {
	var (
		minBound = img.Bounds().Min
		maxBound = img.Bounds().Max
		area     = float64(img.Bounds().Size().X * img.Bounds().Size().Y)
		mse      = float64(0)
	)

	for y := minBound.Y; y < maxBound.Y; y++ {
		for x := minBound.X; x < maxBound.X; x++ {
			mse += calcColorMSE(img.At(x, y), c) / area
		}
	}
	return math.Sqrt(mse) * math.Sqrt(area) // weighted average by size
}

func calcColorMSE(c1, c2 color.Color) float64 {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()
	rDiff, gDiff, bDiff := r1-r2, g1-g2, b1-b2
	// RGB -> Grayscale (Standard NTSC Conversion)
	return 0.299*float64(rDiff*rDiff) + 0.587*float64(gDiff*gDiff) + 0.114*float64(bDiff*bDiff)
}
