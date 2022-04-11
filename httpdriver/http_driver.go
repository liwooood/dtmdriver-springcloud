package httpdriver

type HttpDriver interface {
	// ResolveHttpService to resolve client url to http url
	ResolveHttpService(serviceUrl string) (string, error)
	// RegisterHttpService register endpoint to target with some params, such as username, password.
	RegisterHttpService(target string, endpoint string, options map[string]string) error
}
