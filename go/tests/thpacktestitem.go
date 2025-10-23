package tests

// THPackTestItem is a test item with constrained placement based on allowed vertical orientations
type THPackTestItem struct {
	Description          string
	Width                int
	WidthAllowedVertical bool
	Length               int
	LengthAllowedVertical bool
	Depth                int
	DepthAllowedVertical bool
}

// NewTHPackTestItem creates a new THPackTestItem
func NewTHPackTestItem(description string, width int, widthAllowedVertical bool, length int, lengthAllowedVertical bool, depth int, depthAllowedVertical bool) *THPackTestItem {
	return &THPackTestItem{
		Description:           description,
		Width:                 width,
		WidthAllowedVertical:  widthAllowedVertical,
		Length:                length,
		LengthAllowedVertical: lengthAllowedVertical,
		Depth:                 depth,
		DepthAllowedVertical:  depthAllowedVertical,
	}
}

func (i *THPackTestItem) GetDescription() string {
	return i.Description
}

func (i *THPackTestItem) GetWidth() int {
	return i.Width
}

func (i *THPackTestItem) GetLength() int {
	return i.Length
}

func (i *THPackTestItem) GetDepth() int {
	return i.Depth
}

func (i *THPackTestItem) GetWeight() int {
	return 0
}

func (i *THPackTestItem) GetAllowedRotation() Rotation {
	if !i.WidthAllowedVertical && !i.LengthAllowedVertical && i.DepthAllowedVertical {
		return KeepFlat
	}
	return BestFit
}

// CanBePacked checks if the item can be packed at the proposed position
func (i *THPackTestItem) CanBePacked(packedBox interface{}, proposedX, proposedY, proposedZ, width, length, depth int) bool {
	ok := false
	if depth == i.Width {
		ok = ok || i.WidthAllowedVertical
	}
	if depth == i.Length {
		ok = ok || i.LengthAllowedVertical
	}
	if depth == i.Depth {
		ok = ok || i.DepthAllowedVertical
	}
	return ok
}
