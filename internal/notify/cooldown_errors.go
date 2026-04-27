package notify

import "errors"

// ErrZeroCooldown is returned when a zero or negative cooldown duration is
// provided to NewCooldownNotifier.
var ErrZeroCooldown = errors.New("notify: cooldown duration must be greater than zero")
