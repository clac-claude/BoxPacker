package boxpacker

import "strings"

// Box packing (3D bin packing, knapsack problem).
//
// @author Doug Wright

// DefaultItemSorter - Default implementation of ItemSorter
type DefaultItemSorter struct{}

// Compare compares two items
func (dis *DefaultItemSorter) Compare(itemA Item, itemB Item) int {
	volumeA := itemA.GetWidth() * itemA.GetLength() * itemA.GetDepth()
	volumeB := itemB.GetWidth() * itemB.GetLength() * itemB.GetDepth()

	// larger volume first
	if volumeB < volumeA {
		return -1
	}
	if volumeB > volumeA {
		return 1
	}

	// heavier items first
	if itemB.GetWeight() < itemA.GetWeight() {
		return -1
	}
	if itemB.GetWeight() > itemA.GetWeight() {
		return 1
	}

	// alphabetical by description as final tiebreaker
	return strings.Compare(itemA.GetDescription(), itemB.GetDescription())
}
