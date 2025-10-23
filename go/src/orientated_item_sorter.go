package boxpacker

import (
	"fmt"
	"math"
)

// Box packing (3D bin packing, knapsack problem).
//
// @author Doug Wright

// OrientatedItemSorter - Figure out best choice of orientations for an item and a given context.
type OrientatedItemSorter struct {
	orientatedItemFactory *OrientatedItemFactory
	singlePassMode        bool
	widthLeft             int
	lengthLeft            int
	depthLeft             int
	nextItems             *ItemList
	rowLength             int
	x                     int
	y                     int
	z                     int
	prevPackedItemList    *PackedItemList
}

var lookaheadCache = make(map[string]int)

// NewOrientatedItemSorter creates a new OrientatedItemSorter
func NewOrientatedItemSorter(
	orientatedItemFactory *OrientatedItemFactory,
	singlePassMode bool,
	widthLeft int,
	lengthLeft int,
	depthLeft int,
	nextItems *ItemList,
	rowLength int,
	x int,
	y int,
	z int,
	prevPackedItemList *PackedItemList,
) *OrientatedItemSorter {
	return &OrientatedItemSorter{
		orientatedItemFactory: orientatedItemFactory,
		singlePassMode:        singlePassMode,
		widthLeft:             widthLeft,
		lengthLeft:            lengthLeft,
		depthLeft:             depthLeft,
		nextItems:             nextItems,
		rowLength:             rowLength,
		x:                     x,
		y:                     y,
		z:                     z,
		prevPackedItemList:    prevPackedItemList,
	}
}

// Compare implements comparison logic for sorting orientated items
func (ois *OrientatedItemSorter) Compare(a, b *OrientatedItem) int {
	// Prefer exact fits in width/length/depth order
	orientationAWidthLeft := ois.widthLeft - a.Width
	orientationBWidthLeft := ois.widthLeft - b.Width
	widthDecider := ois.exactFitDecider(orientationAWidthLeft, orientationBWidthLeft)
	if widthDecider != 0 {
		return widthDecider
	}

	orientationALengthLeft := ois.lengthLeft - a.Length
	orientationBLengthLeft := ois.lengthLeft - b.Length
	lengthDecider := ois.exactFitDecider(orientationALengthLeft, orientationBLengthLeft)
	if lengthDecider != 0 {
		return lengthDecider
	}

	orientationADepthLeft := ois.depthLeft - a.Depth
	orientationBDepthLeft := ois.depthLeft - b.Depth
	depthDecider := ois.exactFitDecider(orientationADepthLeft, orientationBDepthLeft)
	if depthDecider != 0 {
		return depthDecider
	}

	// prefer leaving room for next item(s)
	followingItemDecider := ois.lookAheadDecider(a, b, orientationAWidthLeft, orientationBWidthLeft)
	if followingItemDecider != 0 {
		return followingItemDecider
	}

	// otherwise prefer leaving minimum possible gap, or the greatest footprint
	orientationAMinGap := orientationAWidthLeft
	if orientationALengthLeft < orientationAMinGap {
		orientationAMinGap = orientationALengthLeft
	}

	orientationBMinGap := orientationBWidthLeft
	if orientationBLengthLeft < orientationBMinGap {
		orientationBMinGap = orientationBLengthLeft
	}

	if orientationAMinGap != orientationBMinGap {
		if orientationAMinGap < orientationBMinGap {
			return -1
		}
		return 1
	}

	if a.SurfaceFootprint < b.SurfaceFootprint {
		return -1
	}
	if a.SurfaceFootprint > b.SurfaceFootprint {
		return 1
	}

	return 0
}

// lookAheadDecider decides based on looking ahead at next items
func (ois *OrientatedItemSorter) lookAheadDecider(a, b *OrientatedItem, orientationAWidthLeft, orientationBWidthLeft int) int {
	if ois.nextItems.Count() == 0 {
		return 0
	}

	nextItemFitA := ois.orientatedItemFactory.GetPossibleOrientations(ois.nextItems.Top(), a, orientationAWidthLeft, ois.lengthLeft, ois.depthLeft, ois.x, ois.y, ois.z, ois.prevPackedItemList)
	nextItemFitB := ois.orientatedItemFactory.GetPossibleOrientations(ois.nextItems.Top(), b, orientationBWidthLeft, ois.lengthLeft, ois.depthLeft, ois.x, ois.y, ois.z, ois.prevPackedItemList)

	if len(nextItemFitA) > 0 && len(nextItemFitB) == 0 {
		return -1
	}
	if len(nextItemFitB) > 0 && len(nextItemFitA) == 0 {
		return 1
	}

	// if not an easy either/or, do a partial lookahead
	additionalPackedA := ois.calculateAdditionalItemsPackedWithThisOrientation(a)
	additionalPackedB := ois.calculateAdditionalItemsPackedWithThisOrientation(b)

	if additionalPackedB > additionalPackedA {
		return -1
	}
	if additionalPackedB < additionalPackedA {
		return 1
	}

	return 0
}

// calculateAdditionalItemsPackedWithThisOrientation - Approximation of a forward-looking packing.
// Not an actual packing, that has additional logic regarding constraints and stackability, this focuses
// purely on fit.
func (ois *OrientatedItemSorter) calculateAdditionalItemsPackedWithThisOrientation(prevItem *OrientatedItem) int {
	if ois.singlePassMode {
		return 0
	}

	currentRowLength := prevItem.Length
	if ois.rowLength > currentRowLength {
		currentRowLength = ois.rowLength
	}

	itemsToPack := ois.nextItems.TopN(8) // cap lookahead as this gets recursive and slow

	cacheKey := fmt.Sprintf("%d|%d|%d|%d|%d|%d",
		ois.widthLeft,
		ois.lengthLeft,
		prevItem.Width,
		prevItem.Length,
		currentRowLength,
		ois.depthLeft)

	for _, itemToPack := range itemsToPack.GetIterator() {
		cacheKey += fmt.Sprintf("|%d|%d|%d|%d|%d",
			itemToPack.GetWidth(),
			itemToPack.GetLength(),
			itemToPack.GetDepth(),
			itemToPack.GetWeight(),
			itemToPack.GetAllowedRotation())
	}

	if cachedValue, exists := lookaheadCache[cacheKey]; exists {
		return cachedValue
	}

	tempBox := NewWorkingVolume(ois.widthLeft-prevItem.Width, currentRowLength, ois.depthLeft, math.MaxInt32)
	tempPacker := NewVolumePacker(tempBox, itemsToPack)
	tempPacker.SetSinglePassMode(true)
	remainingRowPacked := tempPacker.Pack()

	itemsToPack.RemovePackedItems(remainingRowPacked.Items)

	tempBox = NewWorkingVolume(ois.widthLeft, ois.lengthLeft-currentRowLength, ois.depthLeft, math.MaxInt32)
	tempPacker = NewVolumePacker(tempBox, itemsToPack)
	tempPacker.SetSinglePassMode(true)
	nextRowsPacked := tempPacker.Pack()

	itemsToPack.RemovePackedItems(nextRowsPacked.Items)

	packedCount := ois.nextItems.Count() - itemsToPack.Count()

	lookaheadCache[cacheKey] = packedCount

	return packedCount
}

// exactFitDecider helper for exact fit comparison
func (ois *OrientatedItemSorter) exactFitDecider(dimensionALeft, dimensionBLeft int) int {
	if dimensionALeft == 0 && dimensionBLeft > 0 {
		return -1
	}

	if dimensionALeft > 0 && dimensionBLeft == 0 {
		return 1
	}

	return 0
}
