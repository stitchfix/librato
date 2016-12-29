librato
=======

This package provides an api for publishing gauges and counters to librato asynchonously.

Usage
----

From Go:

```go

  // Configure the authentication credentials
  // See struct docs for additional options
	config := Config{
		Email:  os.Getenv("LIBRATO_EMAIL"),
		APIKey: os.Getenv("LIBRATO_APIKEY"),
	}

  // Create a new librato instance
  // Each instance publishes independently of others, and
  // may be used to connect to multiple accounts
  l := librato.New(config)

  // Add a new gauge measurement.
	l.AddGauge(Gauge{Name: "reticulated.splines", Value: 1, Source: "add-gauge"})

  // add a new counter measurement
  l.AddCounter(Counter{Name: "reticulated.splines.counter", Value: 7, Source: "add-counter"})

  // Create and add an aggregate gauge
	a := Aggregate{Name: "reticulated.splines", Source: "add-aggregate"}
	a.Add(12).Add(6).Add(5).Add(4).Add(5).Add(10).Add(3)
	l.AddGauge(a.ToGauge())

  // When done, call Shutdown().
  // This operation is synchronous and will wait for the last
  // publish of metrics to complete
  l.Shutdown()

  ```


Installation
----

  To install as a library, you can use `go get`:

    go get github.com/stitchfix/librato
