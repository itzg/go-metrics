package lineproto_test

import (
	"context"
	protocol "github.com/influxdata/line-protocol"
	"github.com/itzg/go-metrics"
	"github.com/itzg/go-metrics/reporters/lineproto"
	lpsender "github.com/itzg/line-protocol-sender"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type mockLpSenderClient struct {
	mock.Mock
}

func (c *mockLpSenderClient) Send(m protocol.Metric) {
	c.Called(m)
}

func (c *mockLpSenderClient) Flush() {
	c.Called()
}

type mockRegistry struct {
	mock.Mock
}

func (m *mockRegistry) Measurement(name string, tags ...metrics.Tag) metrics.Measurement {
	args := m.Called(name, tags)
	return args.Get(0).(metrics.Measurement)
}

func (m *mockRegistry) Walk(timestamp time.Time, consumer metrics.SnappedMeasurementConsumer) {
	m.Called(timestamp, consumer)
}

func TestReportToLineProtocolSocket(t *testing.T) {
	client := &mockLpSenderClient{}
	client.On("Send", mock.Anything)
	client.On("Flush")

	registry := &mockRegistry{}

	ctx, cancelFunc := context.WithCancel(context.Background())

	m := &metrics.SnappedMeasurement{
		Name: "testing",
	}

	registry.On("Walk", mock.Anything, mock.Anything).Run(
		func(args mock.Arguments) {
			args.Get(1).(metrics.SnappedMeasurementConsumer)(m)
			cancelFunc()
		})

	go lineproto.ReportToLineProtocolSocket(ctx, registry, 10*time.Millisecond, client)

	select {
	case <-ctx.Done():
		break
	case <-time.After(50 * time.Millisecond):
		t.Fatal("Walk was not called")
	}

	registry.AssertExpectations(t)

	sentMetric := lpsender.NewSimpleMetric("testing")
	client.AssertCalled(t, "Send", sentMetric)
	client.AssertCalled(t, "Flush")
}
