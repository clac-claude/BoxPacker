package main

import (
	"testing"
)

func TestGenerateOrientations(t *testing.T) {
	item := &Item{
		Width:    10,
		Length:   20,
		Depth:    5,
		Weight:   100,
		Rotation: RotationBestFit,
	}

	orientations := GenerateOrientations(item, nil)

	// Should have 6 orientations for BestFit
	if len(orientations) != 6 {
		t.Errorf("Expected 6 orientations, got %d", len(orientations))
	}

	// Check that all dimensions are present in some orientation
	foundBaseOrientation := false
	for _, o := range orientations {
		if o.Width == 10 && o.Length == 20 && o.Depth == 5 {
			foundBaseOrientation = true
			break
		}
	}

	if !foundBaseOrientation {
		t.Error("Base orientation not found")
	}
}

func TestGenerateOrientationsNever(t *testing.T) {
	item := &Item{
		Width:    10,
		Length:   20,
		Depth:    5,
		Weight:   100,
		Rotation: RotationNever,
	}

	orientations := GenerateOrientations(item, nil)

	// Should have only 1 orientation for Never
	if len(orientations) != 1 {
		t.Errorf("Expected 1 orientation, got %d", len(orientations))
	}

	if orientations[0].Width != 10 || orientations[0].Length != 20 || orientations[0].Depth != 5 {
		t.Error("Wrong orientation returned")
	}
}

func TestGetPossibleOrientations(t *testing.T) {
	item := &Item{
		Width:    10,
		Length:   20,
		Depth:    5,
		Weight:   100,
		Rotation: RotationBestFit,
	}

	// Space that fits the item
	orientations := GetPossibleOrientations(item, nil, 25, 25, 10, 0, 1000)

	if len(orientations) == 0 {
		t.Error("Expected at least one possible orientation")
	}

	// Space that doesn't fit
	orientations = GetPossibleOrientations(item, nil, 5, 5, 5, 0, 1000)

	if len(orientations) != 0 {
		t.Errorf("Expected no orientations for too small space, got %d", len(orientations))
	}
}

func TestGetPossibleOrientationsTooHeavy(t *testing.T) {
	item := &Item{
		Width:    10,
		Length:   20,
		Depth:    5,
		Weight:   100,
		Rotation: RotationBestFit,
	}

	// Item is too heavy
	orientations := GetPossibleOrientations(item, nil, 25, 25, 10, 50, 100)

	if len(orientations) != 0 {
		t.Error("Expected no orientations for too heavy item")
	}
}

func TestSimplePacker(t *testing.T) {
	box := &Box{
		InnerWidth:  100,
		InnerLength: 100,
		InnerDepth:  100,
		MaxWeight:   1000,
	}

	items := []*Item{
		{Width: 10, Length: 10, Depth: 10, Weight: 10, Rotation: RotationBestFit},
		{Width: 10, Length: 10, Depth: 10, Weight: 10, Rotation: RotationBestFit},
		{Width: 10, Length: 10, Depth: 10, Weight: 10, Rotation: RotationBestFit},
	}

	packer := &SimplePacker{
		Box:   box,
		Items: items,
	}

	packed := packer.PackSimple()

	// Should pack all 3 items
	if len(packed) != 3 {
		t.Errorf("Expected 3 packed items, got %d", len(packed))
	}

	// Check positions are valid
	for _, p := range packed {
		if p.X < 0 || p.Y < 0 || p.Z < 0 {
			t.Error("Invalid negative position")
		}
		if p.X+p.Width > box.InnerWidth {
			t.Error("Item exceeds box width")
		}
		if p.Y+p.Length > box.InnerLength {
			t.Error("Item exceeds box length")
		}
		if p.Z+p.Depth > box.InnerDepth {
			t.Error("Item exceeds box depth")
		}
	}
}

func TestCalculateLookahead(t *testing.T) {
	prevItem := &OrientatedItem{
		Width:  10,
		Length: 10,
		Depth:  10,
	}

	nextItems := []*Item{
		{Width: 5, Length: 5, Depth: 5, Weight: 5, Rotation: RotationBestFit},
		{Width: 5, Length: 5, Depth: 5, Weight: 5, Rotation: RotationBestFit},
		{Width: 5, Length: 5, Depth: 5, Weight: 5, Rotation: RotationBestFit},
	}

	// Clear cache before test
	ClearLookaheadCache()

	result := CalculateLookahead(prevItem, nextItems, 100, 100, 100, 10, 8)

	// Should be able to pack some items
	if result < 0 {
		t.Error("Lookahead returned negative result")
	}

	// Test caching - second call should be instant
	result2 := CalculateLookahead(prevItem, nextItems, 100, 100, 100, 10, 8)

	if result != result2 {
		t.Error("Cached result differs from original")
	}
}

