package boxpacker

import "sort"

// Box packing (3D bin packing, knapsack problem).
//
// @author Doug Wright

// ItemList - List of items to be packed, ordered by volume.
type ItemList struct {
	list                 []Item
	isSorted             bool
	hasConstrainedItems  *bool
	hasNoRotationItems   *bool
	sorter               ItemSorter
}

// NewItemList creates a new ItemList with optional sorter
func NewItemList(sorter ...ItemSorter) *ItemList {
	var s ItemSorter
	if len(sorter) > 0 {
		s = sorter[0]
	} else {
		s = &DefaultItemSorter{}
	}
	return &ItemList{
		list:                make([]Item, 0),
		isSorted:            false,
		hasConstrainedItems: nil,
		hasNoRotationItems:  nil,
		sorter:              s,
	}
}

// ItemListFromArray - Do a bulk create.
func ItemListFromArray(items []Item, preSorted bool) *ItemList {
	// internal sort is largest at the end
	reversed := make([]Item, len(items))
	for i, item := range items {
		reversed[len(items)-1-i] = item
	}

	list := &ItemList{
		list:                reversed,
		isSorted:            preSorted,
		hasConstrainedItems: nil,
		hasNoRotationItems:  nil,
		sorter:              &DefaultItemSorter{},
	}
	return list
}

// Insert adds items to the list
func (il *ItemList) Insert(item Item, qty int) {
	if qty <= 0 {
		qty = 1
	}
	for i := 0; i < qty; i++ {
		il.list = append(il.list, item)
	}
	il.isSorted = false

	// normally lazy evaluated, override if that's already been done
	if il.hasConstrainedItems != nil {
		_, isConstrained := item.(ConstrainedPlacementItem)
		val := *il.hasConstrainedItems || isConstrained
		il.hasConstrainedItems = &val
	}

	if il.hasNoRotationItems != nil {
		val := *il.hasNoRotationItems || item.GetAllowedRotation() == Never
		il.hasNoRotationItems = &val
	}
}

// Remove removes an item from the list
func (il *ItemList) Remove(item Item) {
	if !il.isSorted {
		il.sort()
	}

	for i := len(il.list) - 1; i >= 0; i-- {
		if il.list[i] == item {
			il.list = append(il.list[:i], il.list[i+1:]...)
			return
		}
	}
}

// RemovePackedItems removes packed items from the list
func (il *ItemList) RemovePackedItems(packedItemList *PackedItemList) {
	for _, packedItem := range packedItemList.GetIterator() {
		for i := len(il.list) - 1; i >= 0; i-- {
			if il.list[i] == packedItem.Item {
				il.list = append(il.list[:i], il.list[i+1:]...)
				break
			}
		}
	}
}

// Extract removes and returns the top item
func (il *ItemList) Extract() Item {
	if !il.isSorted {
		il.sort()
	}

	if len(il.list) == 0 {
		return nil
	}

	item := il.list[len(il.list)-1]
	il.list = il.list[:len(il.list)-1]
	return item
}

// Top returns the top item without removing it
func (il *ItemList) Top() Item {
	if !il.isSorted {
		il.sort()
	}

	if len(il.list) == 0 {
		return nil
	}

	return il.list[len(il.list)-1]
}

// TopN returns the top N items
func (il *ItemList) TopN(n int) *ItemList {
	if !il.isSorted {
		il.sort()
	}

	if n > len(il.list) {
		n = len(il.list)
	}

	topNList := NewItemList(il.sorter)
	if n > 0 {
		topNList.list = make([]Item, n)
		copy(topNList.list, il.list[len(il.list)-n:])
	}
	topNList.isSorted = true

	return topNList
}

// GetIterator returns a sorted slice of items
func (il *ItemList) GetIterator() []Item {
	if !il.isSorted {
		il.sort()
	}

	// return reversed (sorted largest first for iteration)
	reversed := make([]Item, len(il.list))
	for i, item := range il.list {
		reversed[len(il.list)-1-i] = item
	}
	return reversed
}

// Count returns the number of items in the list
func (il *ItemList) Count() int {
	return len(il.list)
}

// HasConstrainedItems checks if this list contains items with constrained placement criteria
func (il *ItemList) HasConstrainedItems() bool {
	if il.hasConstrainedItems == nil {
		hasConstrained := false
		for _, item := range il.list {
			if _, ok := item.(ConstrainedPlacementItem); ok {
				hasConstrained = true
				break
			}
		}
		il.hasConstrainedItems = &hasConstrained
	}

	return *il.hasConstrainedItems
}

// HasNoRotationItems checks if this list contains items which cannot be rotated
func (il *ItemList) HasNoRotationItems() bool {
	if il.hasNoRotationItems == nil {
		hasNoRotation := false
		for _, item := range il.list {
			if item.GetAllowedRotation() == Never {
				hasNoRotation = true
				break
			}
		}
		il.hasNoRotationItems = &hasNoRotation
	}

	return *il.hasNoRotationItems
}

// sort is internal helper to sort the list
func (il *ItemList) sort() {
	sort.Slice(il.list, func(i, j int) bool {
		return il.sorter.Compare(il.list[i], il.list[j]) < 0
	})
	// internal sort is largest at the end
	for i, j := 0, len(il.list)-1; i < j; i, j = i+1, j-1 {
		il.list[i], il.list[j] = il.list[j], il.list[i]
	}
	il.isSorted = true
}
