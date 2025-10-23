package exception

import "fmt"

// Box packing (3D bin packing, knapsack problem).
//
// @author Doug Wright

// TimeoutException - Exception used when the timeout occurred.
type TimeoutException struct {
	message   string
	spentTime float64
	timeout   float64
}

// NewTimeoutException creates a new TimeoutException
func NewTimeoutException(message string, spentTime, timeout float64) *TimeoutException {
	return &TimeoutException{
		message:   message,
		spentTime: spentTime,
		timeout:   timeout,
	}
}

// Error implements the error interface
func (e *TimeoutException) Error() string {
	return e.message
}

// GetTimeout returns the timeout value
func (e *TimeoutException) GetTimeout() float64 {
	return e.timeout
}

// GetSpentTime returns the time spent before timeout
func (e *TimeoutException) GetSpentTime() float64 {
	return e.spentTime
}

// String returns string representation
func (e *TimeoutException) String() string {
	return fmt.Sprintf("TimeoutException: %s (spent: %.2fs, timeout: %.2fs)", e.message, e.spentTime, e.timeout)
}
