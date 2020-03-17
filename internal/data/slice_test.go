// Copyright 2020 OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package data

import (
	"testing"

	otlpmetrics "github.com/open-telemetry/opentelemetry-proto/gen/go/metrics/v1"
)

// InstrumentationLibraryMetrics is a collection of metrics from a Resource.
//
// Must use NewResourceMetrics functions to create new instances.
// Important: zero-initialized instance is not valid for use.
type InstrumentationLibraryMetricsV1 struct {
	// Wrap OTLP InstrumentationLibraryMetric.
	orig *otlpmetrics.InstrumentationLibraryMetrics

	// Override a few fields. These fields are the source of truth. Their counterparts
	// stored in corresponding fields of "orig" are ignored.
	pimpl *internalInstrumentationLibraryMetricsV1
}

type internalInstrumentationLibraryMetricsV1 struct {
	metrics []MetricV1
	// True when the slice was replace.
	sliceChanged bool
	// True if the pimpl was initialized.
	initialized bool
}

func newInstrumentationLibraryMetricsV1(orig *otlpmetrics.InstrumentationLibraryMetrics) InstrumentationLibraryMetricsV1 {
	return InstrumentationLibraryMetricsV1{orig, &internalInstrumentationLibraryMetricsV1{}}
}

func (ilm InstrumentationLibraryMetricsV1) Metrics() []MetricV1 {
	ilm.initInternallIfNeeded()
	return ilm.pimpl.metrics
}

func (ilm InstrumentationLibraryMetricsV1) SetMetrics(ms []MetricV1) {
	ilm.initInternallIfNeeded()
	ilm.pimpl.metrics = ms
	// We don't update the orig slice because this may be called multiple times.
	ilm.pimpl.sliceChanged = true
}

func (ilm InstrumentationLibraryMetricsV1) initInternallIfNeeded() {
	if !ilm.pimpl.initialized {
		ilm.pimpl.metrics = newMetricV1Slice(ilm.orig.Metrics)
		ilm.pimpl.initialized = true
	}
}

func (ilm InstrumentationLibraryMetricsV1) getOrig() *otlpmetrics.InstrumentationLibraryMetrics {
	return ilm.orig
}

func (ilm InstrumentationLibraryMetricsV1) flushInternal() {
	if !ilm.pimpl.initialized {
		// Guaranteed no changes via internal fields.
		return
	}

	if ilm.pimpl.sliceChanged {
		// Reconstruct the slice because we don't know what elements were removed/added.
		// User may have changed internal fields in any Metric, flush all of them.
		ilm.orig.Metrics = make([]*otlpmetrics.Metric, len(ilm.pimpl.metrics))
		for i := range ilm.pimpl.metrics {
			ilm.orig.Metrics[i] = ilm.pimpl.metrics[i].orig
			ilm.pimpl.metrics[i].flushInternal()
		}
	} else {
		// User may have changed internal fields in any Metric, flush all of them.
		for i := range ilm.pimpl.metrics {
			ilm.pimpl.metrics[i].flushInternal()
		}
	}
}

// No pimpl to not complicate the test. Keep it simple.
type MetricV1 struct {
	// Wrap OTLP Metric.
	orig *otlpmetrics.Metric
}

func newMetricV1Slice(origs []*otlpmetrics.Metric) []MetricV1 {
	// Slice for wrappers.
	wrappers := make([]MetricV1, len(origs))
	for i := range origs {
		wrappers[i].orig = origs[i]
	}
	return wrappers
}

func (m MetricV1) MetricDescriptor() MetricDescriptorV2 {
	if m.orig.MetricDescriptor == nil {
		m.orig.MetricDescriptor = &otlpmetrics.MetricDescriptor{}
	}
	return MetricDescriptorV2{m.orig.MetricDescriptor}
}

func (m MetricV1) SetMetricDescriptor(r MetricDescriptorV2) {
	m.orig.MetricDescriptor = r.orig
}

func (m MetricV1) flushInternal() {
	// Do something here to avoid compiler optimizations to not call this.
	// This will not happen in the benchmark
	if m.orig == nil {
		println("ERROR")
	}
}

type MetricDescriptorV1 struct {
	// Wrap OTLP MetricDescriptor.
	orig *otlpmetrics.MetricDescriptor
}

func (md MetricDescriptorV1) Name() string {
	return md.orig.Name
}

func (md MetricDescriptorV1) SetName(v string) {
	md.orig.Name = v
}

type InstrumentationLibraryMetricsV2 struct {
	// Wrap OTLP InstrumentationLibraryMetric.
	orig *otlpmetrics.InstrumentationLibraryMetrics
}

func NewInstrumentationLibraryMetricsV2(metricsCap int) InstrumentationLibraryMetricsV2 {
	orig := &otlpmetrics.InstrumentationLibraryMetrics{Metrics: make([]*otlpmetrics.Metric, 0, metricsCap)}
	return InstrumentationLibraryMetricsV2{orig}
}

func newInstrumentationLibraryMetricsV2(orig *otlpmetrics.InstrumentationLibraryMetrics) InstrumentationLibraryMetricsV2 {
	return InstrumentationLibraryMetricsV2{orig}
}

func (ilm InstrumentationLibraryMetricsV2) MetricsCount() int {
	return len(ilm.orig.Metrics)
}

func (ilm InstrumentationLibraryMetricsV2) ForEachMetric(fn func(MetricV2)) {
	for _, om := range ilm.orig.Metrics {
		fn(MetricV2{om})
	}
}

