package boxpacker

import "sort"

// Box packing (3D bin packing, knapsack problem).
//
// @author Doug Wright

// BoxList - List of boxes available to put items into, ordered by volume.
type BoxList struct {
	list     []Box
	isSorted bool
	sorter   BoxSorter
}

// NewBoxList creates a new BoxList with optional sorter
func NewBoxList(sorter ...BoxSorter) *BoxList {
	var s BoxSorter
	if len(sorter) > 0 {
		s = sorter[0]
	} else {
		s = &DefaultBoxSorter{}
	}
	return &BoxList{
		list:     make([]Box, 0),
		isSorted: false,
		sorter:   s,
	}
}

// FromArray - Do a bulk create.
func BoxListFromArray(boxes []Box, preSorted bool) *BoxList {
	list := &BoxList{
		list:     boxes,
		isSorted: preSorted,
		sorter:   &DefaultBoxSorter{},
	}
	return list
}

// GetIterator returns a sorted slice of boxes
func (bl *BoxList) GetIterator() []Box {
	if !bl.isSorted {
		sort.Slice(bl.list, func(i, j int) bool {
			return bl.sorter.Compare(bl.list[i], bl.list[j]) < 0
		})
		bl.isSorted = true
	}
	return bl.list
}

// Insert adds a box to the list
func (bl *BoxList) Insert(item Box) {
	bl.isSorted = false
	bl.list = append(bl.list, item)
}
