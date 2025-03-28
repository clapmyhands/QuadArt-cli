package main

import "github.com/fogleman/gg"

type Shape interface {
	Draw(dc *gg.Context)
}

func ToShapes[T Shape](items []T) []Shape {
	result := make([]Shape, len(items))
	for i, item := range items {
		result[i] = item
	}
	return result
}
