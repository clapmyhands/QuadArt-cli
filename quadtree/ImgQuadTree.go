package main

import (
	"container/heap"
	"image"
	"quadart-cli/util"
)

type ImgQuadTree struct {
	img     *image.RGBA
	maxHeap *util.MaxHeap
}

func NewImgQuadTree(img *image.RGBA) *ImgQuadTree {
	imgQuadTree := &ImgQuadTree{
		img:     img,
		maxHeap: &util.MaxHeap{},
	}
	heap.Init(imgQuadTree.maxHeap)
	heap.Push(imgQuadTree.maxHeap, extractImageRect(img, img.Bounds()))
	return imgQuadTree
}

// ExtractAndPush - QoL method to push new record
func (h *ImgQuadTree) ExtractAndPush(rect image.Rectangle) ImgRect {
	ire := extractImageRect(h.img, rect)
	heap.Push(h.maxHeap, ire)
	return ire
}

// PopImgRect - QoL method to pop new record
func (h *ImgQuadTree) PopImgRect() ImgRect {
	return heap.Pop(h.maxHeap).(ImgRect)
}
