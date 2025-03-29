package main

import (
	"image"
	"image/color"
	"quadart-cli/util"

	"github.com/fogleman/gg"
)

type OverlappingCircle struct {
	x, y     float64
	color    color.Color
	alpha    float64
	radius   float64
	avgError float64
}

func (h OverlappingCircle) Less(other util.HeapItem) bool {
	return h.avgError < other.(OverlappingCircle).avgError
}

func (h OverlappingCircle) More(other util.HeapItem) bool {
	return h.avgError > other.(OverlappingCircle).avgError
}

func (h OverlappingCircle) Draw(dc *gg.Context) {
	dc.DrawCircle(h.x, h.y, h.radius)
	// dc.SetLineWidth(0.5)
	dc.SetColor(util.SetAlpha(h.color, h.alpha))
	// dc.StrokePreserve()
	// dc.SetColor(util.SetAlpha(h.color, 0.6))
	dc.Fill()
}

func TileWithOverlappingCircle(img *image.RGBA, radius, alpha, spaceMult float64) []OverlappingCircle {
	avgImgColor := util.CalcAvgColor(img)
	// prepare OverlappingCircles
	tiles := make([]OverlappingCircle, 0)
	var y, x int
	px, py := 0.0, 0.0
	for y = 0; float64(y)*radius <= float64(img.Bounds().Dy())+radius; y++ {
		for x = 0; float64(x)*radius <= float64(img.Bounds().Dx())+radius; x++ {
			py = float64(y) * radius * spaceMult
			px = float64(x) * radius * spaceMult

			subImg := util.ExtractRectSubImg(
				img,
				util.CalcRectangle(image.Pt(int(px), int(py)), radius),
			)
			if subImg.Rect.Empty() {
				// extracted rectangle is outside original image
				continue
			}
			avgColor := util.CalcAvgColor(subImg)

			tiles = append(tiles, OverlappingCircle{
				x:        px,
				y:        py,
				color:    avgColor,
				alpha:    alpha,
				radius:   radius,
				avgError: util.CalcImgToColorMSE(subImg, avgImgColor),
			})
		}
	}
	return tiles
}
