package notify

import "fmt"

// batchFormat is separated to avoid import cycle; re-exported via formatBatchSummary.
var _ = fmt.Sprintf // ensure fmt imported
