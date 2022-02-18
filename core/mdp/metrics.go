package mdp

// WorkerMetrics holds metric values about a single worker instance.
type WorkerMetrics struct {
	// connectedAt time.Time
}

// ClientMetrics holds metric values about a single client instance.
type ClientMetrics struct {
}

// ServiceMetrics holds metrics for a service that is connected to a broker and one or more workers.
type ServiceMetrics struct {
	// firstSeen time.Time
}

// BrokerMetrics holds metric values about a broker instance.
type BrokerMetrics struct {
	// upTime time.Time
}
