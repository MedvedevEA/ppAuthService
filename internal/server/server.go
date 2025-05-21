package server

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"ppAuthService/internal/config"
	srvErr "ppAuthService/internal/server/err"
	"ppAuthService/internal/service"

	proto "github.com/MedvedevEA/ppProtos/gen/auth"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	lg         *slog.Logger
	grpcServer *grpc.Server
	cfg        *config.Server
}

func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
func MustNew(service *service.Service, lg *slog.Logger, cfg *config.Server) *Server {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(logging.FinishCall),
	}
	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			lg.Error(p.(string), slog.String("owner", "server"))
			return status.Error(codes.Internal, srvErr.ErrInternalServerError.Error())
		}),
	}
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(InterceptorLogger(lg), loggingOpts...),
	))
	proto.RegisterAuthServiceServer(grpcServer, service)

	return &Server{
		lg:         lg,
		grpcServer: grpcServer,
		cfg:        cfg,
	}
}

func (s *Server) Start() {
	chErr := make(chan error, 1)
	defer close(chErr)

	go func() {
		s.lg.Info("server is started", slog.String("owner", "server"), slog.String("bindAddress", s.cfg.BindAddr))
		listener, err := net.Listen("tcp", s.cfg.BindAddr)
		if err != nil {
			chErr <- err
			return
		}
		chErr <- s.grpcServer.Serve(listener)
	}()
	go func() {
		chQuit := make(chan os.Signal, 1)
		signal.Notify(chQuit, syscall.SIGINT, syscall.SIGTERM)
		<-chQuit
		s.grpcServer.Stop()
		chErr <- nil
	}()
	if err := <-chErr; err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "server"))
		return
	}
	s.lg.Info("server is stoped", slog.String("owner", "server"))

}
