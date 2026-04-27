package notify

import "errors"

var errRedactNoPatterns = errors.New("redact: at least one pattern is required")
