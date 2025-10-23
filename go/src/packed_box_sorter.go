package boxpacker

// Box packing (3D bin packing, knapsack problem).
//
// @author Doug Wright

// PackedBoxSorter - A callback to be used with sort, implementing logic to determine which PackedBox is "better".
type PackedBoxSorter interface {
	// Compare - Return -1 if boxA is "best", 1 if boxB is "best" or 0 if neither is "best".
	Compare(boxA *PackedBox, boxB *PackedBox) int
}
