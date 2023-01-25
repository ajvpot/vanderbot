package systemfx

import (
	"context"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module starts an HTTP server and returns a http.ServeMux to use to register handlers
// for the server.
var Module = fx.Options(fx.Provide(New), fx.Invoke(RegisterHealthCheck))

// Params defines the dependencies of the httpfx module.
type Params struct {
	fx.In

	Lifecycle  fx.Lifecycle
	Shutdowner fx.Shutdowner
	Logger     *zap.Logger
}

// Result defines the objects that the httpfx module provides.
type Result struct {
	fx.Out

	Mux *http.ServeMux `name:"system"`
}

// New exports functionality similar to Module, but allows the caller to wrap
// or modify Result. Most users should use Module instead.
func New(p Params) (Result, error) {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:     "127.0.0.1:8392",
		Handler:  mux,
		ErrorLog: zap.NewStdLog(p.Logger),
	}

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				err := server.ListenAndServe()
				if err != nil && err != http.ErrServerClosed {
					p.Logger.Error("failed to start http server cleanly", zap.Error(err))
					_ = p.Shutdowner.Shutdown()
				}
			}()
			p.Logger.Info("starting HTTP server", zap.String("addr", server.Addr))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			err := server.Shutdown(ctx)
			if err == http.ErrServerClosed {
				return nil
			}
			return err
		},
	})

	return Result{Mux: mux}, nil
}

type RegisterRouteParams struct {
	fx.In

	Mux *http.ServeMux `name:"system"`
}

func RegisterHealthCheck(p RegisterRouteParams) {
	http.NotFoundHandler()
	p.Mux.Handle("/healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
}
