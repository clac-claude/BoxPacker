package boxpacker

import "sort"

// Box packing (3D bin packing, knapsack problem).
//
// @author Doug Wright

// VolumePacker - Actual packer for a single box
type VolumePacker struct {
	box                       Box
	items                     *ItemList
	singlePassMode            bool
	packAcrossWidthOnly       bool
	layerPacker               *LayerPacker
	beStrictAboutItemOrdering bool
	hasConstrainedItems       bool
	hasNoRotationItems        bool
}

// NewVolumePacker creates a new VolumePacker
func NewVolumePacker(box Box, items *ItemList) *VolumePacker {
	// Clone items
	clonedItems := ItemListFromArray(items.GetIterator(), false)

	hasConstrainedItems := items.HasConstrainedItems()
	hasNoRotationItems := items.HasNoRotationItems()

	layerPacker := NewLayerPacker(box)

	return &VolumePacker{
		box:                       box,
		items:                     clonedItems,
		singlePassMode:            false,
		packAcrossWidthOnly:       false,
		layerPacker:               layerPacker,
		beStrictAboutItemOrdering: false,
		hasConstrainedItems:       hasConstrainedItems,
		hasNoRotationItems:        hasNoRotationItems,
	}
}

// PackAcrossWidthOnly sets pack across width only mode
func (vp *VolumePacker) PackAcrossWidthOnly() {
	vp.packAcrossWidthOnly = true
}

// BeStrictAboutItemOrdering sets strict ordering mode
func (vp *VolumePacker) BeStrictAboutItemOrdering(beStrict bool) {
	vp.beStrictAboutItemOrdering = beStrict
	vp.layerPacker.BeStrictAboutItemOrdering(beStrict)
}

// SetSinglePassMode sets single pass mode
func (vp *VolumePacker) SetSinglePassMode(singlePassMode bool) {
	vp.singlePassMode = singlePassMode
	if singlePassMode {
		vp.packAcrossWidthOnly = true
	}
	vp.layerPacker.SetSinglePassMode(singlePassMode)
}

// Pack packs as many items as possible into specific given box
func (vp *VolumePacker) Pack() *PackedBox {
	orientatedItemFactory := NewOrientatedItemFactory(vp.box)

	// Sometimes "space available" decisions depend on orientation of the box, so try both ways
	rotationsToTest := []bool{false}
	if !vp.packAcrossWidthOnly && !vp.hasNoRotationItems {
		rotationsToTest = append(rotationsToTest, true)
	}

	// The orientation of the first item can have an outsized effect on the rest of the placement, so special-case
	// that and try everything

	boxPermutations := make([]*PackedBox, 0)
	for _, rotation := range rotationsToTest {
		var boxWidth, boxLength int
		if rotation {
			boxWidth = vp.box.GetInnerLength()
			boxLength = vp.box.GetInnerWidth()
		} else {
			boxWidth = vp.box.GetInnerWidth()
			boxLength = vp.box.GetInnerLength()
		}

		specialFirstItemOrientations := []*OrientatedItem{nil}
		if !vp.singlePassMode {
			topItem := vp.items.Top()
			if topItem != nil {
				possibleOrientations := orientatedItemFactory.GetPossibleOrientations(
					topItem,
					nil,
					boxWidth,
					boxLength,
					vp.box.GetInnerDepth(),
					0, 0, 0,
					NewPackedItemList(),
				)
				if len(possibleOrientations) > 0 {
					specialFirstItemOrientations = possibleOrientations
				}
			}
		}

		for _, firstItemOrientation := range specialFirstItemOrientations {
			boxPermutation := vp.packRotation(boxWidth, boxLength, firstItemOrientation)
			if boxPermutation.Items.Count() == vp.items.Count() {
				return boxPermutation
			}

			boxPermutations = append(boxPermutations, boxPermutation)
		}
	}

	sort.Slice(boxPermutations, func(i, j int) bool {
		return boxPermutations[j].GetVolumeUtilisation() < boxPermutations[i].GetVolumeUtilisation()
	})

	if len(boxPermutations) > 0 {
		return boxPermutations[0]
	}

	return NewPackedBox(vp.box, NewPackedItemList())
}

