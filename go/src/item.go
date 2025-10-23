package boxpacker

// Box packing (3D bin packing, knapsack problem).
//
// @author Doug Wright

// Item - An item to be packed.
type Item interface {
	// GetDescription - Item SKU etc.
	GetDescription() string

	// GetWidth - Item width in mm.
	GetWidth() int

	// GetLength - Item length in mm.
	GetLength() int

	// GetDepth - Item depth in mm.
	GetDepth() int

	// GetWeight - Item weight in g.
	GetWeight() int

	// GetAllowedRotation - Possible item rotations allowed.
	GetAllowedRotation() Rotation
}
