package api

import (
	"context"
	"time"

	goHttp "net/http"

	"github.com/pterm/pterm"
	httpIface "github.com/taubyte/http"
	http "github.com/taubyte/http/basic"
	"github.com/taubyte/http/options"
	"github.com/taubyte/tau/libdream/common"
	"github.com/taubyte/tau/libdream/services"
)

type multiverseService struct {
	rest httpIface.Service
	common.Multiverse
}

// creates function called 'BigBang' which starts Dreamland server and returns error if unsuccessful
func BigBang() error {
	//initializes variable named 'err' to store and handle errors
	var err error

	//creates an instance of the multiverseService struct named 'srv'
	srv := &multiverseService{
		//initializes the Multiverse field of the 'srv' struct by creating a new multiverse instance
		Multiverse: services.NewMultiVerse(),
	}

	//creates an HTTP server and assigns it to the 'rest' field of the 'srv' struct
	//sets server to listen on the address specified in common.DreamlandApiListen
	//allows requests from the specified origins
	srv.rest, err = http.New(srv.Context(), options.Listen(common.DreamlandApiListen), options.AllowedOrigins(true, []string{".*"}))
	if err != nil {
		return err
	}

	//sets up the HTTP routes and handlers and starts the server
	srv.setUpHttpRoutes().Start()

	//sets a timeout of 10 seconds for the server to start
	waitCtx, waitCtxC := context.WithTimeout(srv.Context(), 10*time.Second)
	defer waitCtxC()

	//for loop that checks the status of Dreamland until it starts successfully or times out
	for {
		select {
		//checks for timeout and returns an error if timeout is reached
		case <-waitCtx.Done():
			return waitCtx.Err()
		//checks if there was a server error and prints error message if so
		case <-time.After(100 * time.Millisecond):
			if srv.rest.Error() != nil {
				pterm.Error.Println("Dreamland failed to start")
				return srv.rest.Error()
			}
			_, err := goHttp.Get("http://" + common.DreamlandApiListen)
			//if server is responding and there is no error, prints message indicating Dreamland is ready
			if err == nil {
				pterm.Info.Println("Dreamland ready")
				return nil
			}
		}
	}
}

/* The above code could be rewritten as follows:

type ServerConfig struct {
	Address        string
	AllowedOrigins []string
	Timeout        time.Duration
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Address:        common.DreamlandApiListen,
		AllowedOrigins: []string{".*"},
		Timeout:        10 * time.Second,
	}
}

func StartServer(config *ServerConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	srv := &multiverseService{
		Multiverse: services.NewMultiVerse(),
	}

	srv.rest, err := http.New(srv.Context(), options.Listen(config.Address), options.AllowedOrigins(true, config.AllowedOrigins))
	if err != nil {
		return fmt.Errorf("failed to create HTTP server: %v", err)
	}

	srv.setUpHttpRoutes().Start()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(100 * time.Millisecond):
		if srv.rest.Error() != nil {
			pterm.Error.Println("Dreamland failed to start")
			return srv.rest.Error()
		}
		_, err := http.Get("http://" + config.Address)
		if err == nil {
			pterm.Info.Println("Dreamland ready")
			return nil
		}
	}

	return nil
}*/

//This version is more modular and readable and provides more detailed error handling messages. However, it is admittedly less fun than the original.
