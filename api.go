package apibuilder

type APIRoutes interface {
	Routes(with *With) []API
}