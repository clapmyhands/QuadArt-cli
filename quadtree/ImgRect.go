package quadtree

import (
	"image"
	"image/color"
	"image/draw"
)

type ImgRect struct {
	// Rect denotes the bounding box of the image
	Rect image.Rectangle
	// AvgColor denotes the average color for an image
	AvgColor color.RGBA64
	// AvgError denotes MSE of image to its color average
	AvgError float64
}

func extractImageRect(img *image.RGBA, rect image.Rectangle) ImgRect {
	trimmedImg := moveTopLeft(img.SubImage(rect), rect)
	avgColor := calcAvgColor(trimmedImg)
	mse := calcImgToColorMSE(trimmedImg, avgColor)

	return ImgRect{
		Rect:     rect,
		AvgColor: avgColor,
		AvgError: mse,
	}
}

func moveTopLeft(img image.Image, rect image.Rectangle) *image.RGBA {
	newImg := image.NewRGBA(rect.Sub(rect.Min))
	draw.Draw(newImg, newImg.Rect, img, rect.Min, draw.Src)
	return newImg
}
