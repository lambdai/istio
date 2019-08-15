// Copyright 2017 Istio Authors
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

/*
Package component defines an in-memory representation of IstioControlPlane.<Feature>.<Component>. It provides functions
for manipulating the component and rendering a manifest from it.
See ../README.md for an architecture overview.
*/
package component

import (
	"fmt"

	"github.com/ghodss/yaml"

	"istio.io/operator/pkg/apis/istio/v1alpha2"
	"istio.io/operator/pkg/helm"
	"istio.io/operator/pkg/name"
	"istio.io/operator/pkg/patch"
	"istio.io/operator/pkg/translate"
	"istio.io/operator/pkg/util"
	"istio.io/pkg/log"
)

const (
	// String to emit for any component which is disabled.
	componentDisabledStr = " component is disabled."
	yamlCommentStr       = "# "
)

// Options defines options for a component.
type Options struct {
	// FeatureName is the name of the feature this component belongs to.
	FeatureName name.FeatureName
	// InstallSpec is the global IstioControlPlaneSpec.
	InstallSpec *v1alpha2.IstioControlPlaneSpec
	// Translator is the translator for this component.
	Translator *translate.Translator
}

// IstioComponent defines the interface for a component.
type IstioComponent interface {
	// name returns the name of the component.
	Name() name.ComponentName
	// Run starts the component. Must me called before the component is used.
	Run() error
	// RenderManifest returns a string with the rendered manifest for the component.
	RenderManifest() (string, error)
}

// CommonComponentFields is a struct common to all components.
type CommonComponentFields struct {
	*Options
	name     name.ComponentName
	started  bool
	renderer helm.TemplateRenderer
}

// CRDComponent is the pilot component.
type CRDComponent struct {
	*CommonComponentFields
}

// NewCRDComponent creates a new CRDComponent and returns a pointer to it.
func NewCRDComponent(opts *Options) *CRDComponent {
	return &CRDComponent{
		&CommonComponentFields{
			Options: opts,
			name:    name.IstioBaseComponentName,
		},
	}
}

// Run implements the IstioComponent interface.
func (c *CRDComponent) Run() error {
	return runComponent(c.CommonComponentFields)
}

// RenderManifest implements the IstioComponent interface.
func (c *CRDComponent) RenderManifest() (string, error) {
	if !c.started {
		return "", fmt.Errorf("component %s not started in RenderManifest", c.Name())
	}
	return renderManifest(c.CommonComponentFields)
}

// Name implements the IstioComponent interface.
func (c *CRDComponent) Name() name.ComponentName {
	return c.CommonComponentFields.name
}

// PilotComponent is the pilot component.
type PilotComponent struct {
	*CommonComponentFields
}

// NewPilotComponent creates a new PilotComponent and returns a pointer to it.
func NewPilotComponent(opts *Options) *PilotComponent {
	return &PilotComponent{
		&CommonComponentFields{
			Options: opts,
			name:    name.PilotComponentName,
		},
	}
}

// Run implements the IstioComponent interface.
func (c *PilotComponent) Run() error {
	return runComponent(c.CommonComponentFields)
}

// RenderManifest implements the IstioComponent interface.
func (c *PilotComponent) RenderManifest() (string, error) {
	if !c.started {
		return "", fmt.Errorf("component %s not started in RenderManifest", c.Name())
	}
	return renderManifest(c.CommonComponentFields)
}

// Name implements the IstioComponent interface.
func (c *PilotComponent) Name() name.ComponentName {
	return c.CommonComponentFields.name
}

// CitadelComponent is the pilot component.
type CitadelComponent struct {
	*CommonComponentFields
}

// NewCitadelComponent creates a new PilotComponent and returns a pointer to it.
func NewCitadelComponent(opts *Options) *CitadelComponent {
	return &CitadelComponent{
		&CommonComponentFields{
			Options: opts,
			name:    name.CitadelComponentName,
		},
	}
}

// Run implements the IstioComponent interface.
func (c *CitadelComponent) Run() error {
	return runComponent(c.CommonComponentFields)
}

// RenderManifest implements the IstioComponent interface.
func (c *CitadelComponent) RenderManifest() (string, error) {
	if !c.started {
		return "", fmt.Errorf("component %s not started in RenderManifest", c.Name())
	}
	return renderManifest(c.CommonComponentFields)
}

