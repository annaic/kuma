// Generated by tools/resource-gen.
// Run "make generate" to update this file.
// nolint:whitespace
package entities

import (
	"github.com/kumahq/kuma/pkg/core/resources/apis/mesh"
)

var meshEntities = []Definition{

	{
		Singular:     "circuit-breaker",
		Plural:       "circuit-breakers",
		ResourceType: mesh.CircuitBreakerType,
		ReadOnly:     false,
	},

	{
		Singular:     "dataplane",
		Plural:       "dataplanes",
		ResourceType: mesh.DataplaneType,
		ReadOnly:     false,
	},

	{
		Singular:     "external-service",
		Plural:       "external-services",
		ResourceType: mesh.ExternalServiceType,
		ReadOnly:     false,
	},

	{
		Singular:     "fault-injection",
		Plural:       "fault-injections",
		ResourceType: mesh.FaultInjectionType,
		ReadOnly:     false,
	},

	{
		Singular:     "health-check",
		Plural:       "health-checks",
		ResourceType: mesh.HealthCheckType,
		ReadOnly:     false,
	},

	{
		Singular:     "mesh",
		Plural:       "meshes",
		ResourceType: mesh.MeshType,
		ReadOnly:     false,
	},

	{
		Singular:     "proxytemplate",
		Plural:       "proxytemplates",
		ResourceType: mesh.ProxyTemplateType,
		ReadOnly:     false,
	},

	{
		Singular:     "rate-limit",
		Plural:       "rate-limits",
		ResourceType: mesh.RateLimitType,
		ReadOnly:     false,
	},

	{
		Singular:     "retry",
		Plural:       "retries",
		ResourceType: mesh.RetryType,
		ReadOnly:     false,
	},

	{
		Singular:     "timeout",
		Plural:       "timeouts",
		ResourceType: mesh.TimeoutType,
		ReadOnly:     false,
	},

	{
		Singular:     "traffic-log",
		Plural:       "traffic-logs",
		ResourceType: mesh.TrafficLogType,
		ReadOnly:     false,
	},

	{
		Singular:     "traffic-permission",
		Plural:       "traffic-permissions",
		ResourceType: mesh.TrafficPermissionType,
		ReadOnly:     false,
	},

	{
		Singular:     "traffic-route",
		Plural:       "traffic-routes",
		ResourceType: mesh.TrafficRouteType,
		ReadOnly:     false,
	},

	{
		Singular:     "traffic-trace",
		Plural:       "traffic-traces",
		ResourceType: mesh.TrafficTraceType,
		ReadOnly:     false,
	},

	{
		Singular:     "zone-ingress",
		Plural:       "zone-ingresses",
		ResourceType: mesh.ZoneIngressType,
		ReadOnly:     false,
	},
}