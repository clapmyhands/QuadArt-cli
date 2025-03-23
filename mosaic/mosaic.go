package main

import (
	"container/heap"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"
	"os"
	"quadart-cli/util"

	flag "github.com/spf13/pflag"

	"github.com/fogleman/gg"
)

var colorWhite = color.RGBA64{
	R: 0xffff,
	G: 0xffff,
	B: 0xffff,
	A: 0xffff,
}

type RunParameter struct {
	inputFilepath string
	outputFolder  string
	errThreshold  float64
	radius        int
}

func parseArgs() *RunParameter {
	param := &RunParameter{}

	flag.StringVarP(&param.inputFilepath, "input", "i", "", "Input Filepath")
	flag.StringVarP(&param.outputFolder, "output", "o", "out/mosaic", "Output Folder")
	flag.Float64VarP(&param.errThreshold, "threshold", "t", 1_000, "Error Threshold before stopping")
	flag.IntVarP(&param.radius, "radius", "r", 10, "Radius for calculating length of edges and error finding")
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

	var hexHeight = float64(runParam.radius) * math.Sqrt(3)

	// setup new virtual img and drawing context
	copyImg := image.NewRGBA(originalImg.Bounds())
	draw.Draw(copyImg, originalImg.Bounds(), originalImg, originalImg.Bounds().Min, draw.Src)

	w := originalImg.Bounds().Dx()
	h := originalImg.Bounds().Dy()

	dc := gg.NewContext(w, h)

	// set background
	dc.DrawRectangle(0, 0, float64(originalImg.Bounds().Dx()), float64(originalImg.Bounds().Dy()))
	dc.SetColor(color.White)
	dc.Fill()

	// prepare hexagons
	HexagonHeap := make(util.MaxHeap, 0)
	var y, x int
	px, py := 0.0, 0.0
	for y = 0; float64(y)*hexHeight < float64(h+runParam.radius); y++ {
		for x = 0; float64(x)*1.5*float64(runParam.radius) < float64(w+runParam.radius); x++ {
			if x%2 == 1 {
				py = float64(y)*hexHeight + (hexHeight / 2)
			} else {
				py = float64(y) * hexHeight
			}
			px = float64(x) * 1.5 * float64(runParam.radius)

			subImg := util.ExtractRectSubImg(
				copyImg,
				util.CalcRectangle(image.Pt(int(px), int(py)), float64(runParam.radius)),
			)
			if subImg.Rect.Empty() {
				// extracted rectangle is outside original image
				continue
			}
			avgColor := util.CalcAvgColor(subImg)

			HexagonHeap = append(HexagonHeap, Hexagon{
				x:        px,
				y:        py,
				color:    avgColor,
				radius:   runParam.radius,
				avgError: util.CalcImgToColorMSE(subImg, colorWhite),
			})
		}
	}
	heap.Init(&HexagonHeap)

	log.Printf("x: %d, y: %d, heapLen: %d", x, y, len(HexagonHeap))

	// draw
	for i := 0; len(HexagonHeap) > 0; i++ {
		hex := heap.Pop(&HexagonHeap).(Hexagon)
		if i%10 == 0 {
			log.Printf("loop: %d,\tx: %f,\ty: %f,\terr: %f,\tcolor: %v,\t len:%d", i, hex.x, hex.y, hex.avgError, hex.color, len(HexagonHeap))
		}

		if hex.avgError < runParam.errThreshold {
			break
		}
		hex.Draw(dc)
	}
	dc.SavePNG(fmt.Sprintf("%s/final.png", runParam.outputFolder))
}