// Name implements the IstioComponent interface.
func (c *CitadelComponent) Name() name.ComponentName {
	return c.CommonComponentFields.name
}

// CertManagerComponent is the pilot component.
type CertManagerComponent struct {
	*CommonComponentFields
}

// NewCertManagerComponent creates a new PilotComponent and returns a pointer to it.
func NewCertManagerComponent(opts *Options) *CertManagerComponent {
	return &CertManagerComponent{
		&CommonComponentFields{
			Options: opts,
			name:    name.CertManagerComponentName,
		},
	}
}

// Run implements the IstioComponent interface.
func (c *CertManagerComponent) Run() error {
	return runComponent(c.CommonComponentFields)
}

// RenderManifest implements the IstioComponent interface.
func (c *CertManagerComponent) RenderManifest() (string, error) {
	if !c.started {
		return "", fmt.Errorf("component %s not started in RenderManifest", c.Name())
	}
	return renderManifest(c.CommonComponentFields)
}

// Name implements the IstioComponent interface.
func (c *CertManagerComponent) Name() name.ComponentName {
	return c.CommonComponentFields.name
}

// NodeAgentComponent is the pilot component.
type NodeAgentComponent struct {
	*CommonComponentFields
}

// NewNodeAgentComponent creates a new PilotComponent and returns a pointer to it.
func NewNodeAgentComponent(opts *Options) *NodeAgentComponent {
	return &NodeAgentComponent{
		&CommonComponentFields{
			Options: opts,
			name:    name.NodeAgentComponentName,
		},
	}
}

// Run implements the IstioComponent interface.
func (c *NodeAgentComponent) Run() error {
	return runComponent(c.CommonComponentFields)
}

// RenderManifest implements the IstioComponent interface.
func (c *NodeAgentComponent) RenderManifest() (string, error) {
	if !c.started {
		return "", fmt.Errorf("component %s not started in RenderManifest", c.Name())
	}
	return renderManifest(c.CommonComponentFields)
}

// Name implements the IstioComponent interface.
func (c *NodeAgentComponent) Name() name.ComponentName {
	return c.CommonComponentFields.name
}

// PolicyComponent is the pilot component.
type PolicyComponent struct {
	*CommonComponentFields
}

// NewPolicyComponent creates a new PilotComponent and returns a pointer to it.
func NewPolicyComponent(opts *Options) *PolicyComponent {
	return &PolicyComponent{
		&CommonComponentFields{
			Options: opts,
			name:    name.PolicyComponentName,
		},
	}
}

// Run implements the IstioComponent interface.
func (c *PolicyComponent) Run() error {
	return runComponent(c.CommonComponentFields)
}

// RenderManifest implements the IstioComponent interface.
func (c *PolicyComponent) RenderManifest() (string, error) {
	if !c.started {
		return "", fmt.Errorf("component %s not started in RenderManifest", c.Name())
	}
	return renderManifest(c.CommonComponentFields)
}

// Name implements the IstioComponent interface.
func (c *PolicyComponent) Name() name.ComponentName {
	return c.CommonComponentFields.name
}

// TelemetryComponent is the pilot component.
type TelemetryComponent struct {
	*CommonComponentFields
}

// NewTelemetryComponent creates a new PilotComponent and returns a pointer to it.
func NewTelemetryComponent(opts *Options) *TelemetryComponent {
	return &TelemetryComponent{
		&CommonComponentFields{
			Options: opts,
			name:    name.TelemetryComponentName,
		},
	}
}

// Run implements the IstioComponent interface.
func (c *TelemetryComponent) Run() error {
	return runComponent(c.CommonComponentFields)
}

// RenderManifest implements the IstioComponent interface.
func (c *TelemetryComponent) RenderManifest() (string, error) {
	if !c.started {
		return "", fmt.Errorf("component %s not started in RenderManifest", c.Name())
	}
	return renderManifest(c.CommonComponentFields)
}

// Name implements the IstioComponent interface.
func (c *TelemetryComponent) Name() name.ComponentName {
	return c.CommonComponentFields.name
}

