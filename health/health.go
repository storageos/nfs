package health

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/storageos/nfs/ganesha"
)

// Health handles health collection and presentation.
type Health struct {
	ganesha *ganesha.Ganesha
}

// New creates a new health instance.
func New(nfs *ganesha.Ganesha) *Health {
	return &Health{
		ganesha: nfs,
	}
}

// Handler returns an http handler for reporting health.
//
// The endpoint will return 200/OK when the NFS server is operational and
// publishing heartbeats.
func (h *Health) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if h.ganesha.IsReady(ctx) {
			w.WriteHeader(200)
			w.Write([]byte("ok"))
		} else {
			w.WriteHeader(503)
			w.Write([]byte(fmt.Sprintf("nfs server not ready")))
		}
	})
}
