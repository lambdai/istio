// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package virtualservice

import (
	"strings"

	"istio.io/api/networking/v1alpha3"
	"istio.io/istio/galley/pkg/config/analysis"
	"istio.io/istio/galley/pkg/config/analysis/analyzers/util"
	"istio.io/istio/galley/pkg/config/analysis/msg"
	"istio.io/istio/pkg/config/resource"
	"istio.io/istio/pkg/config/schema/collection"
	"istio.io/istio/pkg/config/schema/collections"
)

// ConflictingTlsGatewayHostsAnalyzer checks if Mtls virtual services can be accessed by non-mtls method.
type ConflictingTlsGatewayHostsAnalyzer struct{}

var _ analysis.Analyzer = &ConflictingMeshGatewayHostsAnalyzer{}

// Metadata implements Analyzer
func (c *ConflictingTlsGatewayHostsAnalyzer) Metadata() analysis.Metadata {
	return analysis.Metadata{
		Name:        "virtualservice.ConflictingTlsGatewayHostsAnalyzer",
		Description: "Checks if Mtls virtual services can be accessed by non-mtls method",
		Inputs: collection.Names{
			collections.IstioNetworkingV1Alpha3Virtualservices.Name(),
		},
	}
}

// Analyze implements Analyzer
func (c *ConflictingTlsGatewayHostsAnalyzer) Analyze(ctx analysis.Context) {
	mtlsGateways, openGateways := initHostsUnderGateway(ctx)
	ctx.ForEach(collections.IstioNetworkingV1Alpha3Virtualservices.Name(), func(r *resource.Instance) bool {
		vs := r.Message.(*v1alpha3.VirtualService)
		anyFault := false || len(vs.Http) == 0

		for _, httpRoute := range vs.Http {
			if httpRoute.Fault != nil {
				anyFault = true
				break
			}
		}

		// Http Fault is defined. Huristically assume user is aware of the risk.
		if anyFault {
			return true
		}

		// No entry in Gateways implies "mesh" by default
		if len(vs.Gateways) == 0 {
			// Both in open gateway and in mtlsgateway.
			for gatewayName, gView := range mtlsGateways {
				for _, h := range vs.Hosts {
					fqdn := util.ConvertHostToFQDN(r.Metadata.FullName.Namespace, h)
					if overlaps := gView.overlap(fqdn); len(overlaps) > 0 {
						// TODO(lambdai)
						m := msg.NewVirtualServiceMtlsAccessNonMtls(r, fqdn, gatewayName)
						ctx.Report(collections.IstioNetworkingV1Alpha3Virtualservices.Name(), m)
					}
				}
			}
		} else {
			for _, gatewayName := range vs.Gateways {
				if _, inOpenGateway := openGateways[gatewayName]; !inOpenGateway {
					// The attached gateway is not open accessible. Not conflict with mtls gateway.
					continue
				}
				// Both in open gateway and in mtlsgateway.
				if gView, ok := mtlsGateways[gatewayName]; ok {
					for _, h := range vs.Hosts {
						fqdn := util.ConvertHostToFQDN(r.Metadata.FullName.Namespace, h)

						if overlaps := gView.overlap(fqdn); len(overlaps) > 0 {
							m := msg.NewVirtualServiceMtlsAccessNonMtls(r, fqdn, gatewayName)
							ctx.Report(collections.IstioNetworkingV1Alpha3Virtualservices.Name(), m)
						}
					}
				}
			}
		}
		return true
	})
}

type gatewayView struct {
	fqdn []string
}

func (g gatewayView) overlap(fqdn string) []string {
	// TODO(lambdai): Build match tree.
	overlaps := make([]string, 0)
	for _, h := range g.fqdn {
		if fqdnOverlap(h, fqdn) {
			overlaps = append(overlaps, fqdn)
		}
	}
	return overlaps
}

func fqdnOverlap(h1, h2 string) bool {
	if strings.HasPrefix(h1, string("*")) && strings.HasPrefix(h2, "*") {
		return strings.HasSuffix(h1[1:], h2[1:]) || strings.HasSuffix(h2[1:], h1[1:])
	} else if strings.HasPrefix(h1, "*") && !strings.HasPrefix(h2, "*") {
		return strings.HasSuffix(h2, h1[1:])
	} else if !strings.HasPrefix(h1, "*") && strings.HasPrefix(h2, "*") {
		return strings.HasSuffix(h1, h2[1:])
	} else {
		return h1 == h2
	}
}

func allMtlsGateways(ctx analysis.Context) (map[string]gatewayView, map[string]gatewayView) {
	// Mtls Gateway name to hosts. The value will be used in the future.
	mtlsGateways := make(map[string]gatewayView)
	// Open gateways that does not require client validation.
	openGateways := make(map[string]gatewayView)
	ctx.ForEach(collections.IstioNetworkingV1Alpha3Gateways.Name(), func(r *resource.Instance) bool {
		gw := r.Message.(*v1alpha3.Gateway)
		for _, server := range gw.Servers {
			if server.Tls != nil && server.Tls.Mode == v1alpha3.ServerTLSSettings_MUTUAL {
				if _, ok := mtlsGateways[r.Metadata.FullName.String()]; !ok {
					mtlsGateways[r.Metadata.FullName.String()] = gatewayView{}
				}
			} else if server.Tls == nil || server.Tls.Mode == v1alpha3.ServerTLSSettings_SIMPLE {
				if _, ok := openGateways[r.Metadata.FullName.String()]; !ok {
					openGateways[r.Metadata.FullName.String()] = gatewayView{}
				}
			}
		}
		return true
	})
	return mtlsGateways, openGateways
}

func initHostsUnderGateway(ctx analysis.Context) (map[string]gatewayView, map[string]gatewayView) {
	mtlsGateways, openGateways := allMtlsGateways(ctx)
	ctx.ForEach(collections.IstioNetworkingV1Alpha3Virtualservices.Name(), func(r *resource.Instance) bool {
		vs := r.Message.(*v1alpha3.VirtualService)
		attachedGateways := vs.Gateways
		// No entry in Gateways imply "mesh" by default
		if len(attachedGateways) == 0 {
			attachedGateways = []string{util.MeshGateway}
		}
		for _, gw := range attachedGateways {
			if gView, ok := mtlsGateways[gw]; ok {
				for _, h := range vs.Hosts {
					fqdn := util.ConvertHostToFQDN(r.Metadata.FullName.Namespace, h)
					gView.fqdn = append(gView.fqdn, fqdn)
				}
			}
		}
		return true
	})
	return mtlsGateways, openGateways
}
