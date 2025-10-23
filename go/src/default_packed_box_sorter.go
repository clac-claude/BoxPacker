package boxpacker

// Box packing (3D bin packing, knapsack problem).
//
// @author Doug Wright

// DefaultPackedBoxSorter - Default implementation of PackedBoxSorter
type DefaultPackedBoxSorter struct{}

// Compare compares two packed boxes
func (dpbs *DefaultPackedBoxSorter) Compare(boxA *PackedBox, boxB *PackedBox) int {
	// prefer boxes with more items
	countA := boxA.Items.Count()
	countB := boxB.Items.Count()

	if countB < countA {
		return -1
	}
	if countB > countA {
		return 1
	}

	// prefer boxes with higher volume utilisation
	utilA := boxA.GetVolumeUtilisation()
	utilB := boxB.GetVolumeUtilisation()

	if utilB < utilA {
		return -1
	}
	if utilB > utilA {
		return 1
	}

	// prefer boxes with larger used volume
	usedVolA := boxA.GetUsedVolume()
	usedVolB := boxB.GetUsedVolume()

	if usedVolB < usedVolA {
		return -1
	}
	if usedVolB > usedVolA {
		return 1
	}

	return 0
}
