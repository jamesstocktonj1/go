package convert

import (
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.wasmcloud.dev/x/wasitel/wasitelmetric/internal/types"
)

func ResourceMetrics(data *metricdata.ResourceMetrics) (*types.ResourceMetrics, error) {
	return &types.ResourceMetrics{
		Resource:     Resource(data.Resource),
		ScopeMetrics: ScopeMetrics(data.ScopeMetrics),
		SchemaUrl:    data.Resource.SchemaURL(),
	}, nil
}
