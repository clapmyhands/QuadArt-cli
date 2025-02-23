package quadtree

import (
	"container/heap"
	"image"
)

type ImgRectQuadTree struct {
	img         *image.RGBA
	imgRectList []ImgRect
}

func NewImgRectQuadTree(img *image.RGBA) *ImgRectQuadTree {
	imgRectQuadTree := &ImgRectQuadTree{
		img:         img,
		imgRectList: make([]ImgRect, 0),
	}
	heap.Init(imgRectQuadTree)
	heap.Push(imgRectQuadTree, extractImageRect(img, img.Bounds()))
	return imgRectQuadTree
}

func (h ImgRectQuadTree) Len() int { return len(h.imgRectList) }

// Less - max heap for AvgError
func (h ImgRectQuadTree) Less(i, j int) bool {
	return h.imgRectList[i].AvgError > h.imgRectList[j].AvgError
}

func (h ImgRectQuadTree) Swap(i, j int) {
	h.imgRectList[i], h.imgRectList[j] = h.imgRectList[j], h.imgRectList[i]
}

func (h *ImgRectQuadTree) Push(x interface{}) {
	(*h).imgRectList = append((*h).imgRectList, x.(ImgRect))
}

func (h *ImgRectQuadTree) Pop() interface{} {
	old := (*h).imgRectList
	n := len(old)
	x := old[n-1]
	(*h).imgRectList = old[0 : n-1]
	return x
}

// ExtractAndPush - QoL method to push new record
func (h *ImgRectQuadTree) ExtractAndPush(rect image.Rectangle) ImgRect {
	ire := extractImageRect(h.img, rect)
	heap.Push(h, ire)
	return ire
}

// PopImgRect - QoL method to pop new record
func (h *ImgRectQuadTree) PopImgRect() ImgRect {
	return heap.Pop(h).(ImgRect)
}
