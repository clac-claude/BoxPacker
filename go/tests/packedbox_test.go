package tests

import (
	"testing"
)

// PackedBox represents a box with packed items
type PackedBox struct {
	Box   interface{}
	Items []*PackedItem
}

// NewPackedBox creates a new PackedBox
func NewPackedBox(box interface{}, items []*PackedItem) *PackedBox {
	return &PackedBox{
		Box:   box,
		Items: items,
	}
}

func (p *PackedBox) GetWeight() int {
	totalWeight := 0
	if testBox, ok := p.Box.(*TestBox); ok {
		totalWeight += testBox.EmptyWeight
	}
	for _, item := range p.Items {
		if testItem, ok := item.Item.(*TestItem); ok {
			totalWeight += testItem.Weight
		}
	}
	return totalWeight
}

func (p *PackedBox) GetItemWeight() int {
	totalWeight := 0
	for _, item := range p.Items {
		if testItem, ok := item.Item.(*TestItem); ok {
			totalWeight += testItem.Weight
		}
	}
	return totalWeight
}

func (p *PackedBox) GetRemainingWidth() int {
	if testBox, ok := p.Box.(*TestBox); ok {
		usedWidth := 0
		for _, item := range p.Items {
			if item.X+item.Width > usedWidth {
				usedWidth = item.X + item.Width
			}
		}
		return testBox.InnerWidth - usedWidth
	}
	return 0
}

func (p *PackedBox) GetRemainingLength() int {
	if testBox, ok := p.Box.(*TestBox); ok {
		usedLength := 0
		for _, item := range p.Items {
			if item.Y+item.Length > usedLength {
				usedLength = item.Y + item.Length
			}
		}
		return testBox.InnerLength - usedLength
	}
	return 0
}

func (p *PackedBox) GetRemainingDepth() int {
	if testBox, ok := p.Box.(*TestBox); ok {
		usedDepth := 0
		for _, item := range p.Items {
			if item.Z+item.Depth > usedDepth {
				usedDepth = item.Z + item.Depth
			}
		}
		return testBox.InnerDepth - usedDepth
	}
	return 0
}

func (p *PackedBox) GetRemainingWeight() int {
	if testBox, ok := p.Box.(*TestBox); ok {
		return testBox.MaxWeight - p.GetWeight()
	}
	return 0
}

func (p *PackedBox) GetInnerVolume() int {
	if testBox, ok := p.Box.(*TestBox); ok {
		return testBox.InnerWidth * testBox.InnerLength * testBox.InnerDepth
	}
	return 0
}

func (p *PackedBox) GetUsedVolume() int {
	usedVolume := 0
	for _, item := range p.Items {
		usedVolume += item.Width * item.Length * item.Depth
	}
	return usedVolume
}

func (p *PackedBox) GetUnusedVolume() int {
	return p.GetInnerVolume() - p.GetUsedVolume()
}

func (p *PackedBox) GetVolumeUtilisation() int {
	innerVolume := p.GetInnerVolume()
	if innerVolume == 0 {
		return 0
	}
	return (p.GetUsedVolume() * 100) / innerVolume
}

func TestPackedBoxGetters(t *testing.T) {
	box := NewTestBox("Box", 370, 375, 60, 140, 364, 374, 40, 3000)
	item := NewOrientatedItem(NewTestItem("Item", 230, 330, 6, 320, BestFit), 230, 330, 6)

	packedItem := NewPackedItem(item, 0, 0, 0, 230, 330, 6)
	packedItems := []*PackedItem{packedItem}

	packedBox := NewPackedBox(box, packedItems)

	if packedBox.Box != box {
		t.Errorf("Expected box to match")
	}

	if packedBox.GetWeight() != 460 {
		t.Errorf("Expected weight 460, got %d", packedBox.GetWeight())
	}

	if packedBox.GetRemainingWidth() != 134 {
		t.Errorf("Expected remaining width 134, got %d", packedBox.GetRemainingWidth())
	}

	if packedBox.GetRemainingLength() != 44 {
		t.Errorf("Expected remaining length 44, got %d", packedBox.GetRemainingLength())
	}

	if packedBox.GetRemainingDepth() != 34 {
		t.Errorf("Expected remaining depth 34, got %d", packedBox.GetRemainingDepth())
	}

	if packedBox.GetRemainingWeight() != 2540 {
		t.Errorf("Expected remaining weight 2540, got %d", packedBox.GetRemainingWeight())
	}

	if packedBox.GetInnerVolume() != 5445440 {
		t.Errorf("Expected inner volume 5445440, got %d", packedBox.GetInnerVolume())
	}
}

func TestVolumeUtilisation(t *testing.T) {
	box := NewTestBox("Box", 10, 10, 20, 10, 10, 10, 20, 10)
	item := NewOrientatedItem(NewTestItem("Item", 4, 10, 10, 10, BestFit), 4, 10, 10)

	packedItem := NewPackedItem(item, 0, 0, 0, 4, 10, 10)
	boxItems := []*PackedItem{packedItem}

	packedBox := NewPackedBox(box, boxItems)

	if packedBox.GetUsedVolume() != 400 {
		t.Errorf("Expected used volume 400, got %d", packedBox.GetUsedVolume())
	}

	if packedBox.GetUnusedVolume() != 1600 {
		t.Errorf("Expected unused volume 1600, got %d", packedBox.GetUnusedVolume())
	}

	if packedBox.GetVolumeUtilisation() != 20 {
		t.Errorf("Expected volume utilisation 20%%, got %d%%", packedBox.GetVolumeUtilisation())
	}
}

func TestWeightCalc(t *testing.T) {
	box := NewTestBox("Box", 10, 10, 20, 10, 10, 10, 20, 10)
	item := NewOrientatedItem(NewTestItem("Item", 4, 10, 10, 10, BestFit), 4, 10, 10)

	packedItem := NewPackedItem(item, 0, 0, 0, 4, 10, 10)
	boxItems := []*PackedItem{packedItem}

	packedBox := NewPackedBox(box, boxItems)

	if packedBox.GetItemWeight() != 10 {
		t.Errorf("Expected item weight 10, got %d", packedBox.GetItemWeight())
	}
}
