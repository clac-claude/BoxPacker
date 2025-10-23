package boxpacker

import "sort"

// Box packing (3D bin packing, knapsack problem).
//
// @author Doug Wright

// LayerPacker - Layer packer.
type LayerPacker struct {
	box                    Box
	singlePassMode         bool
	orientatedItemFactory  *OrientatedItemFactory
	beStrictAboutItemOrdering bool
	isBoxRotated           bool
}

// NewLayerPacker creates a new LayerPacker
func NewLayerPacker(box Box) *LayerPacker {
	orientatedItemFactory := NewOrientatedItemFactory(box)

	return &LayerPacker{
		box:                       box,
		singlePassMode:            false,
		orientatedItemFactory:     orientatedItemFactory,
		beStrictAboutItemOrdering: false,
		isBoxRotated:              false,
	}
}

// SetSinglePassMode sets single pass mode
func (lp *LayerPacker) SetSinglePassMode(singlePassMode bool) {
	lp.singlePassMode = singlePassMode
	lp.orientatedItemFactory.SetSinglePassMode(singlePassMode)
}

// SetBoxIsRotated sets whether the box is rotated
func (lp *LayerPacker) SetBoxIsRotated(boxIsRotated bool) {
	lp.isBoxRotated = boxIsRotated
	lp.orientatedItemFactory.SetBoxIsRotated(boxIsRotated)
}

// BeStrictAboutItemOrdering sets strict ordering mode
func (lp *LayerPacker) BeStrictAboutItemOrdering(beStrict bool) {
	lp.beStrictAboutItemOrdering = beStrict
}

// PackLayer packs items into an individual vertical layer
func (lp *LayerPacker) PackLayer(
	items *ItemList,
	packedItemList *PackedItemList,
	startX int,
	startY int,
	startZ int,
	widthForLayer int,
	lengthForLayer int,
	depthForLayer int,
	guidelineLayerDepth int,
	considerStability bool,
	firstItem *OrientatedItem,
) *PackedLayer {
	layer := NewPackedLayer()
	x := startX
	y := startY
	z := startZ
	rowLength := 0
	var prevItem *OrientatedItem = nil
	skippedItems := make([]Item, 0)

	for items.Count() > 0 {
		itemToPack := items.Extract()
		if itemToPack == nil {
			break
		}

		// skip items that will never fit e.g. too heavy
		if itemToPack.GetWeight() > (lp.box.GetMaxWeight() - lp.box.GetEmptyWeight() - packedItemList.GetWeight()) {
			continue
		}

		var orientatedItem *OrientatedItem
		if firstItem != nil && firstItem.Item == itemToPack {
			orientatedItem = firstItem
			firstItem = nil
		} else {
			orientatedItem = lp.orientatedItemFactory.GetBestOrientation(
				itemToPack,
				prevItem,
				items,
				widthForLayer-x,
				lengthForLayer-y,
				depthForLayer,
				rowLength,
				x, y, z,
				packedItemList,
				considerStability,
			)
		}

		if orientatedItem != nil {
			packedItem := PackedItemFromOrientatedItem(orientatedItem, x, y, z)
			layer.Insert(packedItem)
			packedItemList.Insert(packedItem)

			if packedItem.Length > rowLength {
				rowLength = packedItem.Length
			}
			prevItem = orientatedItem

			// Figure out if we can stack items on top of this rather than side by side
			// e.g. when we've packed a tall item, and have just put a shorter one next to it.
			layerDepth := layer.GetDepth()
			if guidelineLayerDepth > 0 {
				layerDepth = guidelineLayerDepth
			}
			stackableDepth := layerDepth - packedItem.Depth
			if stackableDepth > 0 {
				stackedLayer := lp.PackLayer(
					items,
					packedItemList,
					x,
					y,
					z+packedItem.Depth,
					x+packedItem.Width,
					y+packedItem.Length,
					stackableDepth,
					stackableDepth,
					considerStability,
					nil,
				)
				layer.Merge(stackedLayer)
			}

			x += packedItem.Width

			// might be space available lengthwise across the width of this item, up to the current layer length
			layer.Merge(lp.PackLayer(
				items,
				packedItemList,
				x-packedItem.Width,
				y+packedItem.Length,
				z,
				x,
				y+rowLength,
				depthForLayer,
				layer.GetDepth(),
				considerStability,
				nil,
			))

			if items.Count() == 0 && len(skippedItems) > 0 {
				items = ItemListFromArray(append(skippedItems, items.GetIterator()...), true)
				skippedItems = make([]Item, 0)
			}

			continue
		}

		if !lp.beStrictAboutItemOrdering && items.Count() > 0 { // skip for now, move on to the next item
			skippedItems = append(skippedItems, itemToPack)
			// abandon here if next item is the same, no point trying to keep going. Last time is not skipped, need that to trigger appropriate reset logic
			for items.Count() > 1 && isSameDimensions(itemToPack, items.Top()) {
				skippedItems = append(skippedItems, items.Extract())
			}
			continue
		}

		if x > startX {
			y += rowLength
			x = startX
			rowLength = 0
			skippedItems = append(skippedItems, itemToPack)
			items = ItemListFromArray(append(skippedItems, items.GetIterator()...), true)
			skippedItems = make([]Item, 0)
			prevItem = nil
			continue
		}

		skippedItems = append(skippedItems, itemToPack)
		items = ItemListFromArray(append(skippedItems, items.GetIterator()...), true)

		return layer
	}

	return layer
}

// isSameDimensions compares two items to see if they have same dimensions
func isSameDimensions(itemA Item, itemB Item) bool {
	if itemA == itemB {
		return true
	}

	itemADimensions := []int{itemA.GetWidth(), itemA.GetLength(), itemA.GetDepth()}
	itemBDimensions := []int{itemB.GetWidth(), itemB.GetLength(), itemB.GetDepth()}
	sort.Ints(itemADimensions)
	sort.Ints(itemBDimensions)

	if len(itemADimensions) != len(itemBDimensions) {
		return false
	}

	for i := range itemADimensions {
		if itemADimensions[i] != itemBDimensions[i] {
			return false
		}
	}

	return true
}
