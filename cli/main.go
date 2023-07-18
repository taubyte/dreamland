package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	// Relative
	"github.com/taubyte/dreamland/cli/common"
	inject "github.com/taubyte/dreamland/cli/inject"
	"github.com/taubyte/dreamland/cli/kill"
	"github.com/taubyte/dreamland/cli/new"
	"github.com/taubyte/dreamland/cli/status"

	// Actual imports

	"bitbucket.org/taubyte/p2p/peer"
	"github.com/taubyte/dreamland/core/services"
	client "github.com/taubyte/dreamland/http"
	"github.com/urfave/cli/v2"

	// Empty imports for initializing fixtures, and client/service run methods
	_ "bitbucket.org/taubyte/billing/api/p2p"
	_ "bitbucket.org/taubyte/billing/service"
	_ "bitbucket.org/taubyte/console/api/p2p"
	_ "bitbucket.org/taubyte/console/ui/service"
	_ "github.com/taubyte/dreamland/fixtures"
	_ "github.com/taubyte/odo/clients/p2p/auth"
	_ "github.com/taubyte/odo/clients/p2p/hoarder"
	_ "github.com/taubyte/odo/clients/p2p/monkey"
	_ "github.com/taubyte/odo/clients/p2p/patrick"
	_ "github.com/taubyte/odo/clients/p2p/seer"
	_ "github.com/taubyte/odo/clients/p2p/tns"
	_ "github.com/taubyte/odo/protocols/auth/service"
	_ "github.com/taubyte/odo/protocols/hoarder/service"
	_ "github.com/taubyte/odo/protocols/monkey/fixtures/compile"
	_ "github.com/taubyte/odo/protocols/monkey/service"
	_ "github.com/taubyte/odo/protocols/node/service"
	_ "github.com/taubyte/odo/protocols/patrick/service"
	_ "github.com/taubyte/odo/protocols/seer/service"
	_ "github.com/taubyte/odo/protocols/tns/service"
)

func main() {
	peer.DevMode = true
	ctx, ctxC := context.WithCancel(context.Background())

	defer func() {
		if common.DoDaemon {
			ctxC()
			services.Zeno()
		}
	}()

	multiverse, err := client.New(ctx, client.URL(common.DefaultDreamlandURL), client.Timeout(300*time.Second))
	if err != nil {
		log.Fatalf("Starting new dreamland client failed with: %s", err.Error())
	}

	err = defineCLI(&common.Context{Ctx: ctx, Multiverse: multiverse}).RunContext(ctx, os.Args)
	if err != nil {
		log.Fatal(err)
	}

	return
}

func defineCLI(ctx *common.Context) *(cli.App) {
	app := &cli.App{
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			new.Command(ctx),
			inject.Command(ctx),
			kill.Command(ctx),
			status.Command(ctx),
		},
		Suggest:              true,
		EnableBashCompletion: true,
	}

	return app
}

func setupTrap(ctxC context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		for range c {
			ctxC()
			// sig is a ^C, handle it
		}
	}()
}
