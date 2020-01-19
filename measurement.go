package metrics

import (
	"sync"
	"time"
)

type valueSnapper interface {
	// returns an int64, float64, string, bool
	snapValue() interface{}
}

type measurement struct {
	name   string
	tags   TagSet
	fields map[string]interface{}
	lock   *sync.Mutex
}

func newMeasurement(name string, tags TagSet, lock *sync.Mutex) *measurement {
	return &measurement{
		name: name,
		tags: tags,
		fields: make(map[string]interface{}),
		lock:lock,
	}
}

func (m *measurement) IntGauge(field string, valueProvider IntValueProvider) {
	m.lock.Lock()
	defer m.lock.Unlock()

	g := &intGauge{valueProvider: valueProvider}
	m.fields[field] = g
}

func (m *measurement) FloatGauge(field string, valueProvider FloatValueProvider) {
	m.lock.Lock()
	defer m.lock.Unlock()

	g := &floatGauge{valueProvider: valueProvider}
	m.fields[field] = g
}

func (m *measurement) Counter(field string) Counter {
	m.lock.Lock()
	defer m.lock.Unlock()

	c := &counter{
		lock:       m.lock,
		cumulative: false,
	}
	m.fields[field] = c
	return c
}

func (m *measurement) CumulativeCounter(field string) Counter {
	m.lock.Lock()
	defer m.lock.Unlock()

	c := &counter{
		lock:       m.lock,
		cumulative: true,
	}
	m.fields[field] = c
	return c
}

// snap gathers a snapshot of this measurement and the fields it contains.
//
// Note
// It is expected that the registry's data mutex has been already locked by the caller.
func (m *measurement) snap(timestamp time.Time, outerTags TagSet) *SnappedMeasurement {
	snapped := &SnappedMeasurement{
		Timestamp: timestamp,
		Name:      m.name,
		Tags:      outerTags.Merge(m.tags),
		Fields:    make(map[string]interface{}),
	}

	for fieldName, field := range m.fields {
		value := field.(valueSnapper).snapValue()
		snapped.Fields[fieldName] = value
	}

	return snapped
}

type intGauge struct {
	valueProvider IntValueProvider
}

func (i *intGauge) snapValue() interface{} {
	return i.valueProvider()
}

type floatGauge struct {
	valueProvider FloatValueProvider
}

func (f *floatGauge) snapValue() interface{} {
	return f.valueProvider()
}

type counter struct {
	lock       *sync.Mutex
	cumulative bool
	value      int64
}

func (c *counter) snapValue() interface{} {
	out := c.value
	if !c.cumulative {
		c.value = 0
	}
	return out
}

func (c *counter) Inc() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.value++
}

func (c *counter) Add(amount int64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.value += amount
}
