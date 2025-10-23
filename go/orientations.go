package main

// GenerateOrientations generates all possible orientations for an item
func GenerateOrientations(item *Item, prevItem *OrientatedItem) []*OrientatedItem {
	// If same as previous item, keep the same orientation
	if prevItem != nil && IsSameDimensions(item, prevItem.Item) {
		return []*OrientatedItem{
			NewOrientatedItem(item, prevItem.Width, prevItem.Length, prevItem.Depth),
		}
	}

	orientations := make([]*OrientatedItem, 0, 6)
	w, l, d := item.Width, item.Length, item.Depth

	// Base orientation
	orientations = append(orientations, NewOrientatedItem(item, w, l, d))

	// Simple 2D rotation
	if item.Rotation != RotationNever {
		orientations = append(orientations, NewOrientatedItem(item, l, w, d))
	}

	// Full 3D rotation if allowed
	if item.Rotation == RotationBestFit {
		orientations = append(orientations,
			NewOrientatedItem(item, w, d, l),
			NewOrientatedItem(item, l, d, w),
			NewOrientatedItem(item, d, w, l),
			NewOrientatedItem(item, d, l, w),
		)
	}

	return orientations
}

// GetPossibleOrientations returns orientations that fit in the given space
func GetPossibleOrientations(item *Item, prevItem *OrientatedItem, widthLeft, lengthLeft, depthLeft int, packedWeight, boxMaxWeight int) []*OrientatedItem {
	// Check if item is too heavy
	if item.Weight > (boxMaxWeight - packedWeight) {
		return nil
	}

	allOrientations := GenerateOrientations(item, prevItem)
	possible := make([]*OrientatedItem, 0, len(allOrientations))

	for _, orientation := range allOrientations {
		if orientation.Width <= widthLeft &&
		   orientation.Length <= lengthLeft &&
		   orientation.Depth <= depthLeft {
			possible = append(possible, orientation)
		}
	}

	return possible
}

// GetUsableOrientations filters orientations by stability
func GetUsableOrientations(orientations []*OrientatedItem, boxDepth int) []*OrientatedItem {
	stable := make([]*OrientatedItem, 0, len(orientations))
	unstable := make([]*OrientatedItem, 0, len(orientations))

	for _, orientation := range orientations {
		if orientation.IsStable() || boxDepth == orientation.Depth {
			stable = append(stable, orientation)
		} else {
			unstable = append(unstable, orientation)
		}
	}

	if len(stable) > 0 {
		return stable
	}

	return unstable
}

// IsSameDimensions checks if two items have the same dimensions
func IsSameDimensions(itemA, itemB *Item) bool {
	if itemA == itemB {
		return true
	}

	// Sort dimensions to compare
	aDims := []int{itemA.Width, itemA.Length, itemA.Depth}
	bDims := []int{itemB.Width, itemB.Length, itemB.Depth}

	sortInts(aDims)
	sortInts(bDims)

	return aDims[0] == bDims[0] && aDims[1] == bDims[1] && aDims[2] == bDims[2]
}

// Simple bubble sort for 3 elements
func sortInts(arr []int) {
	for i := 0; i < len(arr); i++ {
		for j := i + 1; j < len(arr); j++ {
			if arr[i] > arr[j] {
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}
}

// Max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
