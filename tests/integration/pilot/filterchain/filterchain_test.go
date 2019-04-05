package filterchain

import (
	"fmt"
	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/apps"
	"istio.io/istio/pkg/test/framework/components/environment"
	pilotcomponent "istio.io/istio/pkg/test/framework/components/pilot"
	"istio.io/istio/pkg/test/util/retry"
	"testing"
	"time"
)


func TestMain(m *testing.M) {

	// TODO
	// No need to enforce kube for filter chain
	// set up envoy
	// set up client and app
	// set up pilot
	print("In TestMain")
	framework.Main("filterchain_test", m)

	// NOTE: print(m) only run once no matter how many TestXX is given
	print(m)
	// no need to explicitly invoke os.Exit as it's part of framework.Main
	// os.Exit(m.Run())
}

func TestBar(t *testing.T) {
	t.Log("in TestBar")
}

func TestFoo(t *testing.T) {
	t.Log("in TestFoo")
	return
	ctx := framework.NewContext(t)
	defer ctx.Done(t)

	ctx.RequireOrSkip(t, environment.Kube)

	pilot := pilotcomponent.NewOrFail(t, ctx, pilotcomponent.Config{})

	applications := apps.NewOrFail(ctx, t, apps.Config{Pilot: pilot})

	aApp := applications.GetAppOrFail("a", t).(apps.KubeApp)
	bApp := applications.GetAppOrFail("b", t).(apps.KubeApp)

	for _, e := range bApp.Endpoints() {
		t.Logf("b end point %s", e)
	}

	retry.UntilSuccessOrFail(t, func() error {
		bEndpoint := bApp.EndpointForPort(80)
		bProtocol := apps.AppProtocolHTTP
		results, err := aApp.Call(bEndpoint,
			apps.AppCallOptions{
				Protocol: apps.AppProtocol(bProtocol),
			})
		if err != nil || len(results) == 0 || results[0].Code != "200" {
			// Addition log for debugging purpose.
			if err != nil {
				fmt.Printf("Error: %#v\n", err)
			} else if len(results) == 0 {
				fmt.Printf("No result\n")
			} else {
				fmt.Printf("Result: %v\n", results[0])
			}
			return fmt.Errorf("%s to %s:%d using %s: expected success, actually failed",
				aApp.Name(), bApp.Name(), 80, bProtocol)
		}
		return nil
	}, retry.Delay(time.Second), retry.Timeout(10*time.Second))

	t.Logf("stop here")

}
