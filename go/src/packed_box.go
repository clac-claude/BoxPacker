package boxpacker

import (
	"fmt"
	"math"
)

// Box packing (3D bin packing, knapsack problem).
//
// @author Doug Wright

// PackedBox - A "box" with items.
type PackedBox struct {
	Box                Box
	Items              *PackedItemList
	itemWeight         int
	volumeUtilisation  float64
	hasItemWeight      bool
	hasVolumeUtil      bool
}

// NewPackedBox creates a new PackedBox
func NewPackedBox(box Box, items *PackedItemList) *PackedBox {
	pb := &PackedBox{
		Box:                box,
		Items:              items,
		hasItemWeight:      false,
		hasVolumeUtil:      false,
	}
	pb.assertPackingCompliesWithRealWorld()
	return pb
}

// GetWeight returns the packed weight (in grams)
func (pb *PackedBox) GetWeight() int {
	return pb.Box.GetEmptyWeight() + pb.GetItemWeight()
}

// GetItemWeight returns the packed weight of the items only (in grams)
func (pb *PackedBox) GetItemWeight() int {
	if !pb.hasItemWeight {
		itemWeight := 0
		for _, item := range pb.Items.list {
			itemWeight += item.Item.GetWeight()
		}
		pb.itemWeight = itemWeight
		pb.hasItemWeight = true
	}

	return pb.itemWeight
}

// GetRemainingWidth returns the remaining width inside box for another item
func (pb *PackedBox) GetRemainingWidth() int {
	return pb.Box.GetInnerWidth() - pb.GetUsedWidth()
}

// GetRemainingLength returns the remaining length inside box for another item
func (pb *PackedBox) GetRemainingLength() int {
	return pb.Box.GetInnerLength() - pb.GetUsedLength()
}

// GetRemainingDepth returns the remaining depth inside box for another item
func (pb *PackedBox) GetRemainingDepth() int {
	return pb.Box.GetInnerDepth() - pb.GetUsedDepth()
}

// GetUsedWidth returns the used width inside box for packing items
func (pb *PackedBox) GetUsedWidth() int {
	maxWidth := 0

	for _, item := range pb.Items.list {
		end := item.X + item.Width
		if end > maxWidth {
			maxWidth = end
		}
	}

	return maxWidth
}

// GetUsedLength returns the used length inside box for packing items
func (pb *PackedBox) GetUsedLength() int {
	maxLength := 0

	for _, item := range pb.Items.list {
		end := item.Y + item.Length
		if end > maxLength {
			maxLength = end
		}
	}

	return maxLength
}

// GetUsedDepth returns the used depth inside box for packing items
func (pb *PackedBox) GetUsedDepth() int {
	maxDepth := 0

	for _, item := range pb.Items.list {
		end := item.Z + item.Depth
		if end > maxDepth {
			maxDepth = end
		}
	}

	return maxDepth
}

// GetRemainingWeight returns the remaining weight inside box for another item
func (pb *PackedBox) GetRemainingWeight() int {
	return pb.Box.GetMaxWeight() - pb.GetWeight()
}

// GetInnerVolume returns the inner volume of the box
func (pb *PackedBox) GetInnerVolume() int {
	return pb.Box.GetInnerWidth() * pb.Box.GetInnerLength() * pb.Box.GetInnerDepth()
}

// GetUsedVolume returns the used volume of the packed box
func (pb *PackedBox) GetUsedVolume() int {
	return pb.Items.GetVolume()
}

// GetUnusedVolume returns the unused volume of the packed box
func (pb *PackedBox) GetUnusedVolume() int {
	return pb.GetInnerVolume() - pb.GetUsedVolume()
}

// GetVolumeUtilisation returns the volume utilisation of the packed box
func (pb *PackedBox) GetVolumeUtilisation() float64 {
	if !pb.hasVolumeUtil {
		innerVol := pb.GetInnerVolume()
		if innerVol == 0 {
			innerVol = 1
		}
		pb.volumeUtilisation = math.Round(float64(pb.GetUsedVolume())/float64(innerVol)*1000) / 10
		pb.hasVolumeUtil = true
	}

	return pb.volumeUtilisation
}

// assertPackingCompliesWithRealWorld validates that all items are placed solely within the confines of the box,
// and that no two items are placed into the same physical space.
func (pb *PackedBox) assertPackingCompliesWithRealWorld() {
	itemsToCheck := make([]*PackedItem, len(pb.Items.list))
	copy(itemsToCheck, pb.Items.list)

	for len(itemsToCheck) > 0 {
		itemToCheck := itemsToCheck[len(itemsToCheck)-1]
		itemsToCheck = itemsToCheck[:len(itemsToCheck)-1]

		if itemToCheck.X < 0 {
			panic(fmt.Sprintf("Item X coordinate %d is negative", itemToCheck.X))
		}
		if itemToCheck.X+itemToCheck.Width > pb.Box.GetInnerWidth() {
			panic(fmt.Sprintf("Item exceeds box width: %d + %d > %d", itemToCheck.X, itemToCheck.Width, pb.Box.GetInnerWidth()))
		}
		if itemToCheck.Y < 0 {
			panic(fmt.Sprintf("Item Y coordinate %d is negative", itemToCheck.Y))
		}
		if itemToCheck.Y+itemToCheck.Length > pb.Box.GetInnerLength() {
			panic(fmt.Sprintf("Item exceeds box length: %d + %d > %d", itemToCheck.Y, itemToCheck.Length, pb.Box.GetInnerLength()))
		}
		if itemToCheck.Z < 0 {
			panic(fmt.Sprintf("Item Z coordinate %d is negative", itemToCheck.Z))
		}
		if itemToCheck.Z+itemToCheck.Depth > pb.Box.GetInnerDepth() {
			panic(fmt.Sprintf("Item exceeds box depth: %d + %d > %d", itemToCheck.Z, itemToCheck.Depth, pb.Box.GetInnerDepth()))
		}

		for _, otherItem := range itemsToCheck {
			hasXOverlap := itemToCheck.X < (otherItem.X+otherItem.Width) && otherItem.X < (itemToCheck.X+itemToCheck.Width)
			hasYOverlap := itemToCheck.Y < (otherItem.Y+otherItem.Length) && otherItem.Y < (itemToCheck.Y+itemToCheck.Length)
			hasZOverlap := itemToCheck.Z < (otherItem.Z+otherItem.Depth) && otherItem.Z < (itemToCheck.Z+itemToCheck.Depth)

			hasOverlap := hasXOverlap && hasYOverlap && hasZOverlap
			if hasOverlap {
				panic("Items overlap in the packed box")
			}
		}
	}
}
