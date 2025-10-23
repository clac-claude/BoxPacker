package main

// Rotation represents allowed rotation types
type Rotation int

const (
	RotationNever   Rotation = 0
	RotationKeepFlat Rotation = 1
	RotationBestFit  Rotation = 2
)

// Item represents an item to be packed
type Item struct {
	Width    int
	Length   int
	Depth    int
	Weight   int
	Rotation Rotation
}

// Box represents a box/container
type Box struct {
	InnerWidth  int
	InnerLength int
	InnerDepth  int
	MaxWeight   int
}

// OrientatedItem represents an item in a specific orientation
type OrientatedItem struct {
	Item            *Item
	Width           int
	Length          int
	Depth           int
	SurfaceFootprint int
}

// NewOrientatedItem creates a new orientated item
func NewOrientatedItem(item *Item, width, length, depth int) *OrientatedItem {
	return &OrientatedItem{
		Item:            item,
		Width:           width,
		Length:          length,
		Depth:           depth,
		SurfaceFootprint: width * length,
	}
}

// IsStable checks if the orientation is stable (low center of gravity)
func (o *OrientatedItem) IsStable() bool {
	// An item is stable if it's resting on its largest face
	return o.Depth <= o.Width && o.Depth <= o.Length
}

// PackedItem represents a packed item with position
type PackedItem struct {
	Item   *OrientatedItem
	X      int
	Y      int
	Z      int
	Width  int
	Length int
	Depth  int
}

// NewPackedItem creates a new packed item
func NewPackedItem(item *OrientatedItem, x, y, z int) *PackedItem {
	return &PackedItem{
		Item:   item,
		X:      x,
		Y:      y,
		Z:      z,
		Width:  item.Width,
		Length: item.Length,
		Depth:  item.Depth,
	}
}

// PackedItemList is a slice of packed items
type PackedItemList []*PackedItem

// GetWeight calculates total weight of packed items
func (p PackedItemList) GetWeight() int {
	total := 0
	for _, item := range p {
		total += item.Item.Item.Weight
	}
	return total
}
