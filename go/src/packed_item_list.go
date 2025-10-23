package boxpacker

import "sort"

// Box packing (3D bin packing, knapsack problem).
//
// @author Doug Wright

// PackedItemList - List of packed items, ordered by volume.
type PackedItemList struct {
	list     []*PackedItem
	weight   int
	volume   int
	isSorted bool
}

// NewPackedItemList creates a new PackedItemList
func NewPackedItemList() *PackedItemList {
	return &PackedItemList{
		list:     make([]*PackedItem, 0),
		weight:   0,
		volume:   0,
		isSorted: false,
	}
}

// Insert adds a packed item to the list
func (pil *PackedItemList) Insert(item *PackedItem) {
	pil.list = append(pil.list, item)
	pil.weight += item.Item.GetWeight()
	pil.volume += item.Width * item.Length * item.Depth
	pil.isSorted = false
}

// GetIterator returns a sorted slice of packed items
func (pil *PackedItemList) GetIterator() []*PackedItem {
	if !pil.isSorted {
		sort.Slice(pil.list, func(i, j int) bool {
			return pil.compare(pil.list[i], pil.list[j]) < 0
		})
		pil.isSorted = true
	}

	return pil.list
}

// Count returns the number of items in the list
func (pil *PackedItemList) Count() int {
	return len(pil.list)
}

// AsItemArray returns a copy of this list as a standard array of Items
func (pil *PackedItemList) AsItemArray() []Item {
	items := make([]Item, len(pil.list))
	for i, packedItem := range pil.list {
		items[i] = packedItem.Item
	}
	return items
}

// GetVolume returns the total volume of these items
func (pil *PackedItemList) GetVolume() int {
	return pil.volume
}

// GetWeight returns the total weight of these items
func (pil *PackedItemList) GetWeight() int {
	return pil.weight
}

// compare is internal comparator
func (pil *PackedItemList) compare(itemA, itemB *PackedItem) int {
	itemAVolume := itemA.Item.GetWidth() * itemA.Item.GetLength() * itemA.Item.GetDepth()
	itemBVolume := itemB.Item.GetWidth() * itemB.Item.GetLength() * itemB.Item.GetDepth()

	if itemBVolume != itemAVolume {
		if itemBVolume > itemAVolume {
			return -1
		}
		return 1
	}

	if itemB.Item.GetWeight() != itemA.Item.GetWeight() {
		if itemB.Item.GetWeight() > itemA.Item.GetWeight() {
			return -1
		}
		return 1
	}

	return 0
}
