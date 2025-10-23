package tests

import (
	"encoding/json"
	"fmt"
	"testing"
)

// OrientatedItem represents an item with a specific orientation
type OrientatedItem struct {
	Item   interface{}
	Width  int
	Length int
	Depth  int
}

// NewOrientatedItem creates a new OrientatedItem
func NewOrientatedItem(item interface{}, width, length, depth int) *OrientatedItem {
	return &OrientatedItem{
		Item:   item,
		Width:  width,
		Length: length,
		Depth:  depth,
	}
}

func (o *OrientatedItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"item":   o.Item,
		"width":  o.Width,
		"length": o.Length,
		"depth":  o.Depth,
	})
}

func (o *OrientatedItem) String() string {
	return fmt.Sprintf("%d|%d|%d", o.Width, o.Length, o.Depth)
}

func TestOrientatedItemSerialize(t *testing.T) {
	item := NewOrientatedItem(NewTestItem("Test", 1, 2, 3, 4, BestFit), 1, 2, 3)

	data, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("Failed to marshal OrientatedItem: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	expectedKeys := []string{"item", "width", "length", "depth"}
	for _, key := range expectedKeys {
		if _, ok := result[key]; !ok {
			t.Errorf("Expected key %s not found in serialized data", key)
		}
	}

	expected := "1|2|3"
	got := item.String()
	if got != expected {
		t.Errorf("String() = %s, want %s", got, expected)
	}
}
