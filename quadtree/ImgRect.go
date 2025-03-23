package quadtree

import (
	"image"
	"image/color"
	"quadart-cli/util"
)

type ImgRect struct {
	// rect denotes the bounding box of the image
	rect image.Rectangle
	// avgColor denotes the average color for an image
	avgColor color.RGBA64
	// avgError denotes MSE of image to its color average
	avgError float64
}

func (i ImgRect) Less(other util.HeapItem) bool {
	return i.avgError < other.(ImgRect).avgError
}

func (i ImgRect) More(other util.HeapItem) bool {
	return i.avgError > other.(ImgRect).avgError
}

func extractImageRect(img *image.RGBA, rect image.Rectangle) ImgRect {
	trimmedImg := util.ExtractRectSubImg(img, rect)
	avgColor := util.CalcAvgColor(trimmedImg)
	mse := util.CalcImgToColorMSE(trimmedImg, avgColor)

	return ImgRect{
		rect:     rect,
		avgColor: avgColor,
		avgError: mse,
	}
}
