package boxpacker

import "fmt"

// Box packing (3D bin packing, knapsack problem).
//
// @author Doug Wright

// OrientatedItemFactory - Figure out orientations for an item and a given set of dimensions.
type OrientatedItemFactory struct {
	box                                Box
	singlePassMode                     bool
	boxIsRotated                       bool
	emptyBoxStableItemOrientationCache map[string]bool
}

// NewOrientatedItemFactory creates a new OrientatedItemFactory
func NewOrientatedItemFactory(box Box) *OrientatedItemFactory {
	return &OrientatedItemFactory{
		box:                                box,
		singlePassMode:                     false,
		boxIsRotated:                       false,
		emptyBoxStableItemOrientationCache: make(map[string]bool),
	}
}

// SetSinglePassMode sets single pass mode
func (oif *OrientatedItemFactory) SetSinglePassMode(singlePassMode bool) {
	oif.singlePassMode = singlePassMode
}

// SetBoxIsRotated sets whether the box is rotated
func (oif *OrientatedItemFactory) SetBoxIsRotated(boxIsRotated bool) {
	oif.boxIsRotated = boxIsRotated
}

// GetBestOrientation gets the best orientation for an item
func (oif *OrientatedItemFactory) GetBestOrientation(
	item Item,
	prevItem *OrientatedItem,
	nextItems *ItemList,
	widthLeft int,
	lengthLeft int,
	depthLeft int,
	rowLength int,
	x int,
	y int,
	z int,
	prevPackedItemList *PackedItemList,
	considerStability bool,
) *OrientatedItem {
	possibleOrientations := oif.GetPossibleOrientations(item, prevItem, widthLeft, lengthLeft, depthLeft, x, y, z, prevPackedItemList)

	var usableOrientations []*OrientatedItem
	if considerStability {
		usableOrientations = oif.getUsableOrientations(item, possibleOrientations)
	} else {
		usableOrientations = possibleOrientations
	}

	if len(usableOrientations) == 0 {
		return nil
	}

	// Return first orientation (simplified - full version would use OrientatedItemSorter)
	return usableOrientations[0]
}

// GetPossibleOrientations finds all possible orientations for an item
func (oif *OrientatedItemFactory) GetPossibleOrientations(
	item Item,
	prevItem *OrientatedItem,
	widthLeft int,
	lengthLeft int,
	depthLeft int,
	x int,
	y int,
	z int,
	prevPackedItemList *PackedItemList,
) []*OrientatedItem {
	permutations := oif.generatePermutations(item, prevItem)

	// remove any that simply don't fit
	orientations := make([]*OrientatedItem, 0)
	for _, dimensions := range permutations {
		if dimensions[0] <= widthLeft && dimensions[1] <= lengthLeft && dimensions[2] <= depthLeft {
			orientations = append(orientations, NewOrientatedItem(item, dimensions[0], dimensions[1], dimensions[2]))
		}
	}

	// Handle ConstrainedPlacementItem
	if constrainedItem, ok := item.(ConstrainedPlacementItem); ok {
		if _, isWorkingVolume := oif.box.(*WorkingVolume); !isWorkingVolume {
			filtered := make([]*OrientatedItem, 0)
			for _, oi := range orientations {
				var canBePacked bool
				if oif.boxIsRotated {
					// Create rotated packed item list (simplified)
					canBePacked = constrainedItem.CanBePacked(NewPackedBox(oif.box, prevPackedItemList), y, x, z, oi.Length, oi.Width, oi.Depth)
				} else {
					canBePacked = constrainedItem.CanBePacked(NewPackedBox(oif.box, prevPackedItemList), x, y, z, oi.Width, oi.Length, oi.Depth)
				}
				if canBePacked {
					filtered = append(filtered, oi)
				}
			}
			orientations = filtered
		}
	}

	return orientations
}

// getUsableOrientations filters orientations by stability
func (oif *OrientatedItemFactory) getUsableOrientations(item Item, possibleOrientations []*OrientatedItem) []*OrientatedItem {
	stableOrientations := make([]*OrientatedItem, 0)
	unstableOrientations := make([]*OrientatedItem, 0)

	// Divide possible orientations into stable and unstable
	for _, orientation := range possibleOrientations {
		if orientation.IsStable() || oif.box.GetInnerDepth() == orientation.Depth {
			stableOrientations = append(stableOrientations, orientation)
		} else {
			unstableOrientations = append(unstableOrientations, orientation)
		}
	}

	// Prefer stable orientations
	if len(stableOrientations) > 0 {
		return stableOrientations
	}

	if len(unstableOrientations) > 0 && !oif.hasStableOrientationsInEmptyBox(item) {
		return unstableOrientations
	}

	return make([]*OrientatedItem, 0)
}

// hasStableOrientationsInEmptyBox checks if item has stable orientations in empty box
func (oif *OrientatedItemFactory) hasStableOrientationsInEmptyBox(item Item) bool {
	cacheKey := fmt.Sprintf("%d|%d|%d|%d|%d|%d|%d",
		item.GetWidth(), item.GetLength(), item.GetDepth(), item.GetAllowedRotation(),
		oif.box.GetInnerWidth(), oif.box.GetInnerLength(), oif.box.GetInnerDepth())

	if cached, exists := oif.emptyBoxStableItemOrientationCache[cacheKey]; exists {
		return cached
	}

	orientations := oif.GetPossibleOrientations(
		item,
		nil,
		oif.box.GetInnerWidth(),
		oif.box.GetInnerLength(),
		oif.box.GetInnerDepth(),
		0, 0, 0,
		NewPackedItemList(),
	)

	hasStable := false
	for _, orientation := range orientations {
		if orientation.IsStable() {
			hasStable = true
			break
		}
	}

	oif.emptyBoxStableItemOrientationCache[cacheKey] = hasStable
	return hasStable
}

// generatePermutations generates all possible dimension permutations
func (oif *OrientatedItemFactory) generatePermutations(item Item, prevItem *OrientatedItem) [][]int {
	// Special case items that are the same as what we just packed - keep orientation
	if prevItem != nil && prevItem.IsSameDimensions(item) {
		return [][]int{{prevItem.Width, prevItem.Length, prevItem.Depth}}
	}

	permutationsMap := make(map[string][]int)
	w := item.GetWidth()
	l := item.GetLength()
	d := item.GetDepth()

	key := fmt.Sprintf("%d|%d|%d", w, l, d)
	permutationsMap[key] = []int{w, l, d}

	if item.GetAllowedRotation() != Never { // simple 2D rotation
		key = fmt.Sprintf("%d|%d|%d", l, w, d)
		permutationsMap[key] = []int{l, w, d}
	}

	if item.GetAllowedRotation() == BestFit { // add 3D rotation if we're allowed
		permutationsMap[fmt.Sprintf("%d|%d|%d", w, d, l)] = []int{w, d, l}
		permutationsMap[fmt.Sprintf("%d|%d|%d", l, d, w)] = []int{l, d, w}
		permutationsMap[fmt.Sprintf("%d|%d|%d", d, w, l)] = []int{d, w, l}
		permutationsMap[fmt.Sprintf("%d|%d|%d", d, l, w)] = []int{d, l, w}
	}

	permutations := make([][]int, 0, len(permutationsMap))
	for _, perm := range permutationsMap {
		permutations = append(permutations, perm)
	}

	return permutations
}
