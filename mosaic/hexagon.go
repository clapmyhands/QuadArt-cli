package main

import (
	"image/color"
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
	dc.SetRGBA(0, 0, 0, 0.4)
	dc.StrokePreserve()
	dc.SetColor(h.color)
	dc.Fill()
}
