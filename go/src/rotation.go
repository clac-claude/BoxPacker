package boxpacker

// Box packing (3D bin packing, knapsack problem).
//
// @author Doug Wright

// Rotation permutations
type Rotation int

const (
	// Never - Must be placed in it's defined orientation only
	Never Rotation = 1
	// KeepFlat - Can be turned sideways 90°, but cannot be placed *on* it's side e.g. fragile "↑this way up" items
	KeepFlat Rotation = 2
	// BestFit - No handling restrictions, item can be placed in any orientation
	BestFit Rotation = 6
)
