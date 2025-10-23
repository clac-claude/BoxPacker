package boxpacker

// Box packing (3D bin packing, knapsack problem).
//
// @author Doug Wright

// BoxSorter - A callback to be used with sort, implementing logic to determine which Box is "better".
type BoxSorter interface {
	// Compare - Return -1 if boxA is "best", 1 if boxB is "best" or 0 if neither is "best".
	Compare(boxA Box, boxB Box) int
}
