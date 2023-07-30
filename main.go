package main

import (
	"fmt"
	"os"

	apiv1 "github.com/choyri/gugu/api/v1"
	"github.com/choyri/gugu/config"
	"github.com/choyri/gugu/internal/databasex"
	"github.com/choyri/gugu/internal/grpcx"
	"github.com/choyri/gugu/internal/logx"
	"github.com/choyri/gugu/internal/signals"
	"github.com/choyri/gugu/user"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/exp/slog"
)

var (
	version = "dev"
)

func main() {
	conf, err := config.Init()
	checkErr(err)

	logger, err := logx.Init(conf.LogX, slog.String("app.version", version))
	checkErr(err)

	ctx := logx.WithContext(signals.SetupSignalHandler(), logger)

	db, err := databasex.New(ctx, conf.DatabaseX)
	checkErr(err)

	userSvr := user.NewServer(user.NewDBRepo(db))

	err = grpcx.NewGateway(conf.ListenAddr).
		RegisterServiceHandlerServer(func(mux *runtime.ServeMux) {
			_ = apiv1.RegisterUserServiceHandlerServer(nil, mux, userSvr)
		}).
		Run(ctx)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "GuGu:", err)
		os.Exit(1)
	}
}
