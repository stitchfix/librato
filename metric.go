package librato

// Aggregate provides a means for aggregating metrics on the client side, and pushing the
// the calculated value to librato.
type Aggregate struct {
	// Each metric has a name that is unique to its class of metrics e.g. a gauge name must be unique among gauges. The name
	// identifies a metric in subsequent API calls to store/query individual measurements and can be up to 255 characters in
	// length. Valid characters for metric names are A-Za-z0-9.:-_. The metric namespace is case insensitive.
	Name string `json:"name"`

	// Source is an optional property that can be used to subdivide a common gauge/counter among multiple members of a population.
	// For example the number of requests/second serviced by an application could be broken up among a group of server instances in
	// a scale-out tier by setting the hostname as the value of source.
	Source string `json:"source,omitempty"`

	// The epoch time at which an individual measurement occurred with a maximum resolution of seconds.
	MeasureTime int64 `json:"measure_time,omitempty"`

	values []float64
}

// Add a new value to the aggregate
// The method returns itself to allow chaining
func (a *Aggregate) Add(v float64) *Aggregate {
	a.values = append(a.values, v)
	return a
}

// Count is the number of measurements that have been recorded
func (a Aggregate) Count() int {
	return len(a.values)
}

// Sum is the sum of the measurements
func (a Aggregate) Sum() float64 {
	var total float64
	for _, v := range a.values {
		total += v
	}
	return total
}

// Min is the minimum measurement value, or zero of no measurements exist.
func (a Aggregate) Min() float64 {
	if len(a.values) == 0 {
		return 0
	}

	m := a.values[0]
	for _, v := range a.values {
		if v < m {
			m = v
		}
	}
	return m
}

// Max is the maximum measurement value, or zero of no measurements exist.
func (a Aggregate) Max() float64 {
	if len(a.values) == 0 {
		return 0
	}

	m := a.values[0]
	for _, v := range a.values {
		if v > m {
			m = v
		}
	}
	return m
}

// SumSquares is
func (a Aggregate) SumSquares() float64 {
	if len(a.values) == 0 {
		return 0
	}
	var sum, sumOfSquares float64
	for _, v := range a.values {
		sum += v
		sumOfSquares += v * v
	}

	x := (sum * sum) / float64(len(a.values))
	return sumOfSquares - x
}

func (a Aggregate) ToGauge() Gauge {
	return Gauge{
		Name:       a.Name,
		Source:     a.Source,
		Count:      a.Count(),
		Sum:        a.Sum(),
		Min:        a.Min(),
		Max:        a.Max(),
		SumSquares: a.SumSquares(),
	}
}

type Gauge struct {
	// Each metric has a name that is unique to its class of metrics e.g. a gauge name must be unique among gauges. The name
	// identifies a metric in subsequent API calls to store/query individual measurements and can be up to 255 characters in
	// length. Valid characters for metric names are A-Za-z0-9.:-_. The metric namespace is case insensitive.
	Name string `json:"name"`

	// The numeric value of an individual measurement. Multiple formats are supported (e.g. integer, floating point, etc) but the value must be numeric.
	Value interface{} `json:"value"`

	// Source is an optional property that can be used to subdivide a common gauge/counter among multiple members of a population.
	// For example the number of requests/second serviced by an application could be broken up among a group of server instances in
	// a scale-out tier by setting the hostname as the value of source.
	Source string `json:"source,omitempty"`

	// The epoch time at which an individual measurement occurred with a maximum resolution of seconds.
	MeasureTime int64 `json:"measure_time,omitempty"`

	Count      int         `json:"count,omitempty"`
	Sum        interface{} `json:"sum,omitempty"`
	Min        interface{} `json:"min,omitempty"`
	Max        interface{} `json:"max,omitempty"`
	SumSquares interface{} `json:"sum_squares,omitempty"`
}

type Counter struct {
	Name        string      `json:"name"`
	Value       interface{} `json:"value"`
	Source      string      `json:"source,omitempty"`
	MeasureTime int64       `json:"measure_time,omitempty"`
}
