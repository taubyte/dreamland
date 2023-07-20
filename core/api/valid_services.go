package api

import (
	httpIface "github.com/taubyte/http"
)

func validServices() {
	api.GET(&httpIface.RouteDefinition{
		Path:    "/spec/services",
		Handler: servicesHandler,
	})
}

func servicesHandler(ctx httpIface.Context) (interface{}, error) {
	return mv.ValidServices(), nil
}
