package boxpacker

import "sort"

// Box packing (3D bin packing, knapsack problem).
//
// @author Doug Wright

// LayerStabiliser - Applies load stability to generated result.
type LayerStabiliser struct{}

// NewLayerStabiliser creates a new LayerStabiliser
func NewLayerStabiliser() *LayerStabiliser {
	return &LayerStabiliser{}
}

// Stabilise re-orders layers and recalculates z positions
func (ls *LayerStabiliser) Stabilise(packedLayers []*PackedLayer) []*PackedLayer {
	// first re-order according to footprint
	stabilisedLayers := make([]*PackedLayer, 0)
	sort.Slice(packedLayers, func(i, j int) bool {
		return ls.compare(packedLayers[i], packedLayers[j]) < 0
	})

	// then for each item in the layer, re-calculate each item's z position
	currentZ := 0
	for _, oldZLayer := range packedLayers {
		oldZStart := oldZLayer.GetStartZ()
		newZLayer := NewPackedLayer()

		for _, oldZItem := range oldZLayer.GetItems() {
			newZ := oldZItem.Z - oldZStart + currentZ
			newZItem := NewPackedItem(oldZItem.Item, oldZItem.X, oldZItem.Y, newZ, oldZItem.Width, oldZItem.Length, oldZItem.Depth)
			newZLayer.Insert(newZItem)
		}

		stabilisedLayers = append(stabilisedLayers, newZLayer)
		currentZ += newZLayer.GetDepth()
	}

	return stabilisedLayers
}

// compare compares two layers
func (ls *LayerStabiliser) compare(layerA, layerB *PackedLayer) int {
	footprintA := layerA.GetFootprint()
	footprintB := layerB.GetFootprint()

	if footprintB != footprintA {
		if footprintB > footprintA {
			return -1
		}
		return 1
	}

	depthA := layerA.GetDepth()
	depthB := layerB.GetDepth()

	if depthB > depthA {
		return -1
	}
	if depthB < depthA {
		return 1
	}

	return 0
}
