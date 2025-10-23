package main

/*
#include <stdlib.h>

typedef struct {
	int width;
	int length;
	int depth;
	int weight;
	int rotation;
} CItem;

typedef struct {
	int innerWidth;
	int innerLength;
	int innerDepth;
	int maxWeight;
} CBox;

typedef struct {
	int width;
	int length;
	int depth;
	int surfaceFootprint;
} COrientatedItem;
*/
import "C"
import (
	"unsafe"
)

//export CalculateLookaheadFFI
func CalculateLookaheadFFI(
	prevItemWidth, prevItemLength, prevItemDepth C.int,
	items *C.CItem,
	itemCount C.int,
	widthLeft, lengthLeft, depthLeft, rowLength C.int,
	maxLookahead C.int,
) C.int {
	// Convert C structs to Go structs
	prevItem := &OrientatedItem{
		Width:  int(prevItemWidth),
		Length: int(prevItemLength),
		Depth:  int(prevItemDepth),
	}

	// Convert items array
	itemsSlice := unsafe.Slice(items, itemCount)
	goItems := make([]*Item, itemCount)
	for i, cItem := range itemsSlice {
		goItems[i] = &Item{
			Width:    int(cItem.width),
			Length:   int(cItem.length),
			Depth:    int(cItem.depth),
			Weight:   int(cItem.weight),
			Rotation: Rotation(cItem.rotation),
		}
	}

	// Call Go function
	result := CalculateLookahead(
		prevItem,
		goItems,
		int(widthLeft),
		int(lengthLeft),
		int(depthLeft),
		int(rowLength),
		int(maxLookahead),
	)

	return C.int(result)
}

//export GetBestOrientationFFI
func GetBestOrientationFFI(
	item *C.CItem,
	nextItems *C.CItem,
	nextItemCount C.int,
	widthLeft, lengthLeft, depthLeft C.int,
	rowLength C.int,
	packedWeight C.int,
	box *C.CBox,
	resultOrientation *C.COrientatedItem,
) C.int {
	// Convert item
	goItem := &Item{
		Width:    int(item.width),
		Length:   int(item.length),
		Depth:    int(item.depth),
		Weight:   int(item.weight),
		Rotation: Rotation(item.rotation),
	}

	// Convert box
	goBox := &Box{
		InnerWidth:  int(box.innerWidth),
		InnerLength: int(box.innerLength),
		InnerDepth:  int(box.innerDepth),
		MaxWeight:   int(box.maxWeight),
	}

	// Get possible orientations
	orientations := GetPossibleOrientations(
		goItem,
		nil,
		int(widthLeft),
		int(lengthLeft),
		int(depthLeft),
		int(packedWeight),
		goBox.MaxWeight,
	)

	if len(orientations) == 0 {
		return 0 // No valid orientation
	}

	// Filter by stability
	usableOrientations := GetUsableOrientations(orientations, goBox.InnerDepth)
	if len(usableOrientations) == 0 {
		return 0
	}

	// Convert next items for lookahead
	var goNextItems []*Item
	if nextItemCount > 0 {
		nextItemsSlice := unsafe.Slice(nextItems, nextItemCount)
		goNextItems = make([]*Item, nextItemCount)
		for i, cItem := range nextItemsSlice {
			goNextItems[i] = &Item{
				Width:    int(cItem.width),
				Length:   int(cItem.length),
				Depth:    int(cItem.depth),
				Weight:   int(cItem.weight),
				Rotation: Rotation(cItem.rotation),
			}
		}
	}

	// Sort by lookahead if we have next items
	if len(goNextItems) > 0 {
		usableOrientations = SortOrientationsByLookahead(
			usableOrientations,
			goNextItems,
			int(widthLeft),
			int(lengthLeft),
			int(depthLeft),
			int(rowLength),
			8, // Max lookahead depth
		)
	}

	// Return best orientation
	best := usableOrientations[0]
	resultOrientation.width = C.int(best.Width)
	resultOrientation.length = C.int(best.Length)
	resultOrientation.depth = C.int(best.Depth)
	resultOrientation.surfaceFootprint = C.int(best.SurfaceFootprint)

	return 1 // Success
}

//export ClearCacheFFI
func ClearCacheFFI() {
	ClearLookaheadCache()
}

//export GetCacheSizeFFI
func GetCacheSizeFFI() C.int {
	return C.int(len(lookaheadCache))
}

func main() {
	// Required for building shared library
}
