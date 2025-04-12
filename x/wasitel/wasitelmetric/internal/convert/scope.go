package convert

import (
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.wasmcloud.dev/x/wasitel/wasitelmetric/internal/types"
)

func ScopeMetrics(metrics []metricdata.ScopeMetrics) []*types.ScopeMetrics {
	sm := make([]*types.ScopeMetrics, len(metrics))
	for i, metric := range metrics {
		sm[i] = ScopeMetric(metric)
	}
	return sm
}

func ScopeMetric(metric metricdata.ScopeMetrics) *types.ScopeMetrics {
	return &types.ScopeMetrics{
		Scope:     Scope(metric.Scope),
		Metrics:   Metrics(metric.Metrics),
		SchemaUrl: metric.Scope.SchemaURL,
	}
}

func Scope(scope instrumentation.Scope) *types.InstrumentationScope {
	return &types.InstrumentationScope{
		Name:       scope.Name,
		Version:    scope.Version,
		Attributes: KeyValues(scope.Attributes.ToSlice()),
	}
}

func Metrics(metrics []metricdata.Metrics) []*types.Metric {
	met := make([]*types.Metric, len(metrics))
	for i, metric := range metrics {
		met[i] = Metric(metric)
	}
	return met
}

func Metric(metric metricdata.Metrics) *types.Metric {
	met := &types.Metric{
		Name:        metric.Name,
		Description: metric.Description,
		Unit:        metric.Unit,
	}
	switch a := metric.Data.(type) {
	case metricdata.Gauge[int64]:
		met.Gauge = Gauge(a)
	case metricdata.Gauge[float64]:
		met.Gauge = Gauge(a)
	case metricdata.Sum[int64]:
		met.Sum = Sum(a)
	case metricdata.Sum[float64]:
		met.Sum = Sum(a)
	case metricdata.Histogram[int64]:
		met.Histogram = Histogram(a)
	case metricdata.Histogram[float64]:
		met.Histogram = Histogram(a)
		// case metricdata.ExponentialHistogram[int64]:
		// 	met.ExponentialHistogram, err = ExponentialHistogram(a)
		// case metricdata.ExponentialHistogram[float64]:
		// 	met.ExponentialHistogram, err = ExponentialHistogram(a)
		// case metricdata.Summary:
		// 	met.Summary = Summary(a)
	}
	return met
}

func Gauge[N int64 | float64](g metricdata.Gauge[N]) *types.Gauge {
	return &types.Gauge{
		DataPoints: DataPoints(g.DataPoints),
	}
}

func Sum[N int64 | float64](s metricdata.Sum[N]) *types.Sum {
	return &types.Sum{
		DataPoints:             DataPoints(s.DataPoints),
		AggregationTemporality: types.AggregationTemporality(metricdata.DeltaTemporality),
		IsMonotonic:            s.IsMonotonic,
	}
}

func Histogram[N int64 | float64](h metricdata.Histogram[N]) *types.Histogram {
	return &types.Histogram{
		DataPoints:             HistDataPoints(h.DataPoints),
		AggregationTemporality: types.AggregationTemporality(metricdata.DeltaTemporality),
	}
}

func DataPoints[N int64 | float64](datapoints []metricdata.DataPoint[N]) []*types.NumberDataPoint {
	dp := make([]*types.NumberDataPoint, len(datapoints))
	for i, point := range datapoints {
		dp[i] = DataPoint(point)
	}
	return dp
}

func DataPoint[N int64 | float64](datapoint metricdata.DataPoint[N]) *types.NumberDataPoint {
	dp := &types.NumberDataPoint{
		Attributes:        Iterator(datapoint.Attributes.Iter()),
		StartTimeUnixNano: uint64(datapoint.StartTime.UnixNano()),
		TimeUnixNano:      uint64(datapoint.Time.UnixNano()),
		Exemplars:         Exemplars(datapoint.Exemplars),
	}
	switch i := any(datapoint.Value).(type) {
	case int64:
		dp.AsInt = (*types.NumberDataPoint_AsInt)(&i)
	case float64:
		dp.AsDouble = (*types.NumberDataPoint_AsDouble)(&i)
	}
	return dp
}

func Exemplars[N int64 | float64](exemplars []metricdata.Exemplar[N]) []*types.Exemplar {
	exmp := make([]*types.Exemplar, len(exemplars))
	for i, ex := range exemplars {
		exmp[i] = Exemplar(ex)
	}
	return exmp
}

func Exemplar[N int64 | float64](exemplar metricdata.Exemplar[N]) *types.Exemplar {
	ex := &types.Exemplar{
		FilteredAttributes: KeyValues(exemplar.FilteredAttributes),
		TimeUnixNano:       uint64(exemplar.Time.UnixNano()),
		SpanId:             (*types.SpanID)(&exemplar.SpanID),
		TraceId:            (*types.TraceID)(&exemplar.TraceID),
	}
	switch i := any(exemplar.Value).(type) {
	case int64:
		ex.AsInt = (*types.Exemplar_AsInt)(&i)
	case float64:
		ex.AsDouble = (*types.Exemplar_AsDouble)(&i)
	}
	return ex
}

func HistDataPoints[N int64 | float64](datapoints []metricdata.HistogramDataPoint[N]) []*types.HistogramDataPoint {
	dp := make([]*types.HistogramDataPoint, len(datapoints))
	for i, point := range datapoints {
		dp[i] = HistDataPoint(point)
	}
	return dp
}

func HistDataPoint[N int64 | float64](datapoint metricdata.HistogramDataPoint[N]) *types.HistogramDataPoint {
	dp := &types.HistogramDataPoint{
		Attributes:        Iterator(datapoint.Attributes.Iter()),
		StartTimeUnixNano: uint64(datapoint.StartTime.UnixNano()),
		TimeUnixNano:      uint64(datapoint.Time.UnixNano()),
		Count:             datapoint.Count,
		BucketCounts:      datapoint.BucketCounts,
		ExplicitBounds:    datapoint.Bounds,
		Exemplars:         Exemplars(datapoint.Exemplars),
	}
	switch i := any(datapoint.Sum).(type) {
	case float64:
		dp.Sum = &i
	}
	min, _ := datapoint.Min.Value()
	switch i := any(min).(type) {
	case float64:
		dp.Min = &i
	}
	max, _ := datapoint.Max.Value()
	switch i := any(max).(type) {
	case float64:
		dp.Max = &i
	}
	return dp
}
