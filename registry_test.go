package metrics_test

import (
	"github.com/itzg/go-metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRegistry_Walk(t *testing.T) {

	registry := metrics.NewRegistry(metrics.T("host", "h-1"))
	measurement := registry.Measurement("testing", metrics.Tag{Key: "instance", Value: "i-1"})
	counterEach := measurement.Counter("amount_each")
	counterTotal := measurement.CumulativeCounter("amount_total")
	measurement.IntGauge("size", func() int64 {
		return 3
	})
	measurement.FloatGauge("pct", func() float64 {
		return 0.5
	})

	counterEach.Inc()
	counterTotal.Inc()

	timestamp, err := time.ParseInLocation(time.RFC3339, "2006-01-02T15:04:05Z", time.UTC)
	require.NoError(t, err)

	var results []*metrics.SnappedMeasurement
	registry.Walk(timestamp, func(m *metrics.SnappedMeasurement) {
		results = append(results, m)
	})

	require.Len(t, results, 1)
	assert.ElementsMatch(t, results, []*metrics.SnappedMeasurement{
		{
			Timestamp: timestamp,
			Name:      "testing",
			Tags:      metrics.JoinTags(metrics.Tag{"host", "h-1"}, metrics.Tag{"instance", "i-1"}),
			Fields:    map[string]interface{}{"amount_each": int64(1), "amount_total": int64(1), "size": int64(3), "pct": float64(0.5)},
		},
	})

	// now test cumulative counter vs non

	counterTotal.Inc()

	results = results[:0]
	registry.Walk(timestamp, func(m *metrics.SnappedMeasurement) {
		results = append(results, m)
	})

	require.Len(t, results, 1)
	assert.ElementsMatch(t, results, []*metrics.SnappedMeasurement{
		{
			Timestamp: timestamp,
			Name:      "testing",
			Tags:      metrics.JoinTags(metrics.Tag{"host", "h-1"}, metrics.Tag{"instance", "i-1"}),
			Fields: map[string]interface{}{
				"amount_each":  int64(0),
				"amount_total": int64(2),
				"size":         int64(3),
				"pct":          float64(0.5)},
		},
	})
}
