package boxpacker

import (
	"github.com/dvdoug/boxpacker/go/src/exception"
	"math"
	"sort"
)

// Box packing (3D bin packing, knapsack problem).
//
// @author Doug Wright

// Packer - Actual packer
type Packer struct {
	maxBoxesToBalanceWeight  int
	items                    *ItemList
	boxes                    *BoxList
	boxQuantitiesAvailable   map[Box]int
	packedBoxSorter          PackedBoxSorter
	throwOnUnpackableItem    bool
	beStrictAboutItemOrdering bool
	timeoutChecker           TimeoutChecker
}

// NewPacker creates a new Packer
func NewPacker(items *ItemList, boxes *BoxList, packedBoxSorter PackedBoxSorter) *Packer {
	if items == nil {
		items = NewItemList()
	}
	if boxes == nil {
		boxes = NewBoxList()
	}
	if packedBoxSorter == nil {
		packedBoxSorter = &DefaultPackedBoxSorter{}
	}

	return &Packer{
		maxBoxesToBalanceWeight:   12,
		items:                     items,
		boxes:                     boxes,
		boxQuantitiesAvailable:    make(map[Box]int),
		packedBoxSorter:           packedBoxSorter,
		throwOnUnpackableItem:     true,
		beStrictAboutItemOrdering: false,
		timeoutChecker:            nil,
	}
}

// AddItem adds an item to be packed
func (p *Packer) AddItem(item Item, qty int) {
	p.items.Insert(item, qty)
}

// SetItems sets a list of items all at once
func (p *Packer) SetItems(items *ItemList) {
	p.items = ItemListFromArray(items.GetIterator(), false)
}

// AddBox adds a box size
func (p *Packer) AddBox(box Box) {
	p.boxes.Insert(box)
	if limitedBox, ok := box.(LimitedSupplyBox); ok {
		p.SetBoxQuantity(box, limitedBox.GetQuantityAvailable())
	} else {
		p.SetBoxQuantity(box, math.MaxInt32)
	}
}

// SetBoxes adds a pre-prepared set of boxes all at once
func (p *Packer) SetBoxes(boxList *BoxList) {
	p.boxes = boxList
	for _, box := range p.boxes.GetIterator() {
		if limitedBox, ok := box.(LimitedSupplyBox); ok {
			p.SetBoxQuantity(box, limitedBox.GetQuantityAvailable())
		} else {
			p.SetBoxQuantity(box, math.MaxInt32)
		}
	}
}

// SetBoxQuantity sets the quantity of this box type available
func (p *Packer) SetBoxQuantity(box Box, qty int) {
	p.boxQuantitiesAvailable[box] = qty
}

// GetMaxBoxesToBalanceWeight returns number of boxes at which balancing weight is deemed not worth the extra computation time
func (p *Packer) GetMaxBoxesToBalanceWeight() int {
	return p.maxBoxesToBalanceWeight
}

// SetMaxBoxesToBalanceWeight sets number of boxes at which balancing weight is deemed not worth the extra computation time
func (p *Packer) SetMaxBoxesToBalanceWeight(maxBoxesToBalanceWeight int) {
	p.maxBoxesToBalanceWeight = maxBoxesToBalanceWeight
}

// SetPackedBoxSorter sets the packed box sorter
func (p *Packer) SetPackedBoxSorter(packedBoxSorter PackedBoxSorter) {
	p.packedBoxSorter = packedBoxSorter
}

// SetTimeoutChecker sets the timeout checker
func (p *Packer) SetTimeoutChecker(timeoutChecker TimeoutChecker) {
	p.timeoutChecker = timeoutChecker
}

// ThrowOnUnpackableItem sets whether to throw on unpackable item
func (p *Packer) ThrowOnUnpackableItem(throwOnUnpackableItem bool) {
	p.throwOnUnpackableItem = throwOnUnpackableItem
}

