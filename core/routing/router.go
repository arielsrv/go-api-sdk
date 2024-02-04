package routing

type Router interface {
	Register()
	GetRoutes() []Route
	GetStatics() []Static
}
