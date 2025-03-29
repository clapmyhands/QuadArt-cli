package main

import (
	"image"
	"image/color"
	"math"
	"quadart-cli/util"

	"github.com/fogleman/gg"
)

/*
	Length of the edge: r * math.sqrt(3)
	Height of the triangle: r * 1.5
	​Centroid height from the base of the triangle: r/2 = height/3
*/

type Triangle struct {
	x, y       float64
	color      color.Color
	radius     float64
	pointingUp bool
	avgError   float64
}

func (h Triangle) Less(other util.HeapItem) bool {
	return h.avgError < other.(Triangle).avgError
}

func (h Triangle) More(other util.HeapItem) bool {
	return h.avgError > other.(Triangle).avgError
}

func (h Triangle) Draw(dc *gg.Context) {
	rotation := math.Pi
	if h.pointingUp {
		rotation = 0.0
	}
	dc.DrawRegularPolygon(3, h.x, h.y, h.radius, rotation)
	// r, g, b, a := h.color.RGBA()
	// dc.SetColor(color.RGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: uint16(a)})
	dc.SetColor(h.color)
	dc.SetLineWidth(0.5)
	dc.StrokePreserve()
	// dc.SetColor(h.color)
	dc.Fill()
}

func TileWithTriangle(img *image.RGBA, radius float64) []Triangle {
	avgImgColor := util.CalcAvgColor(img)
	var triangleH = radius * 1.5
	var triangleW = radius * math.Sqrt(3)
	// prepare hexagons
	tiles := make([]Triangle, 0)
	var y, x int
	px, py := 0.0, 0.0
	for y = 0; float64(y)*triangleH < float64(img.Bounds().Dy())+radius; y++ {
		for x = 0; float64(x)*triangleW/2.0 < float64(img.Bounds().Dx())+radius; x++ {
			pointingUp := x%2 == 0
			if y%2 == 1 {
				pointingUp = !pointingUp
			}
			if x%2 == y%2 {
				py = float64(y) * triangleH // magic number to allow pixel to breathe
			} else {
				py = float64(y)*triangleH - triangleH/3.0
			}
			px = float64(x) * triangleW / 2.0

			subImg := util.ExtractRectSubImg(
				img,
				util.CalcRectangle(image.Pt(int(px), int(py)), radius),
			)
			if subImg.Rect.Empty() {
				// extracted rectangle is outside original image
				continue
			}
			avgColor := util.CalcAvgColor(subImg)

			tiles = append(tiles, Triangle{
				x:          px,
				y:          py,
				color:      avgColor,
				radius:     radius,
				pointingUp: pointingUp,
				avgError:   util.CalcImgToColorMSE(subImg, avgImgColor),
			})
		}
	}
	return tiles
}
