package middleware

import (
	"context"
	"git.samberi.com/dois/delivery_api/config"
	"github.com/gorilla/sessions"
	"net/http"
)

func SessionsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var store = sessions.NewFilesystemStore("", []byte(config.Config.Session.Key))
		store.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   config.Config.Session.MaxAge,
			HttpOnly: true,
		}
		context.WithValue(r.Context(), "session", store)
		h.ServeHTTP(w, r)
	})
}
