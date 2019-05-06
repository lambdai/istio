// Copyright 2019 Istio Authors
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

package interception

import (
	"fmt"
	"istio.io/istio/pkg/test/framework/components/apps"
	"istio.io/istio/pkg/test/framework/components/galley"
	"istio.io/istio/pkg/test/framework/components/pilot"
	"istio.io/istio/pkg/test/util/retry"
	"testing"
	"time"

	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/environment"
	"istio.io/istio/pkg/test/framework/components/istio"
)

func TestMain(m *testing.M) {
	var ist istio.Instance
	framework.NewSuite("ping_at_traffic_interception", m).
		SetupOnEnv(environment.Kube, istio.Setup(&ist, nil)).
		Run()
}

type appConnectionPair struct {
	src, dst apps.KubeApp
}

func TestReachablity(t *testing.T) {
	// WTF
	ctx := framework.NewContext(t)
	defer ctx.Done(t)
	g := galley.NewOrFail(t, ctx, galley.Config{})
	p := pilot.NewOrFail(t, ctx, pilot.Config{Galley: g})
	appsInstance := apps.NewOrFail(t, ctx, apps.Config{Pilot: p, Galley: g})

	inoutUnitedApp0, _ := appsInstance.GetAppOrFail("a", t).(apps.KubeApp)
	inoutUnitedApp1, _ := appsInstance.GetAppOrFail("b", t).(apps.KubeApp)

	inoutSplitApp0, _ := appsInstance.GetAppOrFail("c", t).(apps.KubeApp)
	inoutSplitApp1, _ := appsInstance.GetAppOrFail("d", t).(apps.KubeApp)

	connectivityPairs := []appConnectionPair{
		// source is inout united
		{inoutUnitedApp0, inoutUnitedApp1},
		{inoutUnitedApp0, inoutSplitApp1},

		// source is inout split
		{inoutSplitApp0, inoutUnitedApp1},
		{inoutSplitApp0, inoutSplitApp1},

		// self connectivity (is it required?)
		{inoutUnitedApp0, inoutUnitedApp0},
		{inoutSplitApp0, inoutSplitApp0},
	}
	for _, pair := range connectivityPairs {
		src := pair.src
		dst := pair.dst
		retry.UntilSuccessOrFail(t, func() error {
			ep := dst.EndpointForPort(80)
			if ep == nil {
				return fmt.Errorf("cannot get upstream endpoint for connection test %s", dst.Name())
			}

			results, err := src.Call(ep, apps.AppCallOptions{Protocol: apps.AppProtocolHTTP, Path: ""})

			if err != nil || len(results) == 0 || results[0].Code != "200" {
				// Addition log for debugging purpose.
				if err != nil {
					t.Logf("Error: %#v\n", err)
				} else if len(results) == 0 {
					t.Logf("No result\n")
				} else {
					t.Logf("Result: %v\n", results[0])
				}
				return fmt.Errorf("%s cannot connect to %s on port 80 using protocol HTTP", src.Name(), dst.Name())
			}
			return nil
		}, retry.Delay(time.Second), retry.Timeout(10*time.Second))
	}
}
