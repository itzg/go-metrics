/*
Package metrics provides Go applications a simple way to report statistics.

Applications will start by creating a Registry and typically just one of those. One or more
Measurement's will be created/obtained from the Registry each Measurement best aligns with a
subsystem of the application since a measurement in turn contains one or more fields.

The fields of the measurement, counters and gauges, are what the application code will directly
use to increment and provide values representing the activity of the application.

Reporters, which may be a subsystem of the application, will periodically utilize the
Registry's Walk function to gather a snapshot of all registered metrics.
 */
package metrics

import "time"

type SnappedMeasurementConsumer func(m *SnappedMeasurement)

type Registry interface {
	Measurement(name string, tags ...Tag) Measurement
	Walk(timestamp time.Time, consumer SnappedMeasurementConsumer)
}

type IntValueProvider func() int64

type FloatValueProvider func() float64

type Measurement interface {
	IntGauge(field string, valueProvider IntValueProvider)
	FloatGauge(field string, valueProvider FloatValueProvider)
	Counter(field string) Counter
	CumulativeCounter(field string) Counter
}

type Counter interface {
	Inc()
	Add(amount int64)
}

type SnappedMeasurement struct {
	Timestamp time.Time
	Name      string
	Tags      TagSet
	// Fields maps named fields to either an int64, float64, string, bool
	Fields map[string]interface{}
}
