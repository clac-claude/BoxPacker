package boxpacker

// Box packing (3D bin packing, knapsack problem).
//
// @author Doug Wright

// DefaultBoxSorter - Default implementation of BoxSorter
type DefaultBoxSorter struct{}

// Compare compares two boxes
func (dbs *DefaultBoxSorter) Compare(boxA Box, boxB Box) int {
	boxAVolume := boxA.GetInnerWidth() * boxA.GetInnerLength() * boxA.GetInnerDepth()
	boxBVolume := boxB.GetInnerWidth() * boxB.GetInnerLength() * boxB.GetInnerDepth()

	// try smallest box first
	if boxAVolume < boxBVolume {
		return -1
	}
	if boxAVolume > boxBVolume {
		return 1
	}

	// with smallest empty weight
	if boxA.GetEmptyWeight() < boxB.GetEmptyWeight() {
		return -1
	}
	if boxA.GetEmptyWeight() > boxB.GetEmptyWeight() {
		return 1
	}

	// maximum weight capacity as fallback decider
	boxACapacity := boxA.GetMaxWeight() - boxA.GetEmptyWeight()
	boxBCapacity := boxB.GetMaxWeight() - boxB.GetEmptyWeight()

	if boxACapacity < boxBCapacity {
		return -1
	}
	if boxACapacity > boxBCapacity {
		return 1
	}

	return 0
}
