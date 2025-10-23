package main

import (
	"sort"
)

// SimplePacker is a simplified packer for lookahead calculations
type SimplePacker struct {
	Box   *Box
	Items []*Item
}

// PackSimple performs a simple single-pass packing
func (p *SimplePacker) PackSimple() PackedItemList {
	packed := make(PackedItemList, 0, len(p.Items))
	x, y, z := 0, 0, 0
	rowLength := 0
	packedWeight := 0

	for _, item := range p.Items {
		// Get possible orientations
		orientations := GetPossibleOrientations(
			item,
			nil,
			p.Box.InnerWidth-x,
			p.Box.InnerLength-y,
			p.Box.InnerDepth-z,
			packedWeight,
			p.Box.MaxWeight,
		)

		if len(orientations) == 0 {
			// Try next row
			if x > 0 {
				y += rowLength
				x = 0
				rowLength = 0

				// Retry with new position
				orientations = GetPossibleOrientations(
					item,
					nil,
					p.Box.InnerWidth,
					p.Box.InnerLength-y,
					p.Box.InnerDepth,
					packedWeight,
					p.Box.MaxWeight,
				)
			}

			if len(orientations) == 0 {
				continue // Skip this item
			}
		}

		// Use first stable orientation
		usable := GetUsableOrientations(orientations, p.Box.InnerDepth)
		if len(usable) == 0 {
			continue
		}

		orientation := usable[0]
		packedItem := NewPackedItem(orientation, x, y, z)
		packed = append(packed, packedItem)

		packedWeight += item.Weight
		rowLength = max(rowLength, orientation.Length)
		x += orientation.Width
	}

	return packed
}

// LookaheadCache stores lookahead results
var lookaheadCache = make(map[string]int)

// CalculateLookahead performs lookahead calculation for orientation selection
func CalculateLookahead(
	prevItem *OrientatedItem,
	nextItems []*Item,
	widthLeft, lengthLeft, depthLeft, rowLength int,
	maxLookahead int,
) int {
	if len(nextItems) == 0 || maxLookahead <= 0 {
		return 0
	}

	// Limit lookahead items
	itemsToCheck := nextItems
	if len(nextItems) > maxLookahead {
		itemsToCheck = nextItems[:maxLookahead]
	}

	// Build cache key
	cacheKey := buildCacheKey(prevItem, itemsToCheck, widthLeft, lengthLeft, depthLeft, rowLength)
	if cached, ok := lookaheadCache[cacheKey]; ok {
		return cached
	}

	currentRowLength := max(prevItem.Length, rowLength)

	// Pack remaining row space
	rowBox := &Box{
		InnerWidth:  widthLeft - prevItem.Width,
		InnerLength: currentRowLength,
		InnerDepth:  depthLeft,
		MaxWeight:   1000000, // Large number for lookahead
	}

	rowPacker := &SimplePacker{
		Box:   rowBox,
		Items: copyItems(itemsToCheck),
	}
	rowPacked := rowPacker.PackSimple()

	// Remove packed items from next items
	remainingItems := removePackedItems(itemsToCheck, rowPacked)

	// Pack next rows
	nextRowBox := &Box{
		InnerWidth:  widthLeft,
		InnerLength: lengthLeft - currentRowLength,
		InnerDepth:  depthLeft,
		MaxWeight:   1000000,
	}

	nextRowPacker := &SimplePacker{
		Box:   nextRowBox,
		Items: remainingItems,
	}
	nextRowPacked := nextRowPacker.PackSimple()

	// Calculate how many items were packed
	packedCount := len(rowPacked) + len(nextRowPacked)

	// Cache the result
	lookaheadCache[cacheKey] = packedCount

	return packedCount
}

// SortOrientationsByLookahead sorts orientations based on lookahead results
func SortOrientationsByLookahead(
	orientations []*OrientatedItem,
	nextItems []*Item,
	widthLeft, lengthLeft, depthLeft, rowLength int,
	maxLookahead int,
) []*OrientatedItem {
	type scoredOrientation struct {
		orientation *OrientatedItem
		score       int
		minGap      int
	}

	scored := make([]scoredOrientation, len(orientations))

	for i, orientation := range orientations {
		lookaheadScore := CalculateLookahead(
			orientation,
			nextItems,
			widthLeft,
			lengthLeft,
			depthLeft,
			rowLength,
			maxLookahead,
		)

		orientationWidthLeft := widthLeft - orientation.Width
		orientationLengthLeft := lengthLeft - orientation.Length
		minGap := min(orientationWidthLeft, orientationLengthLeft)

		scored[i] = scoredOrientation{
			orientation: orientation,
			score:       lookaheadScore,
			minGap:      minGap,
		}
	}

	// Sort: higher lookahead score first, then smaller gap, then larger footprint
	sort.Slice(scored, func(i, j int) bool {
		if scored[i].score != scored[j].score {
			return scored[i].score > scored[j].score
		}
		if scored[i].minGap != scored[j].minGap {
			return scored[i].minGap < scored[j].minGap
		}
		return scored[i].orientation.SurfaceFootprint > scored[j].orientation.SurfaceFootprint
	})

	result := make([]*OrientatedItem, len(orientations))
	for i, s := range scored {
		result[i] = s.orientation
	}

	return result
}

// Helper functions

func buildCacheKey(prevItem *OrientatedItem, items []*Item, widthLeft, lengthLeft, depthLeft, rowLength int) string {
	// Simple string concatenation for cache key
	key := ""
	for i := 0; i < 5 && i < len(items); i++ {
		item := items[i]
		key += string(rune(item.Width)) + string(rune(item.Length)) + string(rune(item.Depth))
	}
	key += string(rune(widthLeft)) + string(rune(lengthLeft)) + string(rune(depthLeft)) + string(rune(rowLength))
	if prevItem != nil {
		key += string(rune(prevItem.Width)) + string(rune(prevItem.Length))
	}
	return key
}

func copyItems(items []*Item) []*Item {
	result := make([]*Item, len(items))
	copy(result, items)
	return result
}

func removePackedItems(items []*Item, packed PackedItemList) []*Item {
	remaining := make([]*Item, 0, len(items))
	packedMap := make(map[*Item]int)

	for _, p := range packed {
		packedMap[p.Item.Item]++
	}

	for _, item := range items {
		if count, exists := packedMap[item]; exists && count > 0 {
			packedMap[item]--
		} else {
			remaining = append(remaining, item)
		}
	}

	return remaining
}

// ClearLookaheadCache clears the lookahead cache
func ClearLookaheadCache() {
	lookaheadCache = make(map[string]int)
}
