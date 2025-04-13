package main

import (
	"image"
	"image/color"
	"quadart-cli/util"
)

type ImgRect struct {
	// Rect denotes the bounding box of the image
	Rect image.Rectangle
	// AvgColor denotes the average color for an image
	AvgColor color.RGBA64
	// AvgError denotes MSE of image to its color average
	AvgError float64
}

func (i ImgRect) Less(other util.HeapItem) bool {
	return i.AvgError < other.(ImgRect).AvgError
}

func (i ImgRect) More(other util.HeapItem) bool {
	return i.AvgError > other.(ImgRect).AvgError
}

func extractImageRect(img *image.RGBA, rect image.Rectangle) ImgRect {
	trimmedImg := util.ExtractRectSubImg(img, rect)
	AvgColor := util.CalcAvgColor(trimmedImg)
	mse := util.CalcImgToColorMSE(trimmedImg, AvgColor)

	return ImgRect{
		Rect:     rect,
		AvgColor: AvgColor,
		AvgError: mse,
	}
}
