package httpdriver

type HttpDriver interface {
	// RegisterHttpResolver register the http resolver to handle custom scheme
	RegisterHttpResolver()
	// RegisterHttpService register endpoint to target with some params, such as username, password.
	RegisterHttpService(target string, endpoint string, options map[string]string, paths []string) error
}