// GalleyComponent is the pilot component.
type GalleyComponent struct {
	*CommonComponentFields
}

// NewGalleyComponent creates a new PilotComponent and returns a pointer to it.
func NewGalleyComponent(opts *Options) *GalleyComponent {
	return &GalleyComponent{
		&CommonComponentFields{
			Options: opts,
			name:    name.GalleyComponentName,
		},
	}
}

// Run implements the IstioComponent interface.
func (c *GalleyComponent) Run() error {
	return runComponent(c.CommonComponentFields)
}

// RenderManifest implements the IstioComponent interface.
func (c *GalleyComponent) RenderManifest() (string, error) {
	if !c.started {
		return "", fmt.Errorf("component %s not started in RenderManifest", c.Name())
	}
	return renderManifest(c.CommonComponentFields)
}

// Name implements the IstioComponent interface.
func (c *GalleyComponent) Name() name.ComponentName {
	return c.CommonComponentFields.name
}

// SidecarInjectorComponent is the pilot component.
type SidecarInjectorComponent struct {
	*CommonComponentFields
}

// NewSidecarInjectorComponent creates a new PilotComponent and returns a pointer to it.
func NewSidecarInjectorComponent(opts *Options) *SidecarInjectorComponent {
	return &SidecarInjectorComponent{
		&CommonComponentFields{
			Options: opts,
			name:    name.SidecarInjectorComponentName,
		},
	}
}

// Run implements the IstioComponent interface.
func (c *SidecarInjectorComponent) Run() error {
	return runComponent(c.CommonComponentFields)
}

// RenderManifest implements the IstioComponent interface.
func (c *SidecarInjectorComponent) RenderManifest() (string, error) {
	if !c.started {
		return "", fmt.Errorf("component %s not started in RenderManifest", c.Name())
	}
	return renderManifest(c.CommonComponentFields)
}

// Name implements the IstioComponent interface.
func (c *SidecarInjectorComponent) Name() name.ComponentName {
	return c.CommonComponentFields.name
}

// IngressComponent is the ingress gateway component.
type IngressComponent struct {
	*CommonComponentFields
}

// NewIngressComponent creates a new IngressComponent and returns a pointer to it.
func NewIngressComponent(opts *Options) *IngressComponent {
	return &IngressComponent{
		&CommonComponentFields{
			Options: opts,
			name:    name.IngressComponentName,
		},
	}
}

// Run implements the IstioComponent interface.
func (c *IngressComponent) Run() error {
	return runComponent(c.CommonComponentFields)
}

// RenderManifest implements the IstioComponent interface.
func (c *IngressComponent) RenderManifest() (string, error) {
	if !c.started {
		return "", fmt.Errorf("component %s not started in RenderManifest", c.Name())
	}
	return renderManifest(c.CommonComponentFields)
}

// Name implements the IstioComponent interface.
func (c *IngressComponent) Name() name.ComponentName {
	return c.CommonComponentFields.name
}

// EgressComponent is the egress gateway component.
type EgressComponent struct {
	*CommonComponentFields
}

// NewEgressComponent creates a new IngressComponent and returns a pointer to it.
func NewEgressComponent(opts *Options) *EgressComponent {
	return &EgressComponent{
		&CommonComponentFields{
			Options: opts,
			name:    name.EgressComponentName,
		},
	}
}

// Run implements the IstioComponent interface.
func (c *EgressComponent) Run() error {
	return runComponent(c.CommonComponentFields)
}

// RenderManifest implements the IstioComponent interface.
func (c *EgressComponent) RenderManifest() (string, error) {
	if !c.started {
		return "", fmt.Errorf("component %s not started in RenderManifest", c.Name())
	}
	return renderManifest(c.CommonComponentFields)
}

// Name implements the IstioComponent interface.
func (c *EgressComponent) Name() name.ComponentName {
	return c.CommonComponentFields.name
}

// PrometheusComponent is the egress gateway component.
type PrometheusComponent struct {
	*CommonComponentFields
}

// NewPrometheusComponent creates a new IngressComponent and returns a pointer to it.
func NewPrometheusComponent(opts *Options) *PrometheusComponent {
	return &PrometheusComponent{
		&CommonComponentFields{
			Options: opts,
			name:    name.PrometheusComponentName,
		},
	}
}

