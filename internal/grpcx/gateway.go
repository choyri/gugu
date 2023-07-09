package grpcx

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/choyri/gugu/internal/logx"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/exp/slog"
)

type Gateway struct {
	svr    *http.Server
	stdMux *http.ServeMux
	gwMux  *runtime.ServeMux
	logger *logx.Logger

	withoutAPIPrefix bool
}

type GatewayOption func(*Gateway)

func WithoutAPIPrefix() GatewayOption {
	return func(gw *Gateway) {
		gw.withoutAPIPrefix = true
	}
}

func NewGateway(addr string, opts ...GatewayOption) *Gateway {
	gw := &Gateway{
		svr:              nil,
		stdMux:           http.NewServeMux(),
		gwMux:            runtime.NewServeMux(),
		logger:           logx.NewPkg(pkgName, "gateway"),
		withoutAPIPrefix: false,
	}

	for _, opt := range opts {
		opt(gw)
	}

	gw.svr = &http.Server{
		Addr:    addr,
		Handler: gw.stdMux,
	}

	if gw.withoutAPIPrefix {
		gw.stdMux.Handle("/", gw.gwMux)
	} else {
		gw.stdMux.Handle("/api/", http.StripPrefix("/api", gw.gwMux))
	}

	return gw
}

func (gw *Gateway) RegisterServiceHandlerServer(register func(*runtime.ServeMux)) *Gateway {
	register(gw.gwMux)
	return gw
}

func (gw *Gateway) Run(ctx context.Context) error {
	errCh := make(chan error)

	go func() {
		gw.logger.Info("gateway start listening and serving", slog.String("address", gw.svr.Addr))
		errCh <- gw.svr.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if err != http.ErrServerClosed {
			return fmt.Errorf("gateway serve: %w", err)
		}
	case <-ctx.Done():
	}

	gw.logger.Info("gateway shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := gw.svr.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("gateway shutdown: %w", err)
	}

	return nil
}