// BeStrictAboutItemOrdering sets whether to be strict about item ordering
func (p *Packer) BeStrictAboutItemOrdering(beStrict bool) {
	p.beStrictAboutItemOrdering = beStrict
}

// GetUnpackedItems returns the items that haven't been packed
func (p *Packer) GetUnpackedItems() *ItemList {
	return p.items
}

// Pack packs items into boxes using built-in heuristics for the best solution
func (p *Packer) Pack() *PackedBoxList {
	if p.timeoutChecker != nil {
		p.timeoutChecker.Start(nil)
	}

	packedBoxes := p.DoBasicPacking(false)

	// If we have multiple boxes, try and optimise/even-out weight distribution
	if !p.beStrictAboutItemOrdering && packedBoxes.Count() > 1 && packedBoxes.Count() <= p.maxBoxesToBalanceWeight {
		redistributor := NewWeightRedistributor(p.boxes, p.packedBoxSorter, p.boxQuantitiesAvailable, p.timeoutChecker)
		packedBoxes = redistributor.RedistributeWeight(packedBoxes)
	}

	return packedBoxes
}

// DoBasicPacking does basic packing without weight redistribution
func (p *Packer) DoBasicPacking(enforceSingleBox bool) *PackedBoxList {
	packedBoxes := NewPackedBoxList(p.packedBoxSorter)

	// Keep going until everything packed
	for p.items.Count() > 0 {
		packedBoxesIteration := make([]*PackedBox, 0)

		// Loop through boxes starting with smallest, see what happens
		for _, box := range p.getBoxList(enforceSingleBox) {
			if p.timeoutChecker != nil {
				if err := p.timeoutChecker.ThrowOnTimeout(nil, ""); err != nil {
					panic(err)
				}
			}

			volumePacker := NewVolumePacker(box, p.items)
			volumePacker.BeStrictAboutItemOrdering(p.beStrictAboutItemOrdering)
			packedBox := volumePacker.Pack()

			if packedBox.Items.Count() > 0 {
				packedBoxesIteration = append(packedBoxesIteration, packedBox)

				// Have we found a single box that contains everything?
				if packedBox.Items.Count() == p.items.Count() {
					break
				}
			}
		}

		if len(packedBoxesIteration) > 0 {
			// Find best box of iteration, and remove packed items from unpacked list
			sort.Slice(packedBoxesIteration, func(i, j int) bool {
				return p.packedBoxSorter.Compare(packedBoxesIteration[i], packedBoxesIteration[j]) < 0
			})
			bestBox := packedBoxesIteration[0]

			p.items.RemovePackedItems(bestBox.Items)

			packedBoxes.Insert(bestBox)
			p.boxQuantitiesAvailable[bestBox.Box]--
		} else if p.throwOnUnpackableItem {
			topItem := p.items.Top()
			if topItem != nil {
				panic(exception.NewNoBoxesAvailableException("No boxes could be found for item '"+topItem.GetDescription()+"'", p.items))
			} else {
				panic(exception.NewNoBoxesAvailableException("No boxes could be found for items", p.items))
			}
		} else {
			break
		}
	}

	return packedBoxes
}