// Run implements the IstioComponent interface.
func (c *PrometheusComponent) Run() error {
	return runComponent(c.CommonComponentFields)
}

// RenderManifest implements the IstioComponent interface.
func (c *PrometheusComponent) RenderManifest() (string, error) {
	if !c.started {
		return "", fmt.Errorf("component %s not started in RenderManifest", c.Name())
	}
	return renderManifest(c.CommonComponentFields)
}

// Name implements the IstioComponent interface.
func (c *PrometheusComponent) Name() name.ComponentName {
	return c.CommonComponentFields.name
}

// PrometheusOperatorComponent is the egress gateway component.
type PrometheusOperatorComponent struct {
	*CommonComponentFields
}

// NewPrometheusOperatorComponent creates a new IngressComponent and returns a pointer to it.
func NewPrometheusOperatorComponent(opts *Options) *PrometheusOperatorComponent {
	return &PrometheusOperatorComponent{
		&CommonComponentFields{
			Options: opts,
			name:    name.PrometheusOperatorComponentName,
		},
	}
}

// Run implements the IstioComponent interface.
func (c *PrometheusOperatorComponent) Run() error {
	return runComponent(c.CommonComponentFields)
}

// RenderManifest implements the IstioComponent interface.
func (c *PrometheusOperatorComponent) RenderManifest() (string, error) {
	if !c.started {
		return "", fmt.Errorf("component %s not started in RenderManifest", c.Name())
	}
	return renderManifest(c.CommonComponentFields)
}

// Name implements the IstioComponent interface.
func (c *PrometheusOperatorComponent) Name() name.ComponentName {
	return c.CommonComponentFields.name
}

// GrafanaComponent is the egress gateway component.
type GrafanaComponent struct {
	*CommonComponentFields
}

// NewGrafanaComponent creates a new IngressComponent and returns a pointer to it.
func NewGrafanaComponent(opts *Options) *GrafanaComponent {
	return &GrafanaComponent{
		&CommonComponentFields{
			Options: opts,
			name:    name.GrafanaComponentName,
		},
	}
}

// Run implements the IstioComponent interface.
func (c *GrafanaComponent) Run() error {
	return runComponent(c.CommonComponentFields)
}

// RenderManifest implements the IstioComponent interface.
func (c *GrafanaComponent) RenderManifest() (string, error) {
	if !c.started {
		return "", fmt.Errorf("component %s not started in RenderManifest", c.Name())
	}
	return renderManifest(c.CommonComponentFields)
}

// Name implements the IstioComponent interface.
func (c *GrafanaComponent) Name() name.ComponentName {
	return c.CommonComponentFields.name
}

// KialiComponent is the egress gateway component.
type KialiComponent struct {
	*CommonComponentFields
}

// NewKialiComponent creates a new IngressComponent and returns a pointer to it.
func NewKialiComponent(opts *Options) *KialiComponent {
	return &KialiComponent{
		&CommonComponentFields{
			Options: opts,
			name:    name.KialiComponentName,
		},
	}
}

// Run implements the IstioComponent interface.
func (c *KialiComponent) Run() error {
	return runComponent(c.CommonComponentFields)
}

// RenderManifest implements the IstioComponent interface.
func (c *KialiComponent) RenderManifest() (string, error) {
	if !c.started {
		return "", fmt.Errorf("component %s not started in RenderManifest", c.Name())
	}
	return renderManifest(c.CommonComponentFields)
}

// Name implements the IstioComponent interface.
func (c *KialiComponent) Name() name.ComponentName {
	return c.CommonComponentFields.name
}

// CNIComponent is the egress gateway component.
type CNIComponent struct {
	*CommonComponentFields
}

// NewCNIComponent creates a new IngressComponent and returns a pointer to it.
func NewCNIComponent(opts *Options) *CNIComponent {
	return &CNIComponent{
		&CommonComponentFields{
			Options: opts,
			name:    name.CNIComponentName,
		},
	}
}

// Run implements the IstioComponent interface.
func (c *CNIComponent) Run() error {
	return runComponent(c.CommonComponentFields)
}

