package main

import (
	"image"
	"image/color"
	"math"
	"quadart-cli/util"

	"github.com/fogleman/gg"
)

type Grill struct {
	x, y          float64
	color         color.Color
	width, height float64
	avgError      float64
}

func (h Grill) Less(other util.HeapItem) bool {
	return h.avgError < other.(Grill).avgError
}

func (h Grill) More(other util.HeapItem) bool {
	return h.avgError > other.(Grill).avgError
}

func (h Grill) Draw(dc *gg.Context) {
	dc.NewSubPath()
	dc.LineTo(h.x+h.width/2.0, h.y+h.height/2.0)
	dc.LineTo(h.x-h.width/2.0, h.y+h.height/2.0)
	dc.LineTo(h.x-h.width/2.0, h.y-h.height/2.0)
	dc.LineTo(h.x+h.width/2.0, h.y-h.height/2.0)
	dc.ClosePath()
	// dc.SetRGBA(0, 0, 0, 0.4)
	dc.SetColor(h.color)
	// dc.SetLineWidth(0.4)
	// dc.StrokePreserve()
	// dc.SetColor(h.color)
	dc.Fill()
}

func TileWithGrill(img *image.RGBA, width, height, shiftRatio float64) []Grill {
	Py := func(x, y int, height, width, shiftRatio float64) float64 {
		if width < height {
			height := height * 1.1
			grillAdjust := math.Mod(shiftRatio*height*float64(x), height)
			return float64(y)*height + grillAdjust
		} else {
			return float64(y) * (height + math.Max(0.15*height, 1))
		}
	}

	Px := func(x, y int, height, width, shiftRatio float64) float64 {
		if width > height {
			width := width * 1.1
			grillAdjust := math.Mod(shiftRatio*width*float64(y), width)
			return float64(x)*width + grillAdjust
		} else {
			return float64(x) * (width + math.Max(0.15*width, 1))
		}
	}
	avgImgColor := util.CalcAvgColor(img)
	// prepare Grills
	tiles := make([]Grill, 0)
	var y, x int
	px, py := 0.0, 0.0
	for y = -1; Py(x, y, height, width, shiftRatio) <= float64(img.Bounds().Dy())+height; y++ {
		for x = -1; Px(x, y, height, width, shiftRatio) <= float64(img.Bounds().Dx())+width; x++ {
			py = Py(x, y, height, width, shiftRatio)
			px = Px(x, y, height, width, shiftRatio)

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

			tiles = append(tiles, Grill{
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
