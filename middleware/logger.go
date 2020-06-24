package middleware

import (
	"context"
	"fmt"
	"git.samberi.com/dois/delivery_api/common"
	"git.samberi.com/dois/delivery_api/config"
	"git.samberi.com/dois/delivery_api/logger"
	"git.samberi.com/dois/delivery_api/models"
	"git.samberi.com/dois/delivery_api/utils"
	"google.golang.org/grpc"
	"net/http"
)

func LoggerMiddlewareHTTP(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				if config.Config.System.Debug {
					logger.Logger.Errorf(
						"500 IP:%s - %s: %s%s - %#v",
						utils.GetIPAddress(r),
						r.Method, r.Host, r.URL.Path, rec,
					)
					panic(rec)
				} else {
					common.InternalServerError(w, r, rec)
				}
			}
		}()
		h.ServeHTTP(w, r)
		who := fmt.Sprintf("IP:%s", utils.GetIPAddress(r))
		if user := r.Context().Value("user"); user != nil {
			who = fmt.Sprintf("%s - %s", who, user.(*models.User).String())
		}
		logger.Logger.Infof("%s - %s: %s%s", who, r.Method, r.Host, r.URL.Path)
	})
}

func LoggerMiddlewareGRPC() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if rec := recover(); rec != nil {
				if config.Config.System.Debug {
					panic(rec)
				}
				logger.Logger.Errorf(
					"GRPC FATAL ERROR %s: - %#v",
					info.FullMethod, rec,
				)
			}
		}()
		resp, err = handler(ctx, req)
		if err != nil {
			logger.Logger.Errorf("%s, %s", info.FullMethod, err.Error())
		} else {
			logger.Logger.Infof("%s", info.FullMethod)
		}
		return resp, err
	}
}
