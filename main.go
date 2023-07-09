package main

import (
	"os"

	apiv1 "github.com/choyri/gugu/api/v1"
	"github.com/choyri/gugu/internal/grpcx"
	"github.com/choyri/gugu/internal/logx"
	"github.com/choyri/gugu/internal/signals"
	"github.com/choyri/gugu/user"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/exp/slog"
)

func main() {
	_, err := logx.NewLogger(logx.Config{Level: "DEBUG", Development: true}, slog.String("app", "gugu"))
	if err != nil {
		slog.Error("new logger failed", logx.Error(err))
		os.Exit(1)
	}

	ctx := signals.SetupSignalHandler()

	userSvr := user.NewServer()

	err = grpcx.NewGateway("localhost:8080").
		RegisterServiceHandlerServer(func(mux *runtime.ServeMux) {
			_ = apiv1.RegisterUserServiceHandlerServer(nil, mux, userSvr)
		}).
		Run(ctx)
	if err != nil {
		slog.Error("gateway run failed", logx.Error(err))
		os.Exit(1)
	}
}
