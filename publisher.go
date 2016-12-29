package librato

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type publisher struct {
	metricsURL *url.URL
	errors     chan<- error
	measures   chan interface{}
	shutdown   chan chan struct{}
	queueSize  int
}

// run is the main run loop, it harvests metrics everytime the harvest ticker is fired,
// and publishes them to librato.
// If the shutdown channel is signaled, one final harvest will be performed, then the
// loop will exit
func (p *publisher) run(harvestOn time.Duration) {
	ticker := time.NewTicker(harvestOn)
	for {
		select {
		case <-ticker.C:
			fmt.Println("ticker fired")
			p.doHarvest()
		case s := <-p.shutdown:
			fmt.Println("shutdown fired")
			ticker.Stop()
			p.doHarvest()
			close(s)
			return
		}
	}
}

func (p *publisher) doHarvest() {
	measures := p.readMeasures()
	if measures == nil {
		return
	}

	client := http.Client{}
	r, w := io.Pipe()

	go func() {
		defer w.Close()
		json.NewEncoder(w).Encode(measures)
	}()

	buf := &bytes.Buffer{}
	tee := io.TeeReader(r, buf)
	req, _ := http.NewRequest("POST", p.metricsURL.String(), tee)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		p.reportError(err)
		return
	}

	io.Copy(os.Stdout, buf)

	defer resp.Body.Close()
	fmt.Println(resp.Status)
	if resp.StatusCode != 200 {
		p.reportError(fmt.Errorf("an error occurred publishing metrics to librato : %v", resp.Status))
	}
}

func (p *publisher) reportError(err error) {
	select {
	case p.errors <- err:
	default:
	}
}

type measurementRequest struct {
	Gauges   []Gauge   `json:"gauges,omitempty"`
	Counters []Counter `json:"counters,omitempty"`
}

// readMeasures reads from the measures channel up to queueSize messages
// The messages are converted and added to the measurement request object
// which will be sent to librato
func (p *publisher) readMeasures() *measurementRequest {
	mr := &measurementRequest{}

	// process at most the number of queued messages
	for i := 0; i <= p.queueSize; i++ {
		select {
		case measure := <-p.measures:
			switch m := measure.(type) {
			case Counter:
				mr.Counters = append(mr.Counters, m)
			case Gauge:
				mr.Gauges = append(mr.Gauges, m)
			}
		default:
			break
		}
	}
	if mr.Gauges == nil && mr.Counters == nil {
		return nil
	}
	return mr
}
