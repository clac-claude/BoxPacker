package boxpacker

import "sort"

// Box packing (3D bin packing, knapsack problem).
//
// @author Doug Wright

// WeightRedistributor - Redistributes weight across packed boxes
type WeightRedistributor struct {
	boxes                  *BoxList
	packedBoxSorter        PackedBoxSorter
	boxQuantitiesAvailable map[Box]int
	timeoutChecker         TimeoutChecker
}

// NewWeightRedistributor creates a new WeightRedistributor
func NewWeightRedistributor(
	boxes *BoxList,
	packedBoxSorter PackedBoxSorter,
	boxQuantitiesAvailable map[Box]int,
	timeoutChecker TimeoutChecker,
) *WeightRedistributor {
	return &WeightRedistributor{
		boxes:                  boxes,
		packedBoxSorter:        packedBoxSorter,
		boxQuantitiesAvailable: boxQuantitiesAvailable,
		timeoutChecker:         timeoutChecker,
	}
}

// RedistributeWeight - Given a solution set of packed boxes, repack them to achieve optimum weight distribution
func (wr *WeightRedistributor) RedistributeWeight(originalBoxes *PackedBoxList) *PackedBoxList {
	targetWeight := originalBoxes.GetMeanItemWeight()

	boxes := originalBoxes.GetIterator()

	sort.Slice(boxes, func(i, j int) bool {
		return boxes[j].GetWeight() < boxes[i].GetWeight()
	})

	iterationSuccessful := true
	for iterationSuccessful {
		iterationSuccessful = false

		for a := 0; a < len(boxes); a++ {
			for b := 0; b < len(boxes); b++ {
				if b <= a || boxes[a] == nil || boxes[b] == nil {
					continue
				}
				if boxes[a].GetWeight() == boxes[b].GetWeight() {
					continue // no need to evaluate
				}

				iterationSuccessful = wr.equaliseWeight(&boxes[a], &boxes[b], targetWeight)
				if iterationSuccessful {
					// remove any now-empty boxes from the list
					filtered := make([]*PackedBox, 0)
					for _, box := range boxes {
						if box != nil && box.Items.Count() > 0 {
							filtered = append(filtered, box)
						}
					}
					boxes = filtered
					break
				}
			}
			if iterationSuccessful {
				break
			}
		}
	}

	// Combine back into a single list
	packedBoxes := NewPackedBoxList(wr.packedBoxSorter)
	packedBoxes.InsertFromArray(boxes)

	return packedBoxes
}

// equaliseWeight attempts to equalise weight distribution between 2 boxes
// Returns true if the weight was rebalanced
func (wr *WeightRedistributor) equaliseWeight(boxA **PackedBox, boxB **PackedBox, targetWeight float64) bool {
	anyIterationSuccessful := false

	var overWeightBox, underWeightBox *PackedBox
	if (*boxA).GetWeight() > (*boxB).GetWeight() {
		overWeightBox = *boxA
		underWeightBox = *boxB
	} else {
		overWeightBox = *boxB
		underWeightBox = *boxA
	}

	overWeightBoxItems := overWeightBox.Items.AsItemArray()
	underWeightBoxItems := underWeightBox.Items.AsItemArray()

	for key, overWeightItem := range overWeightBoxItems {
		if wr.timeoutChecker != nil {
			if err := wr.timeoutChecker.ThrowOnTimeout(nil, ""); err != nil {
				break
			}
		}

		if !wouldRepackActuallyHelp(overWeightBoxItems, overWeightItem, underWeightBoxItems, targetWeight) {
			continue // moving this item would harm more than help
		}

		newItemsForUnderWeightBox := append(underWeightBoxItems, overWeightItem)
		newLighterBoxes := wr.doVolumeRepack(newItemsForUnderWeightBox, underWeightBox.Box)
		if newLighterBoxes.Count() != 1 {
			continue // only want to move this item if it still fits in a single box
		}

		underWeightBoxItems = append(underWeightBoxItems, overWeightItem)

		if len(overWeightBoxItems) == 1 { // sometimes a repack can be efficient enough to eliminate a box
			*boxB = newLighterBoxes.Top()
			*boxA = nil
			wr.boxQuantitiesAvailable[underWeightBox.Box]--
			wr.boxQuantitiesAvailable[overWeightBox.Box]++

			return true
		}

		// Remove item from overWeightBoxItems
		newOverWeightBoxItems := make([]Item, 0)
		for i, item := range overWeightBoxItems {
			if i != key {
				newOverWeightBoxItems = append(newOverWeightBoxItems, item)
			}
		}
		overWeightBoxItems = newOverWeightBoxItems

		newHeavierBoxes := wr.doVolumeRepack(overWeightBoxItems, overWeightBox.Box)
		if newHeavierBoxes.Count() != 1 {
			continue
		}

		wr.boxQuantitiesAvailable[overWeightBox.Box]++
		wr.boxQuantitiesAvailable[underWeightBox.Box]++
		wr.boxQuantitiesAvailable[newHeavierBoxes.Top().Box]--
		wr.boxQuantitiesAvailable[newLighterBoxes.Top().Box]--
		underWeightBox = newLighterBoxes.Top()
		*boxB = underWeightBox
		overWeightBox = newHeavierBoxes.Top()
		*boxA = overWeightBox

		anyIterationSuccessful = true
	}

	return anyIterationSuccessful
}

// doVolumeRepack does a volume repack of a set of items
func (wr *WeightRedistributor) doVolumeRepack(items []Item, currentBox Box) *PackedBoxList {
	packer := NewPacker(NewItemList(), wr.boxes, wr.packedBoxSorter)
	packer.ThrowOnUnpackableItem(false)

	// use the full set of boxes to allow smaller/larger for full efficiency
	for _, box := range wr.boxes.GetIterator() {
		packer.SetBoxQuantity(box, wr.boxQuantitiesAvailable[box])
	}
	packer.SetBoxQuantity(currentBox, wr.boxQuantitiesAvailable[currentBox]+1)

	itemList := NewItemList()
	for _, item := range items {
		itemList.Insert(item, 1)
	}
	packer.SetItems(itemList)

	return packer.DoBasicPacking(true)
}

// wouldRepackActuallyHelp - Not every attempted repack is actually helpful
func wouldRepackActuallyHelp(overWeightBoxItems []Item, overWeightItem Item, underWeightBoxItems []Item, targetWeight float64) bool {
	overWeightItemsWeight := 0
	for _, item := range overWeightBoxItems {
		overWeightItemsWeight += item.GetWeight()
	}

	underWeightItemsWeight := 0
	for _, item := range underWeightBoxItems {
		underWeightItemsWeight += item.GetWeight()
	}

	if float64(overWeightItem.GetWeight()+underWeightItemsWeight) > targetWeight {
		return false
	}

	oldVariance := calculateVariance(overWeightItemsWeight, underWeightItemsWeight)
	newVariance := calculateVariance(overWeightItemsWeight-overWeightItem.GetWeight(), underWeightItemsWeight+overWeightItem.GetWeight())

	return newVariance < oldVariance
}

// calculateVariance calculates variance between two box weights
func calculateVariance(boxAWeight, boxBWeight int) float64 {
	mean := float64(boxAWeight+boxBWeight) / 2.0
	diff := float64(boxAWeight) - mean
	return diff * diff // don't need to calculate B and ÷ 2, for a 2-item population the difference from mean is the same for each box
}
