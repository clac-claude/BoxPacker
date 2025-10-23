package boxpacker

// Box packing (3D bin packing, knapsack problem).
//
// @author Doug Wright

// LimitedSupplyBox - A "box" (or envelope?) to pack items into with limited supply.
type LimitedSupplyBox interface {
	Box

	// GetQuantityAvailable - Quantity of boxes available.
	GetQuantityAvailable() int
}
