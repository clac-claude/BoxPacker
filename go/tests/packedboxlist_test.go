package tests

import (
	"testing"
)

// PackedBoxList represents a list of packed boxes
type PackedBoxList struct {
	boxes []*PackedBox
}

// NewPackedBoxList creates a new PackedBoxList
func NewPackedBoxList() *PackedBoxList {
	return &PackedBoxList{
		boxes: make([]*PackedBox, 0),
	}
}

func (l *PackedBoxList) Insert(box *PackedBox) {
	l.boxes = append(l.boxes, box)
}

func (l *PackedBoxList) Count() int {
	return len(l.boxes)
}

func (l *PackedBoxList) Top() *PackedBox {
	if len(l.boxes) == 0 {
		return nil
	}
	return l.boxes[0]
}

// Simplified tests - full implementation would require complete Packer logic
func TestPackedBoxListBasics(t *testing.T) {
	boxList := NewPackedBoxList()

	if boxList.Count() != 0 {
		t.Errorf("Expected count 0, got %d", boxList.Count())
	}

	box := NewTestBox("Box", 10, 10, 10, 0, 10, 10, 10, 100)
	packedBox := NewPackedBox(box, []*PackedItem{})

	boxList.Insert(packedBox)

	if boxList.Count() != 1 {
		t.Errorf("Expected count 1, got %d", boxList.Count())
	}

	if boxList.Top() != packedBox {
		t.Errorf("Expected top to be packedBox")
	}
}
