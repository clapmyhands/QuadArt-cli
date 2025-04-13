package main

import (
	"container/heap"
	"fmt"
	"image"
	"image/draw"
	"log"
	"os"
	"quadart-cli/util"

	flag "github.com/spf13/pflag"

	"github.com/fogleman/gg"
)

type RunParameter struct {
	inputFilepath       string
	outputFolder        string
	finalOutputFilename string
	color               string
	shape               string
	errThreshold        float64
	radius              float64
	width               float64
	height              float64
	alpha               float64
	spaceMultiplier     float64
	shiftRatio          float64
}

func parseArgs() *RunParameter {
	param := &RunParameter{}

	flag.StringVarP(&param.inputFilepath, "input", "i", "", "Input Filepath")
	flag.StringVarP(&param.outputFolder, "outputFolder", "o", "out/mosaic", "Output Folder")
	flag.StringVarP(&param.finalOutputFilename, "finalOutputFilename", "f", "final", "Final Output Filename")
	flag.StringVarP(&param.color, "color", "c", "ffffff", "Background Color in normal, negative color if negative shape")
	flag.StringVarP(&param.shape, "shape", "s", "hexagon", "Shape for tiling")
	flag.Float64VarP(&param.errThreshold, "threshold", "t", 1_000, "Error Threshold before stopping")
	flag.Float64VarP(&param.radius, "radius", "r", 10, "Radius for calculating length of edges and error finding")
	flag.Float64VarP(&param.width, "width", "w", 10, "Width of shape")
	flag.Float64VarP(&param.height, "height", "h", 10, "Height of shape")
	flag.Float64VarP(&param.alpha, "alpha", "a", 0.6, "Alpha channel 0-1")
	flag.Float64Var(&param.spaceMultiplier, "spaceMultiplier", 1.05, "Space multiplier for overlapping shapes")
	flag.Float64Var(&param.shiftRatio, "shiftRatio", 0.5, "Shift Ratio for Grills")
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

	w := originalImg.Bounds().Dx()
	h := originalImg.Bounds().Dy()
	dc := gg.NewContext(w, h)

	// set background
	dc.DrawRectangle(0, 0, float64(originalImg.Bounds().Dx()), float64(originalImg.Bounds().Dy()))
	dc.SetHexColor(runParam.color)
	dc.Fill()

	switch runParam.shape {
	case "hexagon":
		tiles := TileWithHexagon(copyImg, runParam.radius)
		HeapTiling(dc, util.ToHeap(tiles), runParam.errThreshold)
	case "triangle":
		tiles := TileWithTriangle(copyImg, runParam.radius)
		ExhaustiveTiling(dc, ToShapes(tiles))
	case "square":
		tiles := TileWithSquare(copyImg, runParam.width, runParam.height)
		ExhaustiveTiling(dc, ToShapes(tiles))
	case "diamond":
		tiles := TileWithDiamond(copyImg, runParam.width, runParam.height)
		ExhaustiveTiling(dc, ToShapes(tiles))
	case "oCircle":
		tiles := TileWithOverlappingCircle(copyImg, runParam.radius, runParam.alpha, runParam.spaceMultiplier)
		ExhaustiveTiling(dc, ToShapes(tiles))
	case "grill":
		tiles := TileWithGrill(copyImg, runParam.width, runParam.height, runParam.shiftRatio)
		ExhaustiveTiling(dc, ToShapes(tiles))
	case "negativeGrill":
		dc = gg.NewContextForImage(originalImg)
		negativeColor, _ := util.HexToColorWithAlpha(runParam.color)
		tiles := TileWithNegativeGrill(copyImg, runParam.width, runParam.height, runParam.shiftRatio, negativeColor)
		ExhaustiveTiling(dc, ToShapes(tiles))
	default:
		log.Fatalf("Unsupported shape: %s", runParam.shape)
	}
	dc.SavePNG(fmt.Sprintf("%s/%s.png", runParam.outputFolder, runParam.finalOutputFilename))
}

func HeapTiling(dc *gg.Context, shapeHeap util.MaxHeap, errThreshold float64) {
	var x, y = dc.Width(), dc.Height()
	heap.Init(&shapeHeap)
	log.Printf("x: %d, y: %d, heapLen: %d", x, y, len(shapeHeap))
	// draw
	for i := 0; len(shapeHeap) > 0; i++ {
		hex := heap.Pop(&shapeHeap).(Hexagon)
		if i%10 == 0 {
			log.Printf("loop: %d,\tx: %f,\ty: %f,\terr: %f,\tcolor: %v,\t len:%d", i, hex.x, hex.y, hex.avgError, hex.color, len(shapeHeap))
		}

		if hex.avgError < errThreshold {
			break
		}
		hex.Draw(dc)
	}
}

func ExhaustiveTiling(dc *gg.Context, shapes []Shape) {
	var x, y = dc.Width(), dc.Height()
	log.Printf("x: %d, y: %d, len: %d", x, y, len(shapes))
	// draw
	for _, shape := range shapes {
		// log.Printf("%+v", shape)
		shape.Draw(dc)
	}
}
