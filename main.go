package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"

	"quadart-cli/quadtree"

	"github.com/fogleman/gg"
)

const filename = "./patrick.png"
const errThreshold = float64(50_000)
const drawLoop = 1 // every 100th loop starting from 0

func main() {
	// open image
	reader, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer func(reader *os.File) {
		_ = reader.Close()
	}(reader)

	originalImg, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	// setup new virtual img and drawing context
	copyImg := image.NewRGBA(originalImg.Bounds())
	draw.Draw(copyImg, originalImg.Bounds(), originalImg, originalImg.Bounds().Min, draw.Src)
	dc := gg.NewContext(originalImg.Bounds().Dx(), originalImg.Bounds().Dy())

	imgRectQT := quadtree.NewImgRectQuadTree(copyImg)
	var crossLines []image.Rectangle
	for imgRect, cnt := imgRectQT.PopImgRect(), 0; imgRect.AvgError > errThreshold; imgRect = imgRectQT.PopImgRect() {
		for _, quadrant := range split2Quadrant(imgRect.Rect) {
			imgQuadrant := imgRectQT.ExtractAndPush(quadrant)
			drawRectangle(dc, imgQuadrant.Rect, imgQuadrant.AvgColor)
		}
		crossLines = append(crossLines, imgRect.Rect)

		if cnt%drawLoop == 0 {
			log.Printf("Loop: %d\tAvgError: %f", cnt, imgRect.AvgError)

			// TODO: this could probable be moved to a worker/goroutine
			// copy and fill stroke before printing image
			// this is done such that stroke is always done on the last step as 1 action
			// avoids issue of multiple line stacking causing deeper line
			exportDC := gg.NewContextForImage(dc.Image())
			// for _, cross := range crossLines {
			// 	drawCross(exportDC, cross, 0.5)
			// }
			// exportDC.SetRGBA(0, 0, 0, 0.4)
			// exportDC.Stroke()

			// gg.SaveJPG(fmt.Sprintf("./out/jpeg/%d.jpeg", cnt), dc.Image(), 100)
			gg.SavePNG(fmt.Sprintf("./out/png/%d.png", cnt), exportDC.Image())
		}

		cnt++
		if cnt > 100 {
			break
		}
	}
}

func drawRectangle(dc *gg.Context, rect image.Rectangle, color color.Color) {
	dc.DrawRectangle(
		float64(rect.Min.X),
		float64(rect.Min.Y),
		float64(rect.Dx()),
		float64(rect.Dy()),
	)
	dc.SetColor(color)
	dc.Fill()
}

func split2Quadrant(rect image.Rectangle) []image.Rectangle {
	var (
		minX = rect.Min.X
		minY = rect.Min.Y
		midX = rect.Min.X + rect.Dx()/2
		midY = rect.Min.Y + rect.Dy()/2
		maxX = rect.Max.X
		maxY = rect.Max.Y
	)

	return []image.Rectangle{
		image.Rect(minX, minY, midX, midY),
		image.Rect(midX, minY, maxX, midY),
		image.Rect(minX, midY, midX, maxY),
		image.Rect(midX, midY, maxX, maxY),
	}
}

func drawCross(dc *gg.Context, r image.Rectangle, lineWidth float64) {
	midX := r.Min.X + r.Dx()/2
	midY := r.Min.Y + r.Dy()/2
	dc.DrawLine(float64(midX), float64(r.Min.Y), float64(midX), float64(r.Max.Y))
	dc.DrawLine(float64(r.Min.X), float64(midY), float64(r.Max.X), float64(midY))
	dc.SetLineWidth(lineWidth)
}
