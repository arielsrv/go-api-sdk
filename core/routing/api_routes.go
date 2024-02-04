package routing

import "gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core/container"

type APIRoutes struct {
	routes  []Route
	statics []Static
}

func (r *APIRoutes) AddRoute(method string, path string, action func(ctx *HTTPContext) error) {
	r.routes = append(r.routes, Route{
		Method: method,
		Path:   path,
		Action: action,
	})
}

func (r *APIRoutes) AddStatic(prefix string, route string) {
	r.statics = append(r.statics, Static{
		Prefix: prefix,
		Root:   route,
	})
}

func (r *APIRoutes) GetRoutes() []Route {
	return r.routes
}

func (r *APIRoutes) GetStatics() []Static {
	return r.statics
}

func To[T any]() T {
	return container.Provide[T]()
}
