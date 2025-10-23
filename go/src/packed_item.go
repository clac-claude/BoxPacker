package boxpacker

// Box packing (3D bin packing, knapsack problem).
//
// @author Doug Wright

// PackedItem - A packed item.
type PackedItem struct {
	Item   Item
	X      int
	Y      int
	Z      int
	Width  int
	Length int
	Depth  int
	Volume int
}

// NewPackedItem creates a new PackedItem
func NewPackedItem(item Item, x, y, z, width, length, depth int) *PackedItem {
	return &PackedItem{
		Item:   item,
		X:      x,
		Y:      y,
		Z:      z,
		Width:  width,
		Length: length,
		Depth:  depth,
		Volume: width * length * depth,
	}
}

// FromOrientatedItem creates a PackedItem from an OrientatedItem
func PackedItemFromOrientatedItem(orientatedItem *OrientatedItem, x, y, z int) *PackedItem {
	return NewPackedItem(
		orientatedItem.Item,
		x,
		y,
		z,
		orientatedItem.Width,
		orientatedItem.Length,
		orientatedItem.Depth,
	)
}
