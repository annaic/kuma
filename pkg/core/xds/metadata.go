package xds

import (
	"strconv"
	"strings"

	"google.golang.org/protobuf/types/known/structpb"

	mesh_proto "github.com/kumahq/kuma/api/mesh/v1alpha1"
	"github.com/kumahq/kuma/pkg/core"
	core_mesh "github.com/kumahq/kuma/pkg/core/resources/apis/mesh"
	"github.com/kumahq/kuma/pkg/core/resources/model"
	"github.com/kumahq/kuma/pkg/core/resources/model/rest"
	util_proto "github.com/kumahq/kuma/pkg/util/proto"
)

var metadataLog = core.Log.WithName("xds-server").WithName("metadata-tracker")

const (
	// Supported Envoy node metadata fields.
	fieldDataplaneAdminPort           = "dataplane.admin.port"
	fieldDataplaneApplicationsMetrics = "dataplane.applications.metrics"
	fieldDataplaneDNSPort             = "dataplane.dns.port"
	fieldDataplaneDNSEmptyPort        = "dataplane.dns.empty.port"
	fieldDataplaneDataplaneResource   = "dataplane.resource"
	fieldDynamicMetadata              = "dynamicMetadata"
	fieldDataplaneProxyType           = "dataplane.proxyType"
	fieldVersion                      = "version"
	FieldPrefixDependenciesVersion    = "version.dependencies"
)

// DataplaneMetadata represents environment-specific part of a dataplane configuration.
//
// This information might change from one dataplane run to another,
// and therefore it cannot be a part of Dataplane resource.
//
// On start-up, a dataplane captures its effective configuration (that might come
// from a file, environment variables and command line options) and includes it
// into request for a bootstrap config.
// Control Plane can use this information to fill in node metadata in the bootstrap
// config.
// Envoy will include node metadata from the bootstrap config
// at least into the very first discovery request on every xDS stream.
// This way, xDS server will be able to use Envoy node metadata
// to generate xDS resources that depend on environment-specific configuration.
type DataplaneMetadata struct {
	Resource            model.Resource
	AdminPort           uint32
	DNSPort             uint32
	EmptyDNSPort        uint32
	DynamicMetadata     map[string]string
	ProxyType           mesh_proto.ProxyType
	Version             *mesh_proto.Version
	ApplicationsMetrics []*mesh_proto.PrometheusServicesMetricsAggregateConfig
}

// GetDataplaneResource returns the underlying DataplaneResource, if present.
// If the resource is of a different type, it returns nil.
func (m *DataplaneMetadata) GetDataplaneResource() *core_mesh.DataplaneResource {
	if m != nil {
		if d, ok := m.Resource.(*core_mesh.DataplaneResource); ok {
			return d
		}
	}

	return nil
}

// GetZoneIngressResource returns the underlying ZoneIngressResource, if present.
// If the resource is of a different type, it returns nil.
func (m *DataplaneMetadata) GetZoneIngressResource() *core_mesh.ZoneIngressResource {
	if m != nil {
		if z, ok := m.Resource.(*core_mesh.ZoneIngressResource); ok {
			return z
		}
	}

	return nil
}

// GetZoneEgressResource returns the underlying ZoneEgressResource, if present.
// If the resource is of a different type, it returns nil.
func (m *DataplaneMetadata) GetZoneEgressResource() *core_mesh.ZoneEgressResource {
	if m != nil {
		if z, ok := m.Resource.(*core_mesh.ZoneEgressResource); ok {
			return z
		}
	}

	return nil
}

func (m *DataplaneMetadata) GetProxyType() mesh_proto.ProxyType {
	if m == nil || m.ProxyType == "" {
		return mesh_proto.DataplaneProxyType
	}
	return m.ProxyType
}

func (m *DataplaneMetadata) GetAdminPort() uint32 {
	if m == nil {
		return 0
	}
	return m.AdminPort
}

func (m *DataplaneMetadata) GetDNSPort() uint32 {
	if m == nil {
		return 0
	}
	return m.DNSPort
}

func (m *DataplaneMetadata) GetEmptyDNSPort() uint32 {
	if m == nil {
		return 0
	}
	return m.EmptyDNSPort
}

