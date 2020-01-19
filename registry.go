package metrics

import (
	"sync"
	"time"
)

type registry struct {
	tags             TagSet
	measurementsLock sync.Mutex
	dataLock         sync.Mutex
	measurements     map[string]*measurement
}

func NewRegistry(tags... Tag) Registry {
	return &registry{
		tags: JoinTags(tags...),
		measurements: make(map[string]*measurement),
	}
}

func (r *registry) Measurement(name string, tags ...Tag) Measurement {
	tagSet := JoinTags(tags...)
	key := tagSet.HashKey()

	r.measurementsLock.Lock()
	defer r.measurementsLock.Unlock()

	if m, exists := r.measurements[key]; exists {
		return m
	} else {
		m = newMeasurement(name, tagSet, &r.dataLock)
		r.measurements[key] = m
		return m
	}
}

func (r *registry) Walk(timestamp time.Time, consumer SnappedMeasurementConsumer) {
	r.measurementsLock.Lock()
	defer r.measurementsLock.Unlock()

	r.dataLock.Lock()
	defer r.dataLock.Unlock()

	for _, measurement := range r.measurements {
		consumer(measurement.snap(timestamp, r.tags))
	}
}