// PackAllPermutations packs items into boxes returning "all" possible box combination permutations
// Use with caution (will be slow) with a large number of box types!
func (p *Packer) PackAllPermutations() []*PackedBoxList {
	if p.timeoutChecker != nil {
		p.timeoutChecker.Start(nil)
	}

	boxQuantitiesAvailable := make(map[Box]int)
	for box, qty := range p.boxQuantitiesAvailable {
		boxQuantitiesAvailable[box] = qty
	}

	type WIPPermutation struct {
		permutation *PackedBoxList
		itemsLeft   *ItemList
	}

	wipPermutations := []WIPPermutation{
		{
			permutation: NewPackedBoxList(p.packedBoxSorter),
			itemsLeft:   p.items,
		},
	}
	completedPermutations := make([]*PackedBoxList, 0)

	// Keep going until everything packed
	for len(wipPermutations) > 0 {
		wipPermutation := wipPermutations[len(wipPermutations)-1]
		wipPermutations = wipPermutations[:len(wipPermutations)-1]

		remainingBoxQuantities := make(map[Box]int)
		for box, qty := range boxQuantitiesAvailable {
			remainingBoxQuantities[box] = qty
		}

		for _, packedBox := range wipPermutation.permutation.GetIterator() {
			remainingBoxQuantities[packedBox.Box]--
		}

		if wipPermutation.itemsLeft.Count() == 0 {
			completedPermutations = append(completedPermutations, wipPermutation.permutation)
			continue
		}

		additionalPermutationsForThisPermutation := make([]*PackedBox, 0)
		for _, box := range p.boxes.GetIterator() {
			if p.timeoutChecker != nil {
				if err := p.timeoutChecker.ThrowOnTimeout(nil, ""); err != nil {
					panic(err)
				}
			}

			if remainingBoxQuantities[box] > 0 {
				volumePacker := NewVolumePacker(box, wipPermutation.itemsLeft)
				packedBox := volumePacker.Pack()
				if packedBox.Items.Count() > 0 {
					additionalPermutationsForThisPermutation = append(additionalPermutationsForThisPermutation, packedBox)
				}
			}
		}

		if len(additionalPermutationsForThisPermutation) > 0 {
			for _, additionalPermutationForThisPermutation := range additionalPermutationsForThisPermutation {
				newPermutation := NewPackedBoxList(p.packedBoxSorter)
				for _, box := range wipPermutation.permutation.GetIterator() {
					newPermutation.Insert(box)
				}
				newPermutation.Insert(additionalPermutationForThisPermutation)

				itemsRemainingOnPermutation := ItemListFromArray(wipPermutation.itemsLeft.GetIterator(), false)
				itemsRemainingOnPermutation.RemovePackedItems(additionalPermutationForThisPermutation.Items)

				wipPermutations = append(wipPermutations, WIPPermutation{
					permutation: newPermutation,
					itemsLeft:   itemsRemainingOnPermutation,
				})
			}
		} else if p.throwOnUnpackableItem {
			topItem := wipPermutation.itemsLeft.Top()
			if topItem != nil {
				panic(exception.NewNoBoxesAvailableException("No boxes could be found for item '"+topItem.GetDescription()+"'", wipPermutation.itemsLeft))
			} else {
				panic(exception.NewNoBoxesAvailableException("No boxes could be found for items", wipPermutation.itemsLeft))
			}
		} else {
			if wipPermutation.permutation.Count() > 0 { // don't treat initial empty permutation as completed
				completedPermutations = append(completedPermutations, wipPermutation.permutation)
			}
		}
	}

	for _, completedPermutation := range completedPermutations {
		for _, packedBox := range completedPermutation.GetIterator() {
			p.items.RemovePackedItems(packedBox.Items)
		}
	}

	return completedPermutations
}

// getBoxList gets a "smart" ordering of the boxes to try packing items into
func (p *Packer) getBoxList(enforceSingleBox bool) []Box {
	itemVolume := 0
	for _, item := range p.items.GetIterator() {
		itemVolume += item.GetWidth() * item.GetLength() * item.GetDepth()
	}

	preferredBoxes := make([]Box, 0)
	otherBoxes := make([]Box, 0)

	for _, box := range p.boxes.GetIterator() {
		if p.boxQuantitiesAvailable[box] > 0 {
			boxVolume := box.GetInnerWidth() * box.GetInnerLength() * box.GetInnerDepth()
			if boxVolume >= itemVolume {
				preferredBoxes = append(preferredBoxes, box)
			} else if !enforceSingleBox {
				otherBoxes = append(otherBoxes, box)
			}
		}
	}

	return append(preferredBoxes, otherBoxes...)
}
