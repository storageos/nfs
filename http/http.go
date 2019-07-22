package http

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"sync"
)

// HTTP manages the web server that serves metrics.
type HTTP struct {
	name     string
	server   *http.Server
	handlers map[string]string
	mu       *sync.RWMutex
}

// New creates a new HTTP server which can be Run and Closed.
//
// Endpoint handlers should be registered using RegisterHandler().
func New(listenAddr string, name string) *HTTP {
	return &HTTP{
		name: name,
		server: &http.Server{
			Addr: listenAddr,
		},
		handlers: make(map[string]string),
		mu:       &sync.RWMutex{},
	}
}

// Run starts the HTTP server, returning an immediate error and nil channel if
// the process cannot be started.
//
// If the process is successfully started, the returned channel will have the
// exit error pushed to it (which may be any error, but exec.ExitError is
// returned typically.)
//
// Once the process stops, the returned channel is closed.
func (h *HTTP) Run() (<-chan error, error) {
	errCh := make(chan error)
	go func() {
		err := h.server.ListenAndServe()
		if err != nil {
			errCh <- err
			return
		}
		close(errCh)
	}()

	return errCh, nil
}

// Close stops the HTTP server.
//
// Once the server has stopped, the channel returned from Run is closed.
func (h *HTTP) Close(ctx context.Context) {
	if h.server != nil {
		if err := h.server.Shutdown(ctx); err != nil {
			log.Printf("error shutting down http server: %v", err)
		}
	}
}

// RegisterHandler registers an HTTP handler for an endpoint.
//
// The name is used as an optional human-readable name for the endpoint.
func (h *HTTP) RegisterHandler(name string, endpoint string, handler http.Handler) {

	h.mu.Lock()
	defer h.mu.Unlock()

	// Only register once.
	if _, ok := h.handlers[endpoint]; ok {
		return
	}

	// Keep our own record of endpoints.
	h.handlers[endpoint] = name

	http.Handle(endpoint, handler)

}

// Handler returns the default HTTP handler.
//
// This handler typically respondes to the "/" endpoint and generates a list of
// available endpoints that have been registered with RegisterHandler().
func (h *HTTP) Handler() http.Handler {

	const index = `<html>
<head><title>{{.Name}}</title></head>
<body>
<h1>{{.Name}}</h1>
{{- range $endpoint, $name := .Endpoints}}
<p><a href="{{$endpoint}}">{{$name}}</a></p>
{{- end }}
</body>
</html>`

	type data struct {
		Name      string
		Endpoints map[string]string
	}
	config := data{
		Name:      h.name,
		Endpoints: h.handlers,
	}

	t, err := template.New("response").Parse(index)
	if err != nil {
		log.Fatal(err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.mu.RLock()
		defer h.mu.RUnlock()

		if err := t.Execute(w, config); err != nil {
			log.Printf("failed writing http response: %v", err)
		}
	})
}
