package librato

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func newClient() *Librato {
	err := make(chan error, 10)
	config := Config{
		Email:  os.Getenv("LIBRATO_EMAIL"),
		APIKey: os.Getenv("LIBRATO_APIKEY"),
		Errors: err,
	}
	go func() {
		for e := range err {
			fmt.Println(e)
		}
	}()
	return New(config)
}

// The librato API
func TestAddGauge(t *testing.T) {
	l := newClient()
	pubTime := time.Now().Add(time.Minute * -1)
	l.AddGauge(Gauge{Name: "reticulated.splines", Value: 1, Source: "add-gauge", MeasureTime: pubTime.Unix()})
	l.AddGauge(Gauge{Name: "reticulated.splines", Value: 2, Source: "add-gauge", MeasureTime: pubTime.Add(time.Second * 5).Unix()})
	l.AddGauge(Gauge{Name: "reticulated.splines", Value: 4, Source: "add-gauge", MeasureTime: pubTime.Add(time.Second * 10).Unix()})
	l.AddGauge(Gauge{Name: "reticulated.splines", Value: 8, Source: "add-gauge", MeasureTime: pubTime.Add(time.Second * 15).Unix()})
	l.AddGauge(Gauge{Name: "reticulated.splines", Value: 16, Source: "add-gauge", MeasureTime: pubTime.Add(time.Second * 20).Unix()})
	l.AddGauge(Gauge{Name: "reticulated.splines", Value: 32, Source: "add-gauge", MeasureTime: pubTime.Add(time.Second * 25).Unix()})
	l.AddGauge(Gauge{Name: "reticulated.splines", Value: 64, Source: "add-gauge", MeasureTime: pubTime.Add(time.Second * 30).Unix()})
	l.Shutdown()
}

func TestAddAggregate(t *testing.T) {
	l := newClient()
	a := Aggregate{Name: "reticulated.splines", Source: "add-aggregate"}
	a.Add(12).Add(6).Add(5).Add(4).Add(5).Add(10).Add(3)
	l.AddGauge(a.ToGauge())
	l.Shutdown()
}

// The librato API
func TestAddCounter(t *testing.T) {
	l := newClient()
	pubTime := time.Now().Add(time.Minute * -1)
	for i := 0; i < 10; i++ {
		measureTime := pubTime.Add(time.Second * time.Duration(i*5))
		l.AddCounter(Counter{Name: "reticulated.splines.counter", Value: i * i, Source: "add-counter", MeasureTime: measureTime.Unix()})
	}
	l.Shutdown()
}

func TestAddGuageWithTimings(t *testing.T) {
	l := newClient()
	pubTime := time.Now().Add(time.Minute * -1)
	l.AddGauge(Gauge{Name: "reticulated.splines", Value: 10, Source: "add-gauge-with-timing", MeasureTime: pubTime.Add(time.Second * 30).Unix()})
	l.AddGauge(Gauge{Name: "reticulated.splines", Value: 20, Source: "add-gauge-with-timing", MeasureTime: pubTime.Add(time.Second * 15).Unix()})
	l.AddGauge(Gauge{Name: "reticulated.splines", Value: 40, Source: "add-gauge-with-timing", MeasureTime: pubTime.Add(time.Second * 10).Unix()})
	l.AddGauge(Gauge{Name: "reticulated.splines", Value: 70, Source: "add-gauge-with-timing", MeasureTime: pubTime.Unix()})
	l.Shutdown()
}
