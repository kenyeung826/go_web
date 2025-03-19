package middleware

import (
	"net/http"
	"runtime/debug"

	"app/error"

	fwMiddleware "github.com/go-chi/chi/v5/middleware"
)

func ErrorHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if _, ok := rvr.(*error.ApplicationError); ok {
					logEntry := fwMiddleware.GetLogEntry(r)
					if logEntry != nil {
						logEntry.Panic(rvr, debug.Stack())
					} else {
						fwMiddleware.PrintPrettyStack(rvr)
					}
					http.Redirect(w, r, "/error/500", http.StatusFound)
				} else {
					panic(rvr)
				}
			}
		}()

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