// RenderManifest implements the IstioComponent interface.
func (c *CNIComponent) RenderManifest() (string, error) {
	if !c.started {
		return "", fmt.Errorf("component %s not started in RenderManifest", c.Name())
	}
	return renderManifest(c.CommonComponentFields)
}

// Name implements the IstioComponent interface.
func (c *CNIComponent) Name() name.ComponentName {
	return c.CommonComponentFields.name
}

// TracingComponent is the egress gateway component.
type TracingComponent struct {
	*CommonComponentFields
}

// NewTracingComponent creates a new IngressComponent and returns a pointer to it.
func NewTracingComponent(opts *Options) *TracingComponent {
	return &TracingComponent{
		&CommonComponentFields{
			Options: opts,
			name:    name.TracingComponentName,
		},
	}
}

// Run implements the IstioComponent interface.
func (c *TracingComponent) Run() error {
	return runComponent(c.CommonComponentFields)
}

// RenderManifest implements the IstioComponent interface.
func (c *TracingComponent) RenderManifest() (string, error) {
	if !c.started {
		return "", fmt.Errorf("component %s not started in RenderManifest", c.Name())
	}
	return renderManifest(c.CommonComponentFields)
}

// Name implements the IstioComponent interface.
func (c *TracingComponent) Name() name.ComponentName {
	return c.CommonComponentFields.name
}

// runComponent performs startup tasks for the component defined by the given CommonComponentFields.
func runComponent(c *CommonComponentFields) error {
	r, err := createHelmRenderer(c)
	if err != nil {
		return err
	}
	if err := r.Run(); err != nil {
		return err
	}
	c.renderer = r
	c.started = true
	return nil
}

// TranslateHelmValues creates a Helm values.yaml config data tree from icp using the given translator.
func TranslateHelmValues(icp *v1alpha2.IstioControlPlaneSpec, translator *translate.Translator, componentName name.ComponentName) (string, error) {
	globalVals, apiVals := make(map[string]interface{}), make(map[string]interface{})

	// First, translate the IstioControlPlane API to helm Values.
	apiValsStr, err := translator.ProtoToValues(icp)
	if err != nil {
		return "", err
	}
	err = yaml.Unmarshal([]byte(apiValsStr), &apiVals)
	if err != nil {
		return "", err
	}
	log.Infof("Values translated from IstioControlPlane API:\n%s", apiValsStr)

	// Add global overlay from IstioControlPlaneSpec.Values.
	_, err = name.SetFromPath(icp, "Values", &globalVals)
	if err != nil {
		return "", err
	}
	log.Infof("Values from IstioControlPlaneSpec.Values:\n%s", util.ToYAML(globalVals))

	mergedVals, err := overlayTrees(globalVals, apiVals)
	if err != nil {
		return "", err
	}

	cns := []name.ComponentName{componentName}
	if componentName == "" {
		cns = translator.AllComponentsNames()
	}

	for _, cn := range cns {
		componentVals, componentValsUnvalidated := make(map[string]interface{}), make(map[string]interface{})
		// Add overlay from IstioControlPlaneSpec.<Feature>.Components.<Component>.Common.ValuesOverrides.
		pathToValues := fmt.Sprintf("%s.Components.%s.Common.Values", translator.ToFeature[cn], cn)
		_, err = name.SetFromPath(icp, pathToValues, &componentVals)
		if err != nil {
			return "", err
		}

		// Add overlay from IstioControlPlaneSpec.<Feature>.Components.<Component>.Common.UnvalidatedValuesOverrides.
		pathToUnvalidatedValues := fmt.Sprintf("%s.Components.%s.Common.UnvalidatedValues", translator.ToFeature[cn], cn)
		_, err = name.SetFromPath(icp, pathToUnvalidatedValues, &componentValsUnvalidated)
		if err != nil {
			return "", err
		}

		log.Infof("Untranslated values from %s:\n%s", pathToValues, util.ToYAML(componentVals))
		log.Infof("Untranslated values from %s:\n%s", pathToUnvalidatedValues, util.ToYAML(componentValsUnvalidated))

		// Translate from path in the API to helm paths.
		componentVals = translator.ValuesOverlaysToHelmValues(componentVals, cn)
		componentValsUnvalidated = translator.ValuesOverlaysToHelmValues(componentValsUnvalidated, cn)

		log.Infof("Translated values from %s:\n%s", pathToValues, util.ToYAML(componentVals))
		log.Infof("Translated values from %s:\n%s", pathToUnvalidatedValues, util.ToYAML(componentValsUnvalidated))

		mergedVals, err = overlayTrees(mergedVals, componentVals, componentValsUnvalidated)
		if err != nil {
			return "", err
		}
	}

	mergedYAML, err := yaml.Marshal(mergedVals)
	if err != nil {
		return "", err
	}
	return string(mergedYAML), err
}

