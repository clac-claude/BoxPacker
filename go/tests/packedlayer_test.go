package tests

import "testing"

// PackedLayer represents a layer of packed items
type PackedLayer struct {
	items []*PackedItem
}

// NewPackedLayer creates a new PackedLayer
func NewPackedLayer() *PackedLayer {
	return &PackedLayer{
		items: make([]*PackedItem, 0),
	}
}

func (l *PackedLayer) Insert(item *PackedItem) {
	l.items = append(l.items, item)
}

func (l *PackedLayer) GetStartX() int {
	if len(l.items) == 0 {
		return 0
	}
	minX := l.items[0].X
	for _, item := range l.items {
		if item.X < minX {
			minX = item.X
		}
	}
	return minX
}

func (l *PackedLayer) GetEndX() int {
	if len(l.items) == 0 {
		return 0
	}
	maxX := l.items[0].X + l.items[0].Width
	for _, item := range l.items {
		endX := item.X + item.Width
		if endX > maxX {
			maxX = endX
		}
	}
	return maxX
}

func (l *PackedLayer) GetWidth() int {
	return l.GetEndX() - l.GetStartX()
}

func (l *PackedLayer) GetStartY() int {
	if len(l.items) == 0 {
		return 0
	}
	minY := l.items[0].Y
	for _, item := range l.items {
		if item.Y < minY {
			minY = item.Y
		}
	}
	return minY
}

func (l *PackedLayer) GetEndY() int {
	if len(l.items) == 0 {
		return 0
	}
	maxY := l.items[0].Y + l.items[0].Length
	for _, item := range l.items {
		endY := item.Y + item.Length
		if endY > maxY {
			maxY = endY
		}
	}
	return maxY
}

func (l *PackedLayer) GetLength() int {
	return l.GetEndY() - l.GetStartY()
}

func (l *PackedLayer) GetStartZ() int {
	if len(l.items) == 0 {
		return 0
	}
	minZ := l.items[0].Z
	for _, item := range l.items {
		if item.Z < minZ {
			minZ = item.Z
		}
	}
	return minZ
}

func (l *PackedLayer) GetEndZ() int {
	if len(l.items) == 0 {
		return 0
	}
	maxZ := l.items[0].Z + l.items[0].Depth
	for _, item := range l.items {
		endZ := item.Z + item.Depth
		if endZ > maxZ {
			maxZ = endZ
		}
	}
	return maxZ
}

func (l *PackedLayer) GetDepth() int {
	return l.GetEndZ() - l.GetStartZ()
}

func (l *PackedLayer) GetFootprint() int {
	return l.GetWidth() * l.GetLength()
}

func (l *PackedLayer) GetWeight() int {
	totalWeight := 0
	for _, item := range l.items {
		if testItem, ok := item.Item.(*TestItem); ok {
			totalWeight += testItem.Weight
		}
	}
	return totalWeight
}

func TestGetters(t *testing.T) {
	packedItem := NewPackedItem(NewTestItem("Item", 11, 22, 33, 43, BestFit), 4, 5, 6, 33, 11, 22)
	packedLayer := NewPackedLayer()
	packedLayer.Insert(packedItem)

	tests := []struct {
		name     string
		got      int
		expected int
	}{
		{"StartX", packedLayer.GetStartX(), 4},
		{"EndX", packedLayer.GetEndX(), 37},
		{"Width", packedLayer.GetWidth(), 33},
		{"StartY", packedLayer.GetStartY(), 5},
		{"EndY", packedLayer.GetEndY(), 16},
		{"Length", packedLayer.GetLength(), 11},
		{"StartZ", packedLayer.GetStartZ(), 6},
		{"EndZ", packedLayer.GetEndZ(), 28},
		{"Depth", packedLayer.GetDepth(), 22},
		{"Footprint", packedLayer.GetFootprint(), 363},
		{"Weight", packedLayer.GetWeight(), 43},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s: expected %d, got %d", tt.name, tt.expected, tt.got)
			}
		})
	}
}
