package main

// EmbeddedMetrics is a struct intended to be emdeded in every output format
type EmbeddedMetrics struct {
	// Duration is the smartcontract execution time in millisecond
	Duration int `json:"duration"`
}

// AddDuration set the Duration metrics (in millisecond)
func (m *EmbeddedMetrics) AddDuration(millisecond int) {
	if m == nil {
		m = &EmbeddedMetrics{}
	}
	m.Duration = millisecond
}

// Meter is the interface output struc implement to include metrics
type Meter interface {
	AddDuration(int)
}
