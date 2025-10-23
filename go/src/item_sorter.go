package boxpacker

// Box packing (3D bin packing, knapsack problem).
//
// @author Doug Wright

// ItemSorter - A callback to be used with sort, implementing logic to determine which Item is a higher priority for packing.
type ItemSorter interface {
	// Compare - Return -1 if itemA is preferred, 1 if itemB is preferred or 0 if neither is preferred.
	Compare(itemA Item, itemB Item) int
}
