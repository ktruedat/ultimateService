package mid

import (
	"context"
	"github.com/ktruedat/ultimateService/foundation/web"
	"go.uber.org/zap"
	"net/http"
)

// Logger writes some information about the request to the logs in the
// format: TraceID : (200) GET /foo -> IP ADDR (latency)
func Logger(log *zap.SugaredLogger) web.Middleware {

	// This is the actual middleware function to be executed.
	m := func(handler web.Handler) web.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			log.Infow("request started", "method", r.Method, "path", r.URL.Path,
				"remoteaddr", r.RemoteAddr)

			err := handler(ctx, w, r)

			log.Infow("request completed", "method", r.Method, "path", r.URL.Path,
				"remoteaddr", r.RemoteAddr)
			return err
		}
		return h
	}
	return m
}