func (ilm InstrumentationLibraryMetricsV2) GetMetric(ix int) MetricV2 {
	return MetricV2{ilm.orig.Metrics[ix]}
}

func (ilm InstrumentationLibraryMetricsV2) ForEachMetricWithRemove(fn func(MetricV2) bool) {
	i := 0 // output index
	for _, om := range ilm.orig.Metrics {
		if fn(MetricV2{om}) {
			// copy and increment index
			ilm.orig.Metrics[i] = om
			i++
		}
	}
	ilm.orig.Metrics = ilm.orig.Metrics[:i]
}

func (ilm InstrumentationLibraryMetricsV2) AddMetric(ms MetricV2) {
	ilm.orig.Metrics = append(ilm.orig.Metrics, ms.orig)
}

type MetricV2 struct {
	// Wrap OTLP Metric.
	orig *otlpmetrics.Metric
}

func (m MetricV2) MetricDescriptor() MetricDescriptorV2 {
	if m.orig.MetricDescriptor == nil {
		m.orig.MetricDescriptor = &otlpmetrics.MetricDescriptor{}
	}
	return MetricDescriptorV2{m.orig.MetricDescriptor}
}

func (m MetricV2) SetMetricDescriptor(r MetricDescriptorV2) {
	m.orig.MetricDescriptor = r.orig
}

type MetricDescriptorV2 struct {
	// Wrap OTLP MetricDescriptor.
	orig *otlpmetrics.MetricDescriptor
}

func (md MetricDescriptorV2) Name() string {
	return md.orig.Name
}

func (md MetricDescriptorV2) SetName(v string) {
	md.orig.Name = v
}

func BenchmarkMetricV1(b *testing.B) {
	ils := &otlpmetrics.InstrumentationLibraryMetrics{
		InstrumentationLibrary: generateTestInstrumentationLibrary(),
		Metrics: []*otlpmetrics.Metric{
			generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(),
			generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(),
			generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(),
		},
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ilsV1 := newInstrumentationLibraryMetricsV1(ils)
		// In a processor
		for _, m := range ilsV1.Metrics() {
			if m.MetricDescriptor().Name() == "panic" {
				b.Fatal("This should not happen")
			}
		}

		// In another processor
		for _, m := range ilsV1.Metrics() {
			if m.MetricDescriptor().Name() == "panic" {
				b.Fatal("This should not happen")
			}
		}

		// Need to flush to get orig synchronized
		ilsV1.flushInternal()
	}
}

func BenchmarkMetricV1_NoFlush(b *testing.B) {
	ils := &otlpmetrics.InstrumentationLibraryMetrics{
		InstrumentationLibrary: generateTestInstrumentationLibrary(),
		Metrics: []*otlpmetrics.Metric{
			generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(),
			generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(),
			generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(),
		},
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ilsV1 := newInstrumentationLibraryMetricsV1(ils)
		// In a processor
		for _, m := range ilsV1.Metrics() {
			if m.MetricDescriptor().Name() == "panic" {
				b.Fatal("This should not happen")
			}
		}

		// In another processor
		for _, m := range ilsV1.Metrics() {
			if m.MetricDescriptor().Name() == "panic" {
				b.Fatal("This should not happen")
			}
		}
	}
}

func BenchmarkMetricV2(b *testing.B) {
	ils := &otlpmetrics.InstrumentationLibraryMetrics{
		InstrumentationLibrary: generateTestInstrumentationLibrary(),
		Metrics: []*otlpmetrics.Metric{
			generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(),
			generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(),
			generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(),
		},
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ilsV2 := newInstrumentationLibraryMetricsV2(ils)
		// In a processor
		ilsV2.ForEachMetric(func(m MetricV2) {
			if m.MetricDescriptor().Name() == "panic" {
				b.Fatal("This should not happen")
			}
		})

		// In another processor
		ilsV2.ForEachMetric(func(m MetricV2) {
			if m.MetricDescriptor().Name() == "panic" {
				b.Fatal("This should not happen")
			}
		})
	}
}

func BenchmarkMetricV2_GetMetric(b *testing.B) {
	ils := &otlpmetrics.InstrumentationLibraryMetrics{
		InstrumentationLibrary: generateTestInstrumentationLibrary(),
		Metrics: []*otlpmetrics.Metric{
			generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(),
			generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(),
			generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(),
		},
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ilsV2 := newInstrumentationLibraryMetricsV2(ils)
		// In a processor
		for i := 1; i < ilsV2.MetricsCount(); i++ {
			if ilsV2.GetMetric(i).MetricDescriptor().Name() == "panic" {
				b.Fatal("This should not happen")
			}
		}

		// In another processor
		for i := 1; i < ilsV2.MetricsCount(); i++ {
			if ilsV2.GetMetric(i).MetricDescriptor().Name() == "panic" {
				b.Fatal("This should not happen")
			}
		}
	}
}

func BenchmarkMetricV0_NoWrapper(b *testing.B) {
	ils := &otlpmetrics.InstrumentationLibraryMetrics{
		InstrumentationLibrary: generateTestInstrumentationLibrary(),
		Metrics: []*otlpmetrics.Metric{
			generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(),
			generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(),
			generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(), generateTestIntMetric(),
		},
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		// In a processor
		for _, m := range ils.Metrics {
			if m.MetricDescriptor.Name == "panic" {
				b.Fatal("This should not happen")
			}
		}

		// In another processor
		for _, m := range ils.Metrics {
			if m.MetricDescriptor.Name == "panic" {
				b.Fatal("This should not happen")
			}
		}
	}
}