func TestSortOrientationsByLookahead(t *testing.T) {
	item := &Item{
		Width:    10,
		Length:   20,
		Depth:    5,
		Weight:   10,
		Rotation: RotationBestFit,
	}

	orientations := GenerateOrientations(item, nil)

	nextItems := []*Item{
		{Width: 5, Length: 5, Depth: 5, Weight: 5, Rotation: RotationBestFit},
		{Width: 5, Length: 5, Depth: 5, Weight: 5, Rotation: RotationBestFit},
	}

	sorted := SortOrientationsByLookahead(orientations, nextItems, 100, 100, 100, 0, 8)

	if len(sorted) != len(orientations) {
		t.Error("Sorting changed number of orientations")
	}

	// Verify it's actually sorted (higher lookahead first)
	for i := 0; i < len(sorted)-1; i++ {
		scoreA := CalculateLookahead(sorted[i], nextItems, 100, 100, 100, 0, 8)
		scoreB := CalculateLookahead(sorted[i+1], nextItems, 100, 100, 100, 0, 8)

		if scoreA < scoreB {
			t.Errorf("Sorting failed: orientation %d (score %d) should be before %d (score %d)",
				i, scoreA, i+1, scoreB)
		}
	}
}

func TestIsSameDimensions(t *testing.T) {
	itemA := &Item{Width: 10, Length: 20, Depth: 5, Weight: 100, Rotation: RotationBestFit}
	itemB := &Item{Width: 20, Length: 10, Depth: 5, Weight: 50, Rotation: RotationNever} // Different order, same dims
	itemC := &Item{Width: 10, Length: 20, Depth: 10, Weight: 100, Rotation: RotationBestFit}

	if !IsSameDimensions(itemA, itemB) {
		t.Error("Items with same dimensions (different order) should be considered same")
	}

	if IsSameDimensions(itemA, itemC) {
		t.Error("Items with different dimensions should not be considered same")
	}

	if !IsSameDimensions(itemA, itemA) {
		t.Error("Item should be same as itself")
	}
}

func TestClearLookaheadCache(t *testing.T) {
	prevItem := &OrientatedItem{Width: 10, Length: 10, Depth: 10}
	nextItems := []*Item{
		{Width: 5, Length: 5, Depth: 5, Weight: 5, Rotation: RotationBestFit},
	}

	ClearLookaheadCache()

	// Add something to cache
	CalculateLookahead(prevItem, nextItems, 100, 100, 100, 10, 8)

	if len(lookaheadCache) == 0 {
		t.Error("Cache should not be empty after calculation")
	}

	ClearLookaheadCache()

	if len(lookaheadCache) != 0 {
		t.Error("Cache should be empty after clear")
	}
}

// Benchmarks

func BenchmarkGenerateOrientations(b *testing.B) {
	item := &Item{
		Width:    10,
		Length:   20,
		Depth:    5,
		Weight:   100,
		Rotation: RotationBestFit,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateOrientations(item, nil)
	}
}

func BenchmarkCalculateLookahead(b *testing.B) {
	prevItem := &OrientatedItem{Width: 10, Length: 10, Depth: 10}
	nextItems := make([]*Item, 8)
	for i := 0; i < 8; i++ {
		nextItems[i] = &Item{
			Width:    5 + i,
			Length:   5 + i,
			Depth:    5,
			Weight:   10,
			Rotation: RotationBestFit,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ClearLookaheadCache()
		CalculateLookahead(prevItem, nextItems, 100, 100, 100, 10, 8)
	}
}

func BenchmarkCalculateLookaheadCached(b *testing.B) {
	prevItem := &OrientatedItem{Width: 10, Length: 10, Depth: 10}
	nextItems := make([]*Item, 8)
	for i := 0; i < 8; i++ {
		nextItems[i] = &Item{
			Width:    5,
			Length:   5,
			Depth:    5,
			Weight:   10,
			Rotation: RotationBestFit,
		}
	}

	ClearLookaheadCache()
	// Warm up cache
	CalculateLookahead(prevItem, nextItems, 100, 100, 100, 10, 8)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CalculateLookahead(prevItem, nextItems, 100, 100, 100, 10, 8)
	}
}

func BenchmarkSimplePacker(b *testing.B) {
	box := &Box{
		InnerWidth:  100,
		InnerLength: 100,
		InnerDepth:  100,
		MaxWeight:   1000,
	}

	items := make([]*Item, 20)
	for i := 0; i < 20; i++ {
		items[i] = &Item{
			Width:    10,
			Length:   10,
			Depth:    10,
			Weight:   10,
			Rotation: RotationBestFit,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		packer := &SimplePacker{Box: box, Items: items}
		packer.PackSimple()
	}
}
