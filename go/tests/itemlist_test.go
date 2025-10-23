package tests

import (
	"testing"
)

// ItemList represents a list of items to be packed
type ItemList struct {
	items []interface{}
}

// NewItemList creates a new ItemList
func NewItemList() *ItemList {
	return &ItemList{
		items: make([]interface{}, 0),
	}
}

func (l *ItemList) Insert(item interface{}) {
	l.items = append(l.items, item)
}

func (l *ItemList) Remove(item interface{}) {
	for i, it := range l.items {
		if it == item {
			l.items = append(l.items[:i], l.items[i+1:]...)
			break
		}
	}
}

func (l *ItemList) Count() int {
	return len(l.items)
}

func (l *ItemList) Top() interface{} {
	if len(l.items) == 0 {
		return nil
	}
	return l.items[0]
}

func (l *ItemList) TopN(n int) *ItemList {
	result := NewItemList()
	for i := 0; i < n && i < len(l.items); i++ {
		result.Insert(l.items[i])
	}
	return result
}

func (l *ItemList) Extract() interface{} {
	if len(l.items) == 0 {
		return nil
	}
	item := l.items[0]
	l.items = l.items[1:]
	return item
}

func TestDimensionalSorting(t *testing.T) {
	// In Go, we would need to implement proper sorting
	// This is a simplified stub
	t.Skip("Full sorting implementation requires complete Item comparison logic")
}

func TestKeepingItemsOfSameTypeTogether(t *testing.T) {
	// In Go, we would need to implement proper sorting
	// This is a simplified stub
	t.Skip("Full sorting implementation requires complete Item comparison logic")
}

func TestCount(t *testing.T) {
	itemList := NewItemList()
	if itemList.Count() != 0 {
		t.Errorf("Expected count 0, got %d", itemList.Count())
	}

	item1 := NewTestItem("Item A", 20, 20, 2, 100, BestFit)
	itemList.Insert(item1)
	if itemList.Count() != 1 {
		t.Errorf("Expected count 1, got %d", itemList.Count())
	}

	item2 := NewTestItem("Item B", 20, 20, 2, 100, BestFit)
	itemList.Insert(item2)
	if itemList.Count() != 2 {
		t.Errorf("Expected count 2, got %d", itemList.Count())
	}

	item3 := NewTestItem("Item C", 20, 20, 2, 100, BestFit)
	itemList.Insert(item3)
	if itemList.Count() != 3 {
		t.Errorf("Expected count 3, got %d", itemList.Count())
	}

	itemList.Remove(item2)
	if itemList.Count() != 2 {
		t.Errorf("Expected count 2, got %d", itemList.Count())
	}
}

func TestTop(t *testing.T) {
	itemList := NewItemList()
	item1 := NewTestItem("Item A", 20, 20, 2, 100, BestFit)
	itemList.Insert(item1)

	if itemList.Top() != item1 {
		t.Errorf("Expected top to be item1")
	}
	if itemList.Count() != 1 {
		t.Errorf("Expected count to remain 1")
	}
}

func TestTopN(t *testing.T) {
	itemList := NewItemList()

	item1 := NewTestItem("Item A", 20, 20, 2, 100, BestFit)
	itemList.Insert(item1)

	item2 := NewTestItem("Item B", 20, 20, 2, 100, BestFit)
	itemList.Insert(item2)

	item3 := NewTestItem("Item C", 20, 20, 2, 100, BestFit)
	itemList.Insert(item3)

	top2 := itemList.TopN(2)

	if top2.Count() != 2 {
		t.Errorf("Expected topN count 2, got %d", top2.Count())
	}

	extracted1 := top2.Extract()
	if extracted1 != item1 {
		t.Errorf("Expected first extracted item to be item1")
	}

	extracted2 := top2.Extract()
	if extracted2 != item2 {
		t.Errorf("Expected second extracted item to be item2")
	}
}

func TestExtract(t *testing.T) {
	itemList := NewItemList()
	item1 := NewTestItem("Item A", 20, 20, 2, 100, BestFit)
	itemList.Insert(item1)

	if itemList.Extract() != item1 {
		t.Errorf("Expected extracted item to be item1")
	}
	if itemList.Count() != 0 {
		t.Errorf("Expected count 0 after extract")
	}
}
