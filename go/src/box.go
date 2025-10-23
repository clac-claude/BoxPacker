package boxpacker

// Box packing (3D bin packing, knapsack problem).
//
// @author Doug Wright

// Box - A "box" (or envelope?) to pack items into.
type Box interface {
	// GetReference - Reference for box type (e.g. SKU or description).
	GetReference() string

	// GetOuterWidth - Outer width in mm.
	GetOuterWidth() int

	// GetOuterLength - Outer length in mm.
	GetOuterLength() int

	// GetOuterDepth - Outer depth in mm.
	GetOuterDepth() int

	// GetEmptyWeight - Empty weight in g.
	GetEmptyWeight() int

	// GetInnerWidth - Inner width in mm.
	GetInnerWidth() int

	// GetInnerLength - Inner length in mm.
	GetInnerLength() int

	// GetInnerDepth - Inner depth in mm.
	GetInnerDepth() int

	// GetMaxWeight - Max weight the packaging can hold in g.
	GetMaxWeight() int
}
