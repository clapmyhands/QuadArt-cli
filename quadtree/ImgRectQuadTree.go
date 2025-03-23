package quadtree

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
	imgRectQuadTree := &ImgQuadTree{
		img:     img,
		maxHeap: &util.MaxHeap{},
	}
	heap.Init(imgRectQuadTree.maxHeap)
	heap.Push(imgRectQuadTree.maxHeap, extractImageRect(img, img.Bounds()))
	return imgRectQuadTree
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
