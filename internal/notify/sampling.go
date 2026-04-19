package notify

import (
	"math/rand"
	"time"
)

// SamplingNotifier forwards only a random sample of messages to the inner
// notifier. This is useful for reducing noise in high-volume environments.
type SamplingNotifier struct {
	inner  Notifier
	rate   float64 // 0.0 – 1.0
	rander *rand.Rand
}

// NewSamplingNotifier returns a SamplingNotifier that forwards messages with
// probability rate (0.0 = never, 1.0 = always).
func NewSamplingNotifier(inner Notifier, rate float64) (*SamplingNotifier, error) {
	if inner == nil {
		return nil, errNilInner
	}
	if rate < 0 || rate > 1 {
		return nil, errInvalidRate
	}
	return &SamplingNotifier{
		inner:  inner,
		rate:   rate,
		rander: rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

// Send forwards the message to the inner notifier with probability s.rate.
func (s *SamplingNotifier) Send(msg Message) error {
	if s.rander.Float64() < s.rate {
		return s.inner.Send(msg)
	}
	return nil
}
