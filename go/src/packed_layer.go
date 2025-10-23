package boxpacker

import "math"

// Box packing (3D bin packing, knapsack problem).
//
// @author Doug Wright

// PackedLayer - A packed layer.
type PackedLayer struct {
	items []*PackedItem
}

// NewPackedLayer creates a new PackedLayer
func NewPackedLayer() *PackedLayer {
	return &PackedLayer{
		items: make([]*PackedItem, 0),
	}
}

// Insert adds a packed item to this layer
func (pl *PackedLayer) Insert(packedItem *PackedItem) {
	pl.items = append(pl.items, packedItem)
}

// GetItems returns the packed items
func (pl *PackedLayer) GetItems() []*PackedItem {
	return pl.items
}

// GetFootprint calculates footprint area of this layer (in mm^2)
func (pl *PackedLayer) GetFootprint() int {
	return pl.GetWidth() * pl.GetLength()
}

// GetStartX returns the minimum X coordinate
func (pl *PackedLayer) GetStartX() int {
	if len(pl.items) == 0 {
		return 0
	}

	min := math.MaxInt32
	for _, item := range pl.items {
		if item.X < min {
			min = item.X
		}
	}
	return min
}

// GetEndX returns the maximum X coordinate
func (pl *PackedLayer) GetEndX() int {
	if len(pl.items) == 0 {
		return 0
	}

	max := 0
	for _, item := range pl.items {
		end := item.X + item.Width
		if end > max {
			max = end
		}
	}
	return max
}

// GetWidth returns the width of this layer
func (pl *PackedLayer) GetWidth() int {
	if len(pl.items) == 0 {
		return 0
	}

	return pl.GetEndX() - pl.GetStartX()
}

// GetStartY returns the minimum Y coordinate
func (pl *PackedLayer) GetStartY() int {
	if len(pl.items) == 0 {
		return 0
	}

	min := math.MaxInt32
	for _, item := range pl.items {
		if item.Y < min {
			min = item.Y
		}
	}
	return min
}

// GetEndY returns the maximum Y coordinate
func (pl *PackedLayer) GetEndY() int {
	if len(pl.items) == 0 {
		return 0
	}

	max := 0
	for _, item := range pl.items {
		end := item.Y + item.Length
		if end > max {
			max = end
		}
	}
	return max
}

// GetLength returns the length of this layer
func (pl *PackedLayer) GetLength() int {
	if len(pl.items) == 0 {
		return 0
	}

	return pl.GetEndY() - pl.GetStartY()
}

// GetStartZ returns the minimum Z coordinate
func (pl *PackedLayer) GetStartZ() int {
	if len(pl.items) == 0 {
		return 0
	}

	min := math.MaxInt32
	for _, item := range pl.items {
		if item.Z < min {
			min = item.Z
		}
	}
	return min
}

// GetEndZ returns the maximum Z coordinate
func (pl *PackedLayer) GetEndZ() int {
	if len(pl.items) == 0 {
		return 0
	}

	max := 0
	for _, item := range pl.items {
		end := item.Z + item.Depth
		if end > max {
			max = end
		}
	}
	return max
}

// GetDepth returns the depth of this layer
func (pl *PackedLayer) GetDepth() int {
	if len(pl.items) == 0 {
		return 0
	}

	return pl.GetEndZ() - pl.GetStartZ()
}

// GetWeight returns the total weight of this layer
func (pl *PackedLayer) GetWeight() int {
	weight := 0
	for _, item := range pl.items {
		weight += item.Item.GetWeight()
	}
	return weight
}

// Merge merges another layer into this one
func (pl *PackedLayer) Merge(otherLayer *PackedLayer) {
	pl.items = append(pl.items, otherLayer.items...)
}
