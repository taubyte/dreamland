package http

import (
	"context"
	"testing"
	"time"

	"github.com/taubyte/dreamland/service/api"
	"github.com/taubyte/dreamland/service/inject"
	commonIface "github.com/taubyte/go-interfaces/common"
	"github.com/taubyte/tau/libdream"
	_ "github.com/taubyte/tau/protocols/auth"

	_ "github.com/taubyte/tau/protocols/hoarder"
	_ "github.com/taubyte/tau/protocols/monkey"
	_ "github.com/taubyte/tau/protocols/patrick"
	_ "github.com/taubyte/tau/protocols/seer"
	_ "github.com/taubyte/tau/protocols/tns"

	_ "github.com/taubyte/tau/clients/p2p/monkey"
	_ "github.com/taubyte/tau/clients/p2p/patrick"
	_ "github.com/taubyte/tau/clients/p2p/tns"
)

func TestRoutes(t *testing.T) {
	univerName := "dreamland-http"
	// start multiverse
	err := api.BigBang()
	if err != nil {
		t.Errorf("Failed big bang with error: %v", err)
		return
	}

	u := libdream.New(libdream.UniverseConfig{Name: univerName})
	defer u.Stop()

	err = u.StartWithConfig(&libdream.Config{
		Services: map[string]commonIface.ServiceConfig{
			"monkey":  {},
			"auth":    {},
			"patrick": {},
			"seer":    {},
			"hoarder": {},
			"tns":     {},
		},
		Simples: map[string]libdream.SimpleConfig{
			"client": {
				Clients: libdream.SimpleConfigClients{
					Monkey:  &commonIface.ClientConfig{},
					Patrick: &commonIface.ClientConfig{},
					TNS:     &commonIface.ClientConfig{},
				}.Compat(),
			},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()

	time.Sleep(2 * time.Second)

	client, err := New(ctx, URL("http://localhost:1421"), Timeout(60*time.Second))
	if err != nil {
		t.Errorf("Failed creating http client error: %v", err)
		return
	}

	universe := client.Universe(univerName)

	// Create simple called test1
	err = universe.Inject(inject.Simple("test1", &libdream.SimpleConfig{}))
	if err != nil {
		t.Errorf("Failed simples call with error: %v", err)
		return
	}

	time.Sleep(2 * time.Second)

	// Should not fail
	_, err = u.Simple("test1")
	if err != nil {
		t.Errorf("Failed getting simple with error: %v", err)
		return
	}

	// Should fail
	_, err = u.Simple("dne")
	if err == nil {
		t.Error("Should have failed, expecting to not find dne simple node")
		return
	}

	// Should fail
	err = universe.Inject(inject.Fixture("should fail", "dne"))
	if err == nil {
		t.Error("Expecting fail for fixture not existing")
		return
	}

	test, err := client.Status()
	if err != nil {
		t.Error(err)
		return
	}
	_, ok := test[univerName]
	if ok == false {
		t.Error("Did not find universe in status")
		return
	}

}
