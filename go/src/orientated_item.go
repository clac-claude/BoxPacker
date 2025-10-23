package boxpacker

import (
	"fmt"
	"math"
	"sort"
)

// Box packing (3D bin packing, knapsack problem).
//
// @author Doug Wright

// OrientatedItem - An item to be packed in a specific orientation.
type OrientatedItem struct {
	Item              Item
	Width             int
	Length            int
	Depth             int
	SurfaceFootprint  int
	dimensionsAsArray []int
}

var stabilityCache = make(map[string]bool)

// NewOrientatedItem creates a new OrientatedItem
func NewOrientatedItem(item Item, width, length, depth int) *OrientatedItem {
	oi := &OrientatedItem{
		Item:             item,
		Width:            width,
		Length:           length,
		Depth:            depth,
		SurfaceFootprint: width * length,
	}

	oi.dimensionsAsArray = []int{width, length, depth}
	sort.Ints(oi.dimensionsAsArray)

	return oi
}

// IsStable - Is this item stable (low centre of gravity), calculated as if the tipping point is >15 degrees.
// N.B. Assumes equal weight distribution.
func (oi *OrientatedItem) IsStable() bool {
	cacheKey := fmt.Sprintf("%d|%d|%d", oi.Width, oi.Length, oi.Depth)

	if val, exists := stabilityCache[cacheKey]; exists {
		return val
	}

	depth := oi.Depth
	if depth == 0 {
		depth = 1
	}

	minDim := oi.Length
	if oi.Width < minDim {
		minDim = oi.Width
	}

	stable := math.Atan(float64(minDim)/float64(depth)) > 0.261
	stabilityCache[cacheKey] = stable

	return stable
}

// IsSameDimensions - Is the supplied item the same size as this one?
func (oi *OrientatedItem) IsSameDimensions(item Item) bool {
	if item == oi.Item {
		return true
	}

	itemDimensions := []int{item.GetWidth(), item.GetLength(), item.GetDepth()}
	sort.Ints(itemDimensions)

	if len(itemDimensions) != len(oi.dimensionsAsArray) {
		return false
	}

	for i, dim := range itemDimensions {
		if dim != oi.dimensionsAsArray[i] {
			return false
		}
	}

	return true
}

// String returns string representation
func (oi *OrientatedItem) String() string {
	return fmt.Sprintf("%d|%d|%d", oi.Width, oi.Length, oi.Depth)
}
