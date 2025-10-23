package boxpacker

import (
	"math"
	"sort"
)

// Box packing (3D bin packing, knapsack problem).
//
// @author Doug Wright

// PackedBoxList - List of packed boxes.
type PackedBoxList struct {
	list     []*PackedBox
	isSorted bool
	sorter   PackedBoxSorter
}

// NewPackedBoxList creates a new PackedBoxList with optional sorter
func NewPackedBoxList(sorter ...PackedBoxSorter) *PackedBoxList {
	var s PackedBoxSorter
	if len(sorter) > 0 {
		s = sorter[0]
	} else {
		s = &DefaultPackedBoxSorter{}
	}
	return &PackedBoxList{
		list:     make([]*PackedBox, 0),
		isSorted: false,
		sorter:   s,
	}
}

// GetIterator returns a sorted slice of packed boxes
func (pbl *PackedBoxList) GetIterator() []*PackedBox {
	if !pbl.isSorted {
		sort.Slice(pbl.list, func(i, j int) bool {
			return pbl.sorter.Compare(pbl.list[i], pbl.list[j]) < 0
		})
		pbl.isSorted = true
	}

	return pbl.list
}

// Count returns the number of boxes in the list
func (pbl *PackedBoxList) Count() int {
	return len(pbl.list)
}

// Insert adds a packed box to the list
func (pbl *PackedBoxList) Insert(item *PackedBox) {
	pbl.list = append(pbl.list, item)
	pbl.isSorted = false
}

// InsertFromArray does a bulk insert
func (pbl *PackedBoxList) InsertFromArray(boxes []*PackedBox) {
	for _, box := range boxes {
		pbl.Insert(box)
	}
}

// Top returns the top (first) packed box
func (pbl *PackedBoxList) Top() *PackedBox {
	if !pbl.isSorted {
		sort.Slice(pbl.list, func(i, j int) bool {
			return pbl.sorter.Compare(pbl.list[i], pbl.list[j]) < 0
		})
		pbl.isSorted = true
	}

	if len(pbl.list) > 0 {
		return pbl.list[0]
	}
	return nil
}

// GetMeanWeight calculates the average (mean) weight of the boxes
func (pbl *PackedBoxList) GetMeanWeight() float64 {
	if len(pbl.list) == 0 {
		return 0
	}

	meanWeight := 0.0

	for _, box := range pbl.list {
		meanWeight += float64(box.GetWeight())
	}

	return meanWeight / float64(len(pbl.list))
}

// GetMeanItemWeight calculates the average (mean) weight of the items in the boxes
func (pbl *PackedBoxList) GetMeanItemWeight() float64 {
	if len(pbl.list) == 0 {
		return 0
	}

	meanWeight := 0.0

	for _, box := range pbl.list {
		meanWeight += float64(box.GetItemWeight())
	}

	return meanWeight / float64(len(pbl.list))
}

// GetWeightVariance calculates the variance in weight between these boxes
func (pbl *PackedBoxList) GetWeightVariance() float64 {
	if len(pbl.list) == 0 {
		return 0
	}

	mean := pbl.GetMeanWeight()

	weightVariance := 0.0
	for _, box := range pbl.list {
		diff := float64(box.GetWeight()) - mean
		weightVariance += diff * diff
	}

	return math.Round(weightVariance/float64(len(pbl.list))*10) / 10
}

// GetVolumeUtilisation returns the volume utilisation of the set of packed boxes
func (pbl *PackedBoxList) GetVolumeUtilisation() float64 {
	itemVolume := 0
	boxVolume := 0

	for _, box := range pbl.list {
		boxVolume += box.GetInnerVolume()

		for _, item := range box.Items.list {
			itemVolume += (item.Item.GetWidth() * item.Item.GetLength() * item.Item.GetDepth())
		}
	}

	if boxVolume == 0 {
		return 0
	}

	return math.Round(float64(itemVolume)/float64(boxVolume)*1000) / 10
}
