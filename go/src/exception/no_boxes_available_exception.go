package exception

import "fmt"

// Box packing (3D bin packing, knapsack problem).
//
// @author Doug Wright

// NoBoxesAvailableException - Exception used when an item cannot be packed into any box.
type NoBoxesAvailableException struct {
	message  string
	itemList interface{} // ItemList - using interface{} to avoid circular dependency
}

// NewNoBoxesAvailableException creates a new NoBoxesAvailableException
func NewNoBoxesAvailableException(message string, itemList interface{}) *NoBoxesAvailableException {
	return &NoBoxesAvailableException{
		message:  message,
		itemList: itemList,
	}
}

// Error implements the error interface
func (e *NoBoxesAvailableException) Error() string {
	return e.message
}

// GetAffectedItems returns the item list that couldn't be packed
func (e *NoBoxesAvailableException) GetAffectedItems() interface{} {
	return e.itemList
}

// String returns string representation
func (e *NoBoxesAvailableException) String() string {
	return fmt.Sprintf("NoBoxesAvailableException: %s", e.message)
}
