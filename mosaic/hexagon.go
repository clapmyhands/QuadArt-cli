package main

import (
	"image"
	"image/color"
	"math"
	"quadart-cli/util"

	"github.com/fogleman/gg"
)

type Hexagon struct {
	x, y     float64
	color    color.Color
	radius   int
	avgError float64
}

func (h Hexagon) Less(other util.HeapItem) bool {
	return h.avgError < other.(Hexagon).avgError
}

func (h Hexagon) More(other util.HeapItem) bool {
	return h.avgError > other.(Hexagon).avgError
}

func (h Hexagon) Draw(dc *gg.Context) {
	dc.DrawRegularPolygon(6, h.x, h.y, float64(h.radius), 0)
	// dc.SetRGBA(0, 0, 0, 0.4)
	dc.SetColor(h.color)
	dc.StrokePreserve()
	dc.Fill()
}

func TileWithHexagon(img *image.RGBA, radius int) []Hexagon {
	avgImgColor := util.CalcAvgColor(img)
	var hexHeight = float64(radius) * math.Sqrt(3)
	// prepare hexagons
	tiles := make([]Hexagon, 0)
	var y, x int
	px, py := 0.0, 0.0
	for y = 0; float64(y)*hexHeight < float64(img.Bounds().Dy()+radius); y++ {
		for x = 0; float64(x)*1.5*float64(radius) < float64(img.Bounds().Dx()+radius); x++ {
			if x%2 == 1 {
				py = float64(y)*hexHeight + (hexHeight / 2)
			} else {
				py = float64(y) * hexHeight
			}
			px = float64(x) * 1.5 * float64(radius)

			subImg := util.ExtractRectSubImg(
				img,
				util.CalcRectangle(image.Pt(int(px), int(py)), float64(radius)),
			)
			if subImg.Rect.Empty() {
				// extracted rectangle is outside original image
				continue
			}
			avgColor := util.CalcAvgColor(subImg)

			tiles = append(tiles, Hexagon{
				x:        px,
				y:        py,
				color:    avgColor,
				radius:   radius,
				avgError: util.CalcImgToColorMSE(subImg, avgImgColor),
			})
		}
	}
	return tiles
}