// renderManifest renders the manifest for the component defined by c and returns the resulting string.
func renderManifest(c *CommonComponentFields) (string, error) {
	e, err := c.Translator.IsComponentEnabled(c.name, c.InstallSpec)
	if err != nil {
		return "", err
	}
	if !e {
		return disabledYAMLStr(c.name), nil
	}

	mergedYAML, err := TranslateHelmValues(c.InstallSpec, c.Translator, c.name)
	if err != nil {
		return "", err
	}

	log.Infof("Merged values:\n%s\n", mergedYAML)

	my, err := c.renderer.RenderManifest(mergedYAML)
	if err != nil {
		log.Errorf("Error rendering the manifest: %s", err)
		return "", err
	}
	my += helm.YAMLSeparator + "\n"
	log.Infof("Initial manifest with merged values:\n%s\n", my)

	// Add the k8s resources from IstioControlPlaneSpec.
	my, err = c.Translator.OverlayK8sSettings(my, c.InstallSpec, c.name)
	if err != nil {
		log.Errorf("Error in OverlayK8sSettings: %s", err)
		return "", err
	}
	my = "# Resources for " + string(c.name) + " component\n\n" + my
	log.Infof("Manifest after k8s API settings:\n%s\n", my)

	// Add the k8s resource overlays from IstioControlPlaneSpec.
	pathToK8sOverlay := fmt.Sprintf("%s.Components.%s.Common.K8S.Overlays", c.FeatureName, c.name)
	var overlays []*v1alpha2.K8SObjectOverlay
	found, err := name.SetFromPath(c.InstallSpec, pathToK8sOverlay, &overlays)
	if err != nil {
		return "", err
	}
	if !found {
		return my, nil
	}
	kyo, err := yaml.Marshal(overlays)
	if err != nil {
		return "", err
	}
	log.Infof("Applying kubernetes overlay: \n%s\n", kyo)
	ns, err := name.Namespace(c.FeatureName, c.name, c.InstallSpec)
	if err != nil {
		return "", err
	}
	ret, err := patch.YAMLManifestPatch(my, ns, overlays)
	if err != nil {
		return "", err
	}

	log.Infof("After overlay: \n%s\n", ret)
	return ret, nil
}

// overlayTrees overlays component validated and unvalidated values over base.
func overlayTrees(base map[string]interface{}, overlays ...map[string]interface{}) (map[string]interface{}, error) {
	bby, err := yaml.Marshal(base)
	if err != nil {
		return nil, err
	}
	by := string(bby)

	for _, o := range overlays {
		oy, err := yaml.Marshal(o)
		if err != nil {
			return nil, err
		}

		by, err = helm.OverlayYAML(by, string(oy))
		if err != nil {
			return nil, err
		}
	}

	out := make(map[string]interface{})
	err = yaml.Unmarshal([]byte(by), &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// createHelmRenderer creates a helm renderer for the component defined by c and returns a ptr to it.
func createHelmRenderer(c *CommonComponentFields) (helm.TemplateRenderer, error) {
	icp := c.InstallSpec
	ns, err := name.Namespace(c.FeatureName, c.name, c.InstallSpec)
	if err != nil {
		return nil, err
	}
	return helm.NewHelmRenderer(icp.InstallPackagePath, c.Translator.ComponentMaps[c.name].HelmSubdir,
		string(c.name), ns)
}

// disabledYAMLStr returns the YAML comment string that the given component is disabled.
func disabledYAMLStr(componentName name.ComponentName) string {
	return yamlCommentStr + string(componentName) + componentDisabledStr + "\n"
}
