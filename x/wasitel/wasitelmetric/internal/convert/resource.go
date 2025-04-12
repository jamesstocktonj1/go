package convert

import (
	"go.opentelemetry.io/otel/sdk/resource"
	"go.wasmcloud.dev/x/wasitel/wasitelmetric/internal/types"
)

// Resource transforms a Resource into an OTLP Resource.
func Resource(r *resource.Resource) *types.Resource {
	if r == nil {
		return nil
	}
	return &types.Resource{Attributes: ResourceAttributes(r)}
}
