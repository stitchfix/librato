/*
Package librato is a pure go client for publishing metrics to Librato.

The package publishes metrics asynchronously, at a regular interval. If the package is unable to
log metrics (network issues, service outage, or the app being overloaded), it will drop metrics instead of degrading the
application's performance. The package allows some control over this behavior by allowing the developer the option of
configuring the queue size. Once this queue size is exceeded, messages will be dropped until the publisher catches up.

The package also provides an Aggregator struct which can be used to aggregate the gauge measurements on the client. For
applications that need to log a substantial number of metrics, this will be preferable to publishing the individual metrics.


*/
package librato

import (
	"fmt"
	"net/url"
	"time"
)

const apiEndpoint = "https://metrics-api.librato.com"

type Librato struct {
	publisher *publisher
}

type Config struct {
	// Email used for logging into your librato account
	Email string
	// The Key used to access the librato api.
	APIKey string
	// An optional Queue size. By default, this will be 600
	QueueSize int
}

// New creates a new librato client. The client will harvest metrics and publish
// them every second. You can specify the QueueSize to control how many metrics
// the client will batch. If you exceed the queue size, the measures will be silently
// dropped.
func New(config Config, errCh chan<- error) *Librato {
	u, _ := url.Parse(apiEndpoint)
	u.User = url.UserPassword(config.Email, config.APIKey)
	u.Path = "/v1/metrics"

	// determine queue size
	queueSize := 600
	if config.QueueSize > 0 {
		queueSize = config.QueueSize
	}

	// start the publisher
	p := &publisher{
		metricsURL: u,
		queueSize:  queueSize,
		measures:   make(chan interface{}, queueSize),
		shutdown:   make(chan chan struct{}),
		errors:     errCh,
	}
	go p.run(time.Second * 1)

	return &Librato{publisher: p}
}

// Adds a Gauge measurement to librato. If the queue is full, the measure will be dropped,
// but an error will be published to the error channel if it was configured.
func (l *Librato) AddGauge(g Gauge) {
	select {
	case l.publisher.measures <- g:
	default:
		l.publisher.reportError(fmt.Errorf("gauge could not be added to the metrics queue"))
	}
}

// Adds a Counter measurement to librato. If the queue is full, the measure will be dropped,
// but an error will be published to the error channel if it was configured.
func (l *Librato) AddCounter(c Counter) {
	select {
	case l.publisher.measures <- c:
	default:
		l.publisher.reportError(fmt.Errorf("counter could not be added to the metrics queue"))
	}
}

// Shutdown stops the librato client. The operation is blocking, and will make one final attempt
// to harvest measures and send them to librato.
func (l *Librato) Shutdown() {
	close(l.publisher.measures)

	s := make(chan struct{})
	l.publisher.shutdown <- s
	<-s
}
