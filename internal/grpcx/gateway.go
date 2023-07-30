package grpcx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/textproto"
	"strconv"
	"time"

	"github.com/choyri/gugu/internal/logx"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Gateway struct {
	svr    *http.Server
	stdMux *http.ServeMux
	gwMux  *runtime.ServeMux
	logger *slog.Logger

	// for runtime.NewServeMux
	allowedHeaders     []string
	headerMatcher      runtime.HeaderMatcherFunc
	metadataAnnotators []func(context.Context, *http.Request) metadata.MD
	responseDisposers  []func(context.Context, http.ResponseWriter, proto.Message) error
	gwOpts             []runtime.ServeMuxOption

	middlewares      []GatewayMiddleware
	withoutAPIPrefix bool
}

type GatewayOption func(*Gateway)

type GatewayMiddleware func(http.Handler) http.Handler

const (
	MDHTTPStatusCode = "Http-Status-Code"
)

var (
	DefaultGatewayResponseDisposer = gatewayResponseDisposer

	DefaultGatewayMarshaler = &runtime.HTTPBodyMarshaler{
		Marshaler: &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		},
	}
)

func WithAllowedHeaders(headers []string) GatewayOption {
	return func(gw *Gateway) {
		gw.allowedHeaders = append(gw.allowedHeaders, headers...)
	}
}

func WithHeaderMatcher(matcher runtime.HeaderMatcherFunc) GatewayOption {
	return func(gw *Gateway) {
		gw.headerMatcher = matcher
	}
}

func WithMetadataAnnotator(annotator func(context.Context, *http.Request) metadata.MD) GatewayOption {
	return func(gw *Gateway) {
		gw.metadataAnnotators = append(gw.metadataAnnotators, annotator)
	}
}

func WithResponseDisposer(disposer func(context.Context, http.ResponseWriter, proto.Message) error) GatewayOption {
	return func(gw *Gateway) {
		gw.responseDisposers = append(gw.responseDisposers, disposer)
	}
}

func WithGatewayMuxOptions(opts ...runtime.ServeMuxOption) GatewayOption {
	return func(gw *Gateway) {
		gw.gwOpts = append(gw.gwOpts, opts...)
	}
}

func NewGateway(addr string, opts ...GatewayOption) *Gateway {
	stdMux := http.NewServeMux()

	gw := &Gateway{
		svr: &http.Server{
			Addr:    addr,
			Handler: stdMux,
		},
		stdMux: stdMux,
		logger: logx.NewPkg(pkgName, "gateway"),
	}

	for _, opt := range opts {
		opt(gw)
	}

	var (
		headerMatcher = gatewayHeaderMatcher(gw.allowedHeaders, gw.headerMatcher)
		gwOpts        = []runtime.ServeMuxOption{
			runtime.WithIncomingHeaderMatcher(headerMatcher),
			runtime.WithOutgoingHeaderMatcher(headerMatcher),
			runtime.WithMarshalerOption(runtime.MIMEWildcard, DefaultGatewayMarshaler),
			runtime.WithForwardResponseOption(DefaultGatewayResponseDisposer),
		}
	)

	for _, v := range gw.metadataAnnotators {
		gwOpts = append(gwOpts, runtime.WithMetadata(v))
	}
	for _, v := range gw.responseDisposers {
		gwOpts = append(gwOpts, runtime.WithForwardResponseOption(v))
	}

	gw.gwMux = runtime.NewServeMux(append(gwOpts, gw.gwOpts...)...)

	return gw
}

func (gw *Gateway) RegisterServiceHandlerServer(register func(*runtime.ServeMux)) *Gateway {
	register(gw.gwMux)
	return gw
}

func (gw *Gateway) Handle(pattern string, handler func(http.ResponseWriter, *http.Request)) *Gateway {
	gw.stdMux.HandleFunc(pattern, handler)
	return gw
}

func (gw *Gateway) WithMiddleware(mw GatewayMiddleware) *Gateway {
	gw.middlewares = append(gw.middlewares, mw)
	return gw
}

func (gw *Gateway) RemoveAPIPrefix() *Gateway {
	gw.withoutAPIPrefix = true
	return gw
}

func (gw *Gateway) Run(ctx context.Context) error {
	var (
		mainPattern string
		mainHandler http.Handler
	)

	if gw.withoutAPIPrefix {
		mainPattern = "/"
		mainHandler = gw.gwMux
	} else {
		mainPattern = "/api/"
		mainHandler = http.StripPrefix("/api", gw.gwMux)
	}

	for i := range gw.middlewares {
		mainHandler = gw.middlewares[len(gw.middlewares)-1-i](mainHandler)
	}

	gw.stdMux.Handle(mainPattern, mainHandler)

	errCh := make(chan error)

	go func() {
		gw.logger.Info("gateway start listening and serving", slog.String("address", gw.svr.Addr))
		errCh <- gw.svr.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
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

func GatewayLoggerMiddleware(logger *slog.Logger) GatewayMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			next.ServeHTTP(w, r.WithContext(logx.WithContext(r.Context(), logger)))
			logger.Info(
				"request finished",
				slog.String(pkgName+".request_method", r.Method),
				slog.String(pkgName+".request_path", r.URL.Path),
				slog.String(pkgName+".start_time", startTime.Format(time.RFC3339)),
				slog.Float64(pkgName+".cost_time_ms", float64(time.Since(startTime).Nanoseconds()/1000)/1000),
			)
		})
	}
}

func gatewayHeaderMatcher(allowedHeaders []string, customMatcher runtime.HeaderMatcherFunc) runtime.HeaderMatcherFunc {
	keyMap := make(map[string]struct{}, len(allowedHeaders))

	for _, v := range allowedHeaders {
		keyMap[textproto.CanonicalMIMEHeaderKey(v)] = struct{}{}
	}

	return func(key string) (string, bool) {
		if _, ok := keyMap[textproto.CanonicalMIMEHeaderKey(key)]; ok {
			return key, true
		}
		if customMatcher != nil {
			if v, ok := customMatcher(key); ok {
				return v, true
			}
		}
		return runtime.DefaultHeaderMatcher(key)
	}
}

func gatewayResponseDisposer(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
	md, _ := runtime.ServerMetadataFromContext(ctx)

	if sc := md.HeaderMD.Get(MDHTTPStatusCode); len(sc) > 0 && len(sc[0]) > 0 {
		code, err := strconv.Atoi(sc[0])
		if err != nil {
			return err
		}
		delete(w.Header(), runtime.MetadataHeaderPrefix+MDHTTPStatusCode)
		w.WriteHeader(code)
	}

	if _, ok := resp.(*emptypb.Empty); ok {
		w.WriteHeader(http.StatusNoContent)
	}

	return nil
}
