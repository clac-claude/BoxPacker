package tests

import (
	"testing"
)

// BoxList represents a list of available boxes
type BoxList struct {
	boxes []interface{}
}

// NewBoxList creates a new BoxList
func NewBoxList() *BoxList {
	return &BoxList{
		boxes: make([]interface{}, 0),
	}
}

func (l *BoxList) Insert(box interface{}) {
	l.boxes = append(l.boxes, box)
}

func (l *BoxList) Count() int {
	return len(l.boxes)
}

// Simplified tests
func TestBoxListBasics(t *testing.T) {
	boxList := NewBoxList()

	if boxList.Count() != 0 {
		t.Errorf("Expected count 0, got %d", boxList.Count())
	}

	box := NewTestBox("Box", 10, 10, 10, 0, 10, 10, 10, 100)
	boxList.Insert(box)

	if boxList.Count() != 1 {
		t.Errorf("Expected count 1, got %d", boxList.Count())
	}
}