func (m *DataplaneMetadata) GetDynamicMetadata(key string) string {
	if m == nil || m.DynamicMetadata == nil {
		return ""
	}
	return m.DynamicMetadata[key]
}

func (m *DataplaneMetadata) GetVersion() *mesh_proto.Version {
	if m == nil {
		return nil
	}
	return m.Version
}

func DataplaneMetadataFromXdsMetadata(xdsMetadata *structpb.Struct) *DataplaneMetadata {
	// Be extra careful here about nil checks since xdsMetadata is a "user" input.
	// Even if we know that something should not be nil since we are generating metadata,
	// the DiscoveryRequest can still be crafted manually to crash the CP.
	metadata := DataplaneMetadata{}
	if xdsMetadata == nil {
		return &metadata
	}
	if field := xdsMetadata.Fields[fieldDataplaneProxyType]; field != nil {
		metadata.ProxyType = mesh_proto.ProxyType(field.GetStringValue())
	}
	metadata.AdminPort = uint32Metadata(xdsMetadata, fieldDataplaneAdminPort)
	metadata.DNSPort = uint32Metadata(xdsMetadata, fieldDataplaneDNSPort)
	metadata.EmptyDNSPort = uint32Metadata(xdsMetadata, fieldDataplaneDNSEmptyPort)
	if value := xdsMetadata.Fields[fieldDataplaneDataplaneResource]; value != nil {
		res, err := rest.UnmarshallToCore([]byte(value.GetStringValue()))
		if err != nil {
			metadataLog.Error(err, "invalid value in dataplane metadata", "field", fieldDataplaneDataplaneResource, "value", value)
		}
		switch r := res.(type) {
		case *core_mesh.DataplaneResource,
			*core_mesh.ZoneIngressResource,
			*core_mesh.ZoneEgressResource:
			metadata.Resource = r
		default:
			metadataLog.Error(err, "invalid dataplane resource type",
				"resource", r.Descriptor().Name,
				"field", fieldDataplaneDataplaneResource,
				"value", value)
		}
	}

	if value := xdsMetadata.Fields[fieldVersion]; value.GetStructValue() != nil {
		version := &mesh_proto.Version{}
		if err := util_proto.ToTyped(value.GetStructValue(), version); err != nil {
			metadataLog.Error(err, "invalid value in dataplane metadata", "field", fieldVersion, "value", value)
		}
		metadata.Version = version
	}

	if value := xdsMetadata.Fields[fieldDynamicMetadata]; value != nil {
		dynamicMetadata := map[string]string{}
		for field, val := range value.GetStructValue().GetFields() {
			if strings.HasPrefix(field, FieldPrefixDependenciesVersion) {
				dependencyName := strings.TrimPrefix(field, FieldPrefixDependenciesVersion+".")
				if metadata.GetVersion().GetDependencies() != nil {
					metadata.Version.Dependencies[dependencyName] = val.GetStringValue()
				}
			} else {
				dynamicMetadata[field] = val.GetStringValue()
			}
		}
		metadata.DynamicMetadata = dynamicMetadata
	}

	if value := xdsMetadata.Fields[fieldDataplaneApplicationsMetrics]; value != nil {
		applicationMetricsConfig := []*mesh_proto.PrometheusServicesMetricsAggregateConfig{}
		for _, val := range value.GetListValue().GetValues() {
			appConfig := mesh_proto.PrometheusServicesMetricsAggregateConfig{}
			if scrapeConfig := val.GetStructValue(); scrapeConfig != nil {
				if path := scrapeConfig.Fields["path"].GetStringValue(); path != "" {
					appConfig.Path = path
				}
				if name := scrapeConfig.Fields["name"].GetStringValue(); name != "" {
					appConfig.Name = name
				}
				appConfig.Port = uint32Metadata(scrapeConfig, "port")
			}
		}
		metadata.ApplicationsMetrics = applicationMetricsConfig
	}

	return &metadata
}

func uint32Metadata(xdsMetadata *structpb.Struct, field string) uint32 {
	value := xdsMetadata.Fields[field]
	if value == nil {
		return 0
	}
	port, err := strconv.Atoi(value.GetStringValue())
	if err != nil {
		metadataLog.Error(err, "invalid value in dataplane metadata", "field", field, "value", value)
		return 0
	}
	return uint32(port)
}
