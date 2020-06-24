package main

import (
	"context"
	"fmt"
	"git.samberi.com/dois/delivery_api/admin"
	"git.samberi.com/dois/delivery_api/common"
	"git.samberi.com/dois/delivery_api/config"
	"git.samberi.com/dois/delivery_api/database"
	"git.samberi.com/dois/delivery_api/deliveryman"
	"git.samberi.com/dois/delivery_api/logger"
	"git.samberi.com/dois/delivery_api/middleware"
	"git.samberi.com/dois/delivery_api/proto/gen"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_runtime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"io"
	"math/rand"
	"mime"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"
)

var GRPCServer *grpc.Server
var HTTPProxyServer *http.Server

func Init() {
	config.Init()
	logger.Init()
	database.MigrateDB()
}

func serveStatic(prefix, name string, router *mux.Router) {
	mime.AddExtensionType(".svg", "image/svg+xml")
	fileServer := http.FileServer(http.Dir(path.Join(config.Config.Static.StaticRoot, name)))
	router.PathPrefix(prefix).Handler(http.StripPrefix(prefix, fileServer))
}

func StartGRPCServer() {
	defer func() {
		fmt.Println("GRPC Server closed")
	}()
	GRPCServer = grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			middleware.LoggerMiddlewareGRPC(),
			middleware.DbSetContextUnary(),
			grpc_auth.UnaryServerInterceptor(common.CheckAuthByToken),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)
	gen.RegisterAdminServer(GRPCServer, admin.NewSever())
	gen.RegisterApiServer(GRPCServer, deliveryman.NewServer())
	listener, err := net.Listen("tcp", config.Config.System.GRPCServer)
	if err != nil {
		logger.HandleError(err)
		return
	}
	logger.Logger.Infof("Run GRPC server %s", config.Config.System.GRPCServer)
	logger.HandleError(GRPCServer.Serve(listener))
}

func StartHTTPProxyServer() {
	registerEndpointsHandlers := []func(context.Context, *grpc_runtime.ServeMux, string, []grpc.DialOption) error{
		gen.RegisterApiHandlerFromEndpoint,
		gen.RegisterAdminHandlerFromEndpoint,
	}
	cancels := make([]func(), 0)
	defer func() {
		fmt.Println("HTTP ProxyServer closed")
		for _, cancel := range cancels {
			cancel()
		}
	}()
	gmux := grpc_runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	for _, registerEndpointsHandler := range registerEndpointsHandlers {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		cancels = append(cancels, cancel)
		err := registerEndpointsHandler(ctx, gmux, config.Config.System.GRPCServer, opts)
		if err != nil {
			logger.HandleError(err)
			return
		}
	}
	router := mux.NewRouter()
	router.HandleFunc("/schema/{file_name}", func(w http.ResponseWriter, r *http.Request) {
		if fileName, ok := mux.Vars(r)["file_name"]; ok {
			if f, err := os.Open(path.Join(config.Config.Static.StaticRoot, "schema", fileName)); err == nil {
				io.Copy(w, f)
			}
		}
	})
	_router := router.NewRoute().Subrouter()
	_router.Use(middleware.LoggerMiddlewareHTTP)
	serveStatic("/swagger-ui/", "swagger_ui", _router)
	serveStatic("/admin-ui/", "admin_ui", _router)
	router.PathPrefix("/").Handler(gmux)
	HTTPProxyServer = &http.Server{
		Addr:    config.Config.System.HTTPServer,
		Handler: router,
	}
	logger.Logger.Infof("Run ProxyServer %s", config.Config.System.HTTPServer)
	logger.HandleError(HTTPProxyServer.ListenAndServe())
}

func main() {
	rand.Seed(time.Now().UnixNano())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		_ = HTTPProxyServer.Close()
		GRPCServer.Stop()
	}()
	Init()
	go StartGRPCServer()
	StartHTTPProxyServer()
}
