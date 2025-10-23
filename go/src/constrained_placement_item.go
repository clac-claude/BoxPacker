package boxpacker

// Box packing (3D bin packing, knapsack problem).
//
// @author Doug Wright

// ConstrainedPlacementItem - An item to be packed where additional constraints need to be considered.
// Only implement this interface if you actually need this additional functionality as it will slow down
// the packing algorithm.
type ConstrainedPlacementItem interface {
	Item

	// CanBePacked - Hook for user implementation of item-specific constraints, e.g. max <x> batteries per box.
	CanBePacked(
		packedBox *PackedBox,
		proposedX int,
		proposedY int,
		proposedZ int,
		width int,
		length int,
		depth int,
	) bool
}
