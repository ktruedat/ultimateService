package mid

import (
	"context"
	"fmt"
	"github.com/ktruedat/ultimateService/foundation/web"
	"net/http"
	"runtime/debug"
)

// Panics recovers from panics and converts the panic to an error, so it is
// reported in Metrics and handled in Errors.
func Panics() web.Middleware {

	// This is the actual middleware function to be executed.
	m := func(handler web.Handler) web.Handler {

		// Create the handler that will be attached in the middleware chain.
		// Now, this function signature returns a named variable err of type error
		// Go has this feature (named return types) and this is a really rare case
		// when we should use it to return the error if a panic is registered.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {

			// Defer a function to recover from a panic and set the err return
			// variable after the fact.
			defer func() {
				if rec := recover(); rec != nil {

					// we can use the debug.Stack() function to print the stack trace
					// in the error: "PANIC [%v] TRACE [%v]", rec, string(trace)
					trace := debug.Stack()
					// Stack trace will be provided
					err = fmt.Errorf("PANIC [%v] TRACE [%v]", rec, string(trace))
				}
			}()

			// Call the next handler and set its return value in the err variable.
			return handler(ctx, w, r)
		}
		return h
	}
	return m
}
