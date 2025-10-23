package tests

import "encoding/json"

// TestBox is a test implementation of the Box interface
type TestBox struct {
	Reference   string
	OuterWidth  int
	OuterLength int
	OuterDepth  int
	EmptyWeight int
	InnerWidth  int
	InnerLength int
	InnerDepth  int
	MaxWeight   int

	jsonSerializeOverride interface{}
}

// NewTestBox creates a new TestBox
func NewTestBox(reference string, outerWidth, outerLength, outerDepth, emptyWeight, innerWidth, innerLength, innerDepth, maxWeight int) *TestBox {
	return &TestBox{
		Reference:   reference,
		OuterWidth:  outerWidth,
		OuterLength: outerLength,
		OuterDepth:  outerDepth,
		EmptyWeight: emptyWeight,
		InnerWidth:  innerWidth,
		InnerLength: innerLength,
		InnerDepth:  innerDepth,
		MaxWeight:   maxWeight,
	}
}

func (b *TestBox) GetReference() string {
	return b.Reference
}

func (b *TestBox) GetOuterWidth() int {
	return b.OuterWidth
}

func (b *TestBox) GetOuterLength() int {
	return b.OuterLength
}

func (b *TestBox) GetOuterDepth() int {
	return b.OuterDepth
}

func (b *TestBox) GetEmptyWeight() int {
	return b.EmptyWeight
}

func (b *TestBox) GetInnerWidth() int {
	return b.InnerWidth
}

func (b *TestBox) GetInnerLength() int {
	return b.InnerLength
}

func (b *TestBox) GetInnerDepth() int {
	return b.InnerDepth
}

func (b *TestBox) GetMaxWeight() int {
	return b.MaxWeight
}

func (b *TestBox) MarshalJSON() ([]byte, error) {
	if b.jsonSerializeOverride != nil {
		return json.Marshal(b.jsonSerializeOverride)
	}

	return json.Marshal(map[string]interface{}{
		"reference":   b.Reference,
		"innerWidth":  b.InnerWidth,
		"innerLength": b.InnerLength,
		"innerDepth":  b.InnerDepth,
		"emptyWeight": b.EmptyWeight,
		"maxWeight":   b.MaxWeight,
	})
}

func (b *TestBox) SetJsonSerializeOverride(override interface{}) {
	b.jsonSerializeOverride = override
}
