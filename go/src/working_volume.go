package boxpacker

import "fmt"

// Box packing (3D bin packing, knapsack problem).
//
// @author Doug Wright

// WorkingVolume - Internal working volume used for packing calculations.
type WorkingVolume struct {
	width     int
	length    int
	depth     int
	maxWeight int
}

// NewWorkingVolume creates a new WorkingVolume
func NewWorkingVolume(width, length, depth, maxWeight int) *WorkingVolume {
	return &WorkingVolume{
		width:     width,
		length:    length,
		depth:     depth,
		maxWeight: maxWeight,
	}
}

// GetReference returns the reference string for this working volume
func (wv *WorkingVolume) GetReference() string {
	return fmt.Sprintf("Working Volume %dx%dx%d", wv.width, wv.length, wv.depth)
}

// GetOuterWidth returns the outer width (same as inner for working volume)
func (wv *WorkingVolume) GetOuterWidth() int {
	return wv.width
}

// GetOuterLength returns the outer length (same as inner for working volume)
func (wv *WorkingVolume) GetOuterLength() int {
	return wv.length
}

// GetOuterDepth returns the outer depth (same as inner for working volume)
func (wv *WorkingVolume) GetOuterDepth() int {
	return wv.depth
}

// GetEmptyWeight returns the empty weight (always 0 for working volume)
func (wv *WorkingVolume) GetEmptyWeight() int {
	return 0
}

// GetInnerWidth returns the inner width
func (wv *WorkingVolume) GetInnerWidth() int {
	return wv.width
}

// GetInnerLength returns the inner length
func (wv *WorkingVolume) GetInnerLength() int {
	return wv.length
}

// GetInnerDepth returns the inner depth
func (wv *WorkingVolume) GetInnerDepth() int {
	return wv.depth
}

// GetMaxWeight returns the maximum weight capacity
func (wv *WorkingVolume) GetMaxWeight() int {
	return wv.maxWeight
}
