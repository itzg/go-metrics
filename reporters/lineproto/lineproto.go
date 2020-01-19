package lineproto

import (
	"context"
	"github.com/itzg/go-metrics"
	lpsender "github.com/itzg/line-protocol-sender"
	"time"
)

// ReportToLineProtocolSocket will periodically walk the registry and report the snapped metrics
// via the given line protocol sender. This function should be called in a go routine since it
// will block until the given context is done.
func ReportToLineProtocolSocket(ctx context.Context, registry metrics.Registry,
	interval time.Duration, client lpsender.Client) {

	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ctx.Done():
			return

		case timestamp := <-ticker.C:
			registry.Walk(timestamp, func(m *metrics.SnappedMeasurement) {
				client.Send(convertMetric(m))
			})
			client.Flush()
		}
	}
}

func convertMetric(m *metrics.SnappedMeasurement) *lpsender.SimpleMetric {
	converted := lpsender.NewSimpleMetric(m.Name)
	for _, tag := range m.Tags {
		converted.AddTag(tag.Key, tag.Value)
	}
	for k, v := range m.Fields {
		converted.AddField(k, v)
	}

	return converted
}
