package main

import (
	"image"
	"image/color"
	"quadart-cli/util"

	"github.com/fogleman/gg"
)

type Square struct {
	x, y          float64
	color         color.Color
	width, height float64
	avgError      float64
}

func (h Square) Less(other util.HeapItem) bool {
	return h.avgError < other.(Square).avgError
}

func (h Square) More(other util.HeapItem) bool {
	return h.avgError > other.(Square).avgError
}

func (h Square) Draw(dc *gg.Context) {
	dc.NewSubPath()
	dc.LineTo(h.x+h.width/2.0, h.y+h.height/2.0)
	dc.LineTo(h.x-h.width/2.0, h.y+h.height/2.0)
	dc.LineTo(h.x-h.width/2.0, h.y-h.height/2.0)
	dc.LineTo(h.x+h.width/2.0, h.y-h.height/2.0)
	dc.ClosePath()
	// dc.SetRGBA(0, 0, 0, 0.4)
	dc.SetColor(h.color)
	dc.SetLineWidth(0.5)
	dc.StrokePreserve()
	// dc.SetColor(h.color)
	dc.Fill()
}

func TileWithSquare(img *image.RGBA, width, height float64) []Square {
	avgImgColor := util.CalcAvgColor(img)
	// prepare Squares
	tiles := make([]Square, 0)
	var y, x int
	px, py := 0.0, 0.0
	for y = 0; float64(y)*height <= float64(img.Bounds().Dy()); y++ {
		for x = 0; float64(x)*width <= float64(img.Bounds().Dx()); x++ {
			py = float64(y) * height
			px = float64(x) * width

			subImg := util.ExtractRectSubImg(
				img,
				image.Rect(
					int(px-width/2.0),
					int(py-height/2.0),
					int(px+width/2.0),
					int(py+height/2.0),
				),
			)
			if subImg.Rect.Empty() {
				// extracted rectangle is outside original image
				continue
			}
			avgColor := util.CalcAvgColor(subImg)

			tiles = append(tiles, Square{
				x:        px,
				y:        py,
				color:    avgColor,
				width:    width,
				height:   height,
				avgError: util.CalcImgToColorMSE(subImg, avgImgColor),
			})
		}
	}
	return tiles
}
