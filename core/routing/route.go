package routing

type Route struct {
	Method string
	Path   string
	Action func(ctx *HTTPContext) error
}

type Static struct {
	Prefix string
	Root   string
}
