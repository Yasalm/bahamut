package maker

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Api struct {
	Address string
	Port    int
	Router  *chi.Mux
	Maker   *Maker
}

type ErrResposne struct {
	HTTPStatusCode int
	Message        string
}

// start;stop;schedule;
func (api *Api) initRouter() {
	api.Router = chi.NewRouter()
	api.Router.Route("/components", func(r chi.Router) {
		r.Post("/", api.scheduleHandler)
		// r.Get("/", api.getComponents)
		r.Route("/{ID}", func(r chi.Router) {
			r.Delete("/", api.stopHandler)
		})
		// r.Route("/stats", func(r chi.Router) {
		// 	r.Get("/", api.getStatus)
		// })
	})
}

func (api *Api) Start() {
	api.initRouter()
	http.ListenAndServe(fmt.Sprintf("%s:%d", api.Address, api.Port), api.Router)
}
