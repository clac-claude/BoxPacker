package boxpacker

// Box packing (3D bin packing, knapsack problem).
//
// @author Doug Wright

// TimeoutChecker - Interface for timeout checking
type TimeoutChecker interface {
	// Start begins the timeout timer
	Start(startTime *float64)

	// ThrowOnTimeout checks if timeout has been exceeded and panics if so
	ThrowOnTimeout(currentTime *float64, message string) error
}
