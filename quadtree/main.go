package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"

	flag "github.com/spf13/pflag"

	"github.com/fogleman/gg"
)

type RunParameter struct {
	inputFilepath       string
	outputFolder        string
	finalOutputFilename string
	backgroundColor     string
	errThreshold        float64
	radius              float64
	alpha               float64
	drawIteration       int
}

func parseArgs() *RunParameter {
	param := &RunParameter{}

	flag.StringVarP(&param.inputFilepath, "input", "i", "", "Input Filepath")
	flag.StringVarP(&param.outputFolder, "outputFolder", "o", "out/quadtree", "Output Folder")
	flag.StringVarP(&param.finalOutputFilename, "finalOutputFilename", "f", "final", "Final Output Filename")
	flag.StringVarP(&param.backgroundColor, "backgroundColor", "b", "ffffff", "Background Color in Hex")
	flag.Float64VarP(&param.errThreshold, "threshold", "t", 1_000, "Error Threshold before stopping")
	flag.Float64VarP(&param.radius, "radius", "r", 0, "Radius for rounded rectangle")
	flag.Float64VarP(&param.alpha, "alpha", "a", 0.6, "Alpha channel 0-1")
	flag.IntVarP(&param.drawIteration, "drawIteration", "d", 10, "Save image every drawIteration-th")
	flag.Parse()

	if len(param.inputFilepath) == 0 {
		log.Fatalf("Input filepath cannot be empty")
	}

	return param
}

func main() {
	runParam := parseArgs()

	// open image
	reader, err := os.Open(runParam.inputFilepath)
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

	dc.DrawRectangle(0, 0, float64(originalImg.Bounds().Dx()), float64(originalImg.Bounds().Dy()))
	dc.SetHexColor(runParam.backgroundColor)
	dc.Fill()

	imgRectQT := NewImgQuadTree(copyImg)
	var crossLines []image.Rectangle
	for imgRect, cnt := imgRectQT.PopImgRect(), 0; imgRect.AvgError > runParam.errThreshold; imgRect = imgRectQT.PopImgRect() {
		for _, quadrant := range split2Quadrant(imgRect.Rect) {
			imgQuadrant := imgRectQT.ExtractAndPush(quadrant)
			drawRectangle(dc, imgQuadrant.Rect, imgQuadrant.AvgColor, runParam.radius)
		}
		crossLines = append(crossLines, imgRect.Rect)

		if cnt%runParam.drawIteration == 0 {
			log.Printf("Loop: %d\tAvgError: %f", cnt, imgRect.AvgError)

			// TODO: this could probable be moved to a worker/goroutine
			// copy and fill stroke before printing image
			// this is done such that stroke is always done on the last step as 1 action
			// avoids issue of multiple line stacking causing deeper line
			exportDC := drawCrossLines(dc, crossLines, runParam.alpha)
			exportDC.SavePNG(fmt.Sprintf("%s/%d.png", runParam.outputFolder, cnt))
		}

		if imgRect.AvgError < runParam.errThreshold {
			break
		}
		cnt++
	}
	exportDC := drawCrossLines(dc, crossLines, runParam.alpha)
	exportDC.SavePNG(fmt.Sprintf("%s/%s.png", runParam.outputFolder, runParam.finalOutputFilename))
}

func drawCrossLines(dc *gg.Context, crossLines []image.Rectangle, alpha float64) *gg.Context {
	exportDC := gg.NewContextForImage(dc.Image())
	for _, cross := range crossLines {
		drawCross(exportDC, cross, 0.5)
	}
	exportDC.SetRGBA(0, 0, 0, alpha)
	exportDC.Stroke()
	return exportDC
}

func drawRectangle(dc *gg.Context, rect image.Rectangle, color color.Color, radius float64) {
	dc.DrawRoundedRectangle(
		float64(rect.Min.X),
		float64(rect.Min.Y),
		float64(rect.Dx()),
		float64(rect.Dy()),
		radius,
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
