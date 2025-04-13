package main

import (
	"image"
	"image/color"
	"quadart-cli/util"

	"github.com/fogleman/gg"
)

type Diamond struct {
	x, y          float64
	width, height float64
	color         color.Color
	avgError      float64
}

func (h Diamond) Less(other util.HeapItem) bool {
	return h.avgError < other.(Diamond).avgError
}

func (h Diamond) More(other util.HeapItem) bool {
	return h.avgError > other.(Diamond).avgError
}

func (h Diamond) Draw(dc *gg.Context) {
	dc.NewSubPath()
	dc.LineTo(h.x+h.width/2.0, h.y)
	dc.LineTo(h.x, h.y+h.height/2.0)
	dc.LineTo(h.x-h.width/2.0, h.y)
	dc.LineTo(h.x, h.y-h.height/2.0)
	dc.ClosePath()
	// dc.SetRGBA(0, 0, 0, 0.4)
	// dc.SetColor(h.color)
	dc.SetColor(h.color)
	dc.SetLineWidth(0.2)
	dc.StrokePreserve()
	dc.Fill()
}

func TileWithDiamond(img *image.RGBA, width, height float64) []Diamond {
	avgImgColor := util.CalcAvgColor(img)
	// prepare Diamonds
	tiles := make([]Diamond, 0)
	var y, x int
	px, py := 0.0, 0.0
	for y = 0; float64(y)*height <= float64(img.Bounds().Dy()); y++ {
		for x = 0; float64(x)*width/2.0 <= float64(img.Bounds().Dx()); x++ {
			if x%2 == 1 {
				py = float64(y)*height + (height / 2.0)
			} else {
				py = float64(y) * height
			}
			px = float64(x) * width / 2.0

			subImg := util.ExtractRectSubImg(
				img,
				image.Rect(
					// this will take the inside square of the diamond
					// more contrasty
					int(px-width/2.0),
					int(py-height/2.0),
					int(px+width/2.0),
					int(py+height/2.0),

					// this will take the average of width+height square outside the diamond
					// bigger input square but less contrast as result
					// int(px-(width+height/4.0)),
					// int(py-(width+height/4.0)),
					// int(px+(width+height/4.0)),
					// int(py+(width+height/4.0)),
				),
			)
			if subImg.Rect.Empty() {
				// extracted rectangle is outside original image
				continue
			}
			avgColor := util.CalcAvgColor(subImg)

			tiles = append(tiles, Diamond{
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
