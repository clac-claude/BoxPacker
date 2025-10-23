package tests

// LimitedSupplyTestBox is a test implementation of LimitedSupplyBox interface
type LimitedSupplyTestBox struct {
	Reference   string
	OuterWidth  int
	OuterLength int
	OuterDepth  int
	EmptyWeight int
	InnerWidth  int
	InnerLength int
	InnerDepth  int
	MaxWeight   int
	Quantity    int
}

// NewLimitedSupplyTestBox creates a new LimitedSupplyTestBox
func NewLimitedSupplyTestBox(reference string, outerWidth, outerLength, outerDepth, emptyWeight, innerWidth, innerLength, innerDepth, maxWeight, quantity int) *LimitedSupplyTestBox {
	return &LimitedSupplyTestBox{
		Reference:   reference,
		OuterWidth:  outerWidth,
		OuterLength: outerLength,
		OuterDepth:  outerDepth,
		EmptyWeight: emptyWeight,
		InnerWidth:  innerWidth,
		InnerLength: innerLength,
		InnerDepth:  innerDepth,
		MaxWeight:   maxWeight,
		Quantity:    quantity,
	}
}

func (b *LimitedSupplyTestBox) GetReference() string {
	return b.Reference
}

func (b *LimitedSupplyTestBox) GetOuterWidth() int {
	return b.OuterWidth
}

func (b *LimitedSupplyTestBox) GetOuterLength() int {
	return b.OuterLength
}

func (b *LimitedSupplyTestBox) GetOuterDepth() int {
	return b.OuterDepth
}

func (b *LimitedSupplyTestBox) GetEmptyWeight() int {
	return b.EmptyWeight
}

func (b *LimitedSupplyTestBox) GetInnerWidth() int {
	return b.InnerWidth
}

func (b *LimitedSupplyTestBox) GetInnerLength() int {
	return b.InnerLength
}

func (b *LimitedSupplyTestBox) GetInnerDepth() int {
	return b.InnerDepth
}

func (b *LimitedSupplyTestBox) GetMaxWeight() int {
	return b.MaxWeight
}

func (b *LimitedSupplyTestBox) GetQuantityAvailable() int {
	return b.Quantity
}