// packRotation packs as many items as possible into specific given box with specific rotation
func (vp *VolumePacker) packRotation(boxWidth, boxLength int, firstItemOrientation *OrientatedItem) *PackedBox {
	vp.layerPacker.SetBoxIsRotated(vp.box.GetInnerWidth() != boxWidth)

	layers := make([]*PackedLayer, 0)
	items := ItemListFromArray(vp.items.GetIterator(), false)

	for items.Count() > 0 {
		layerStartDepth := getCurrentPackedDepth(layers)
		packedItemList := getPackedItemList(layers)

		if packedItemList.Count() > 0 {
			firstItemOrientation = nil
		}

		// do a preliminary layer pack to get the depth used
		preliminaryItems := ItemListFromArray(items.GetIterator(), false)
		preliminaryPackedItemList := NewPackedItemList()
		for _, item := range packedItemList.list {
			preliminaryPackedItemList.Insert(item)
		}

		preliminaryLayer := vp.layerPacker.PackLayer(
			preliminaryItems,
			preliminaryPackedItemList,
			0,
			0,
			layerStartDepth,
			boxWidth,
			boxLength,
			vp.box.GetInnerDepth()-layerStartDepth,
			0,
			true,
			firstItemOrientation,
		)

		if len(preliminaryLayer.GetItems()) == 0 {
			break
		}

		preliminaryLayerDepth := preliminaryLayer.GetDepth()
		if preliminaryLayerDepth == preliminaryLayer.GetItems()[0].Depth { // preliminary === final
			layers = append(layers, preliminaryLayer)
			items = preliminaryItems
		} else { // redo with now-known-depth so that we can stack to that height from the first item
			layer := vp.layerPacker.PackLayer(
				items,
				packedItemList,
				0,
				0,
				layerStartDepth,
				boxWidth,
				boxLength,
				vp.box.GetInnerDepth()-layerStartDepth,
				preliminaryLayerDepth,
				true,
				firstItemOrientation,
			)
			layers = append(layers, layer)
		}
	}

	if !vp.singlePassMode && len(layers) > 0 {
		layers = vp.stabiliseLayers(layers)

		// having packed layers, there may be tall, narrow gaps at the ends that can be utilised
		maxLayerWidth := 0
		for _, layer := range layers {
			endX := layer.GetEndX()
			if endX > maxLayerWidth {
				maxLayerWidth = endX
			}
		}

		layer1 := vp.layerPacker.PackLayer(
			items,
			getPackedItemList(layers),
			maxLayerWidth,
			0,
			0,
			boxWidth,
			boxLength,
			vp.box.GetInnerDepth(),
			vp.box.GetInnerDepth(),
			false,
			nil,
		)
		layers = append(layers, layer1)

		maxLayerLength := 0
		for _, layer := range layers {
			endY := layer.GetEndY()
			if endY > maxLayerLength {
				maxLayerLength = endY
			}
		}

		layer2 := vp.layerPacker.PackLayer(
			items,
			getPackedItemList(layers),
			0,
			maxLayerLength,
			0,
			boxWidth,
			boxLength,
			vp.box.GetInnerDepth(),
			vp.box.GetInnerDepth(),
			false,
			nil,
		)
		layers = append(layers, layer2)
	}

	layers = vp.correctLayerRotation(layers, boxWidth)

	return NewPackedBox(vp.box, getPackedItemList(layers))
}

// stabiliseLayers reorders layers so that the ones with the greatest surface area are placed at the bottom
func (vp *VolumePacker) stabiliseLayers(oldLayers []*PackedLayer) []*PackedLayer {
	if vp.hasConstrainedItems || vp.beStrictAboutItemOrdering { // constraints include position, so cannot change
		return oldLayers
	}

	stabiliser := NewLayerStabiliser()
	return stabiliser.Stabilise(oldLayers)
}

// correctLayerRotation swaps back width/length of the packed items to match orientation of the box if needed
func (vp *VolumePacker) correctLayerRotation(oldLayers []*PackedLayer, boxWidth int) []*PackedLayer {
	if vp.box.GetInnerWidth() == boxWidth {
		return oldLayers
	}

	newLayers := make([]*PackedLayer, 0)
	for _, originalLayer := range oldLayers {
		newLayer := NewPackedLayer()
		for _, item := range originalLayer.GetItems() {
			packedItem := NewPackedItem(
				item.Item,
				item.Y,
				item.X,
				item.Z,
				item.Length,
				item.Width,
				item.Depth,
			)
			newLayer.Insert(packedItem)
		}
		newLayers = append(newLayers, newLayer)
	}

	return newLayers
}

// getPackedItemList generates a single list of items packed from layers
func getPackedItemList(layers []*PackedLayer) *PackedItemList {
	packedItemList := NewPackedItemList()
	for _, layer := range layers {
		for _, packedItem := range layer.GetItems() {
			packedItemList.Insert(packedItem)
		}
	}

	return packedItemList
}

// getCurrentPackedDepth returns the current packed depth
func getCurrentPackedDepth(layers []*PackedLayer) int {
	depth := 0
	for _, layer := range layers {
		depth += layer.GetDepth()
	}

	return depth
}
