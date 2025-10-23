package tests

import "encoding/json"

// Rotation represents allowed rotation for items
type Rotation int

const (
	BestFit  Rotation = 6
	KeepFlat Rotation = 3
)

// TestItem is a test implementation of the Item interface
type TestItem struct {
	Description     string
	Width           int
	Length          int
	Depth           int
	Weight          int
	AllowedRotation Rotation

	jsonSerializeOverride interface{}
}

// NewTestItem creates a new TestItem
func NewTestItem(description string, width, length, depth, weight int, allowedRotation Rotation) *TestItem {
	return &TestItem{
		Description:     description,
		Width:           width,
		Length:          length,
		Depth:           depth,
		Weight:          weight,
		AllowedRotation: allowedRotation,
	}
}

func (i *TestItem) GetDescription() string {
	return i.Description
}

func (i *TestItem) GetWidth() int {
	return i.Width
}

func (i *TestItem) GetLength() int {
	return i.Length
}

func (i *TestItem) GetDepth() int {
	return i.Depth
}

func (i *TestItem) GetWeight() int {
	return i.Weight
}

func (i *TestItem) GetAllowedRotation() Rotation {
	return i.AllowedRotation
}

func (i *TestItem) MarshalJSON() ([]byte, error) {
	if i.jsonSerializeOverride != nil {
		return json.Marshal(i.jsonSerializeOverride)
	}

	return json.Marshal(map[string]interface{}{
		"description":     i.Description,
		"width":           i.Width,
		"length":          i.Length,
		"depth":           i.Depth,
		"weight":          i.Weight,
		"allowedRotation": i.AllowedRotation,
	})
}

func (i *TestItem) SetJsonSerializeOverride(override interface{}) {
	i.jsonSerializeOverride = override
}
