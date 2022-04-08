package httpdriver

type HttpDriver interface {
	// RegisterHttpService register endpoint to target with some params, such as username, password.
	RegisterHttpService(target string, endpoint string, options map[string]string, paths []string) error
}
