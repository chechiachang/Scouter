package server

import (
	"github.com/chechiachang/scouter/serviceprovider"
	"github.com/emicklei/go-restful"
	"github.com/gorilla/mux"
)

// AppRoute will add router
func (a *Apiserver) AppRoute() *mux.Router {
	router := mux.NewRouter()

	container := restful.NewContainer()

	container.Filter(globalLogging)

	container.Add(newVersionService(a.ServiceProvider))

	router.PathPrefix("/v1/").Handler(container)
	return router
}

func newVersionService(sp *serviceprovider.Container) *restful.WebService {
	webService := new(restful.WebService)
	webService.Path("/v1/version").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	//  webService.Filter(validateTokenMiddleware)
	webService.Route(webService.GET("/").To(RESTfulServiceHandler(sp, versionHandler)))
	return webService
}

// RESTfulContextHandler is the interface for restfuul handler(restful.Request,restful.Response)
type RESTfulContextHandler func(*Context)

// RESTfulServiceHandler is the wrapper to combine the RESTfulContextHandler with our serviceprovider object
func RESTfulServiceHandler(sp *serviceprovider.Container, handler RESTfulContextHandler) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
		ctx := Context{
			ServiceProvider: sp,
			Request:         req,
			Response:        resp,
		}
		handler(&ctx)
	}
}