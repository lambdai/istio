package filterchain

import (
	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/apps"
	"istio.io/istio/pkg/test/framework/components/environment"
	pilotcomponent "istio.io/istio/pkg/test/framework/components/pilot"
	"testing"
)


func TestMain(m *testing.M) {
	framework.Main("filterchain_test", m)
	// TODO
	// enforce kube
	// set up envoy
	// set up client and app
	// set up pilot
}

func TestFoo(t *testing.T) {
	ctx := framework.NewContext(t)
	defer ctx.Done(t)

	// TODO kube only
	ctx.RequireOrSkip(t, environment.Native)

	t.Logf("Before pilot")
	pilot := pilotcomponent.NewOrFail(t, ctx, pilotcomponent.Config{})

	t.Logf("After pilot")
	applications := apps.NewOrFail(ctx, t, apps.Config{Pilot: pilot})

	t.Logf("Before Get app A")
	appA := applications.GetAppOrFail("a", t)

	t.Logf( "appA name : %s\n", appA.Name())

}
