package tests

import (
	"testing"
)

// NoBoxesAvailableException represents an exception when no boxes are available
type NoBoxesAvailableException struct {
	message        string
	affectedItems  []interface{}
}

// NewNoBoxesAvailableException creates a new exception
func NewNoBoxesAvailableException(message string, items []interface{}) *NoBoxesAvailableException {
	return &NoBoxesAvailableException{
		message:       message,
		affectedItems: items,
	}
}

func (e *NoBoxesAvailableException) Error() string {
	return e.message
}

func (e *NoBoxesAvailableException) GetAffectedItems() []interface{} {
	return e.affectedItems
}

// TestCanGetItem tests that the offending items can be retrieved from the exception
func TestCanGetItem(t *testing.T) {
	item1 := NewTestItem("Item 1", 2500, 2500, 20, 2000, BestFit)
	item2 := NewTestItem("Item 2", 2500, 2500, 20, 2000, BestFit)

	itemList := []interface{}{item1, item2}
	exception := NewNoBoxesAvailableException("Just testing...", itemList)

	affectedItems := exception.GetAffectedItems()
	if len(affectedItems) != 2 {
		t.Errorf("Expected 2 affected items, got %d", len(affectedItems))
	}

	if affectedItems[0] != item1 {
		t.Errorf("Expected first item to be item1")
	}

	if affectedItems[1] != item2 {
		t.Errorf("Expected second item to be item2")
	}
}
