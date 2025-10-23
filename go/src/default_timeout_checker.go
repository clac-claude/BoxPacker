package boxpacker

import (
	"github.com/dvdoug/boxpacker/go/src/exception"
	"time"
)

// Box packing (3D bin packing, knapsack problem).
//
// @author Doug Wright

// DefaultTimeoutChecker - Default implementation of TimeoutChecker
type DefaultTimeoutChecker struct {
	timeout   float64
	startTime float64
}

// NewDefaultTimeoutChecker creates a new DefaultTimeoutChecker
func NewDefaultTimeoutChecker(timeout float64) *DefaultTimeoutChecker {
	return &DefaultTimeoutChecker{
		timeout: timeout,
	}
}

// Start begins the timeout timer
func (dtc *DefaultTimeoutChecker) Start(startTime *float64) {
	if startTime != nil {
		dtc.startTime = *startTime
	} else {
		dtc.startTime = float64(time.Now().UnixNano()) / 1e9
	}
}

// ThrowOnTimeout checks if timeout has been exceeded and returns error if so
func (dtc *DefaultTimeoutChecker) ThrowOnTimeout(currentTime *float64, message string) error {
	var current float64
	if currentTime != nil {
		current = *currentTime
	} else {
		current = float64(time.Now().UnixNano()) / 1e9
	}

	spentTime := current - dtc.startTime
	isTimeout := spentTime >= dtc.timeout

	if isTimeout {
		if message == "" {
			message = "Exceeded the timeout"
		}
		return exception.NewTimeoutException(message, spentTime, dtc.timeout)
	}

	return nil
}
