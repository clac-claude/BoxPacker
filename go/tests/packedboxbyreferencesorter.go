package tests

// PackedBoxByReferenceSorter sorts packed boxes by reference
type PackedBoxByReferenceSorter struct {
	Reference string
}

// NewPackedBoxByReferenceSorter creates a new sorter
func NewPackedBoxByReferenceSorter(reference string) *PackedBoxByReferenceSorter {
	return &PackedBoxByReferenceSorter{
		Reference: reference,
	}
}

// Compare compares two packed boxes, prioritizing the one matching the reference
func (s *PackedBoxByReferenceSorter) Compare(boxA, boxB interface{}) int {
	// This would need access to the actual PackedBox structure
	// For now, returning a simplified implementation
	return 0 // Simplified for now
}
