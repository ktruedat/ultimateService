package mid

import (
	"context"
	"github.com/ktruedat/ultimateService/business/sys/metrics"
	"github.com/ktruedat/ultimateService/foundation/web"
	"net/http"
)

// Metrics updates app counters.
func Metrics() web.Middleware {

	// This is the actual middleware function to be executed.
	m := func(handler web.Handler) web.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// Add the metrics into the context for metric gathering.
			ctx = metrics.Set(ctx)

			// Call the handler of the route.
			err := handler(ctx, w, r)

			// Handle updating the metrics that handled here.

			// Increase the request and goroutines counter.
			metrics.AddRequests(ctx)
			metrics.AddGoroutines(ctx)

			// Increment if there is an error flowing through the request.

			if err != nil {
				metrics.AddErrors(ctx)
			}

			return err
		}
		return h
	}
	return m
}
