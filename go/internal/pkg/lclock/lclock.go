package lclock

import "time"

// Used to standardise the way we get the current timestamp
// Returns the current timestamp in nanoseconds
func GetCurrentTimestamp() int64 {
	return time.Now().UnixNano()
}
