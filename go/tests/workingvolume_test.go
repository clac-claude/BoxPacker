package tests

import (
	"encoding/json"
	"testing"
)

// WorkingVolume represents a working volume for packing
type WorkingVolume struct {
	innerWidth  int
	innerLength int
	innerDepth  int
	maxWeight   int
}

// NewWorkingVolume creates a new WorkingVolume
func NewWorkingVolume(innerWidth, innerLength, innerDepth, maxWeight int) *WorkingVolume {
	return &WorkingVolume{
		innerWidth:  innerWidth,
		innerLength: innerLength,
		innerDepth:  innerDepth,
		maxWeight:   maxWeight,
	}
}

func (v *WorkingVolume) GetInnerWidth() int {
	return v.innerWidth
}

func (v *WorkingVolume) GetOuterWidth() int {
	return v.innerWidth
}

func (v *WorkingVolume) GetInnerLength() int {
	return v.innerLength
}

func (v *WorkingVolume) GetOuterLength() int {
	return v.innerLength
}

func (v *WorkingVolume) GetInnerDepth() int {
	return v.innerDepth
}

func (v *WorkingVolume) GetOuterDepth() int {
	return v.innerDepth
}

func (v *WorkingVolume) GetEmptyWeight() int {
	return 0
}

func (v *WorkingVolume) GetMaxWeight() int {
	return v.maxWeight
}

func (v *WorkingVolume) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"reference": "working-volume",
		"width":     v.innerWidth,
		"length":    v.innerLength,
		"depth":     v.innerDepth,
		"maxWeight": v.maxWeight,
	})
}

func TestDimensions(t *testing.T) {
	volume := NewWorkingVolume(1, 2, 3, 4)

	if volume.GetInnerWidth() != 1 {
		t.Errorf("Expected innerWidth 1, got %d", volume.GetInnerWidth())
	}
	if volume.GetOuterWidth() != 1 {
		t.Errorf("Expected outerWidth 1, got %d", volume.GetOuterWidth())
	}
	if volume.GetInnerLength() != 2 {
		t.Errorf("Expected innerLength 2, got %d", volume.GetInnerLength())
	}
	if volume.GetOuterLength() != 2 {
		t.Errorf("Expected outerLength 2, got %d", volume.GetOuterLength())
	}
	if volume.GetInnerDepth() != 3 {
		t.Errorf("Expected innerDepth 3, got %d", volume.GetInnerDepth())
	}
	if volume.GetOuterDepth() != 3 {
		t.Errorf("Expected outerDepth 3, got %d", volume.GetOuterDepth())
	}
	if volume.GetEmptyWeight() != 0 {
		t.Errorf("Expected emptyWeight 0, got %d", volume.GetEmptyWeight())
	}
	if volume.GetMaxWeight() != 4 {
		t.Errorf("Expected maxWeight 4, got %d", volume.GetMaxWeight())
	}
}

func TestSerialize(t *testing.T) {
	volume := NewWorkingVolume(1, 2, 3, 4)

	data, err := json.Marshal(volume)
	if err != nil {
		t.Fatalf("Failed to marshal WorkingVolume: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	expectedKeys := []string{"reference", "width", "length", "depth", "maxWeight"}
	for _, key := range expectedKeys {
		if _, ok := result[key]; !ok {
			t.Errorf("Expected key %s not found in serialized data", key)
		}
	}
}
