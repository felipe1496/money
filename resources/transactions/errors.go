package transactions

import "fmt"

var (
	ErrFailedToFetchEntries = fmt.Errorf("failed to fetch entries")
	ErrToCountEntries       = fmt.Errorf("failed to count entries")
)
