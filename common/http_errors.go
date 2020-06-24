package common

import (
	"git.samberi.com/dois/delivery_api/logger"
	"git.samberi.com/dois/delivery_api/utils"
	"net/http"
)

func InternalServerError(w http.ResponseWriter, r *http.Request, rec interface{}) {
	w.WriteHeader(http.StatusInternalServerError)
	_, e := w.Write([]byte("500 Internal Server Error"))
	logger.Logger.Errorf("500 - IP:%s - %s: %s%s - %v", utils.GetIPAddress(r), r.Method, r.Host, r.URL.Path, rec)
	logger.HandleError(e)
}

func Error404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_, err := w.Write([]byte("Page not found"))
	logger.Logger.Warningf("404 - IP:%s - %s: %s%s", utils.GetIPAddress(r), r.Method, r.Host, r.URL.Path)
	logger.HandleError(err)
}

func Error405(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	_, err := w.Write([]byte("Method not allowed"))
	logger.Logger.Warningf("405 - IP:%s - %s: %s%s", utils.GetIPAddress(r), r.Method, r.Host, r.URL.Path)
	logger.HandleError(err)
}

func Forbidden(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
	logger.Logger.Warningf("403 - IP:%s - %s: %s%s", utils.GetIPAddress(r), r.Method, r.Host, r.URL.Path)
	_, err := w.Write([]byte("Forbidden"))
	logger.HandleError(err)
}

func Unauthorized(w http.ResponseWriter, r *http.Request, message string) {
	w.WriteHeader(http.StatusUnauthorized)
	logger.Logger.Warningf("401 - IP:%s - %s: %s%s", utils.GetIPAddress(r), r.Method, r.Host, r.URL.Path)
	if message == "" {
		message = "Unauthorized"
	}
	_, err := w.Write([]byte(message))
	logger.HandleError(err)
}
