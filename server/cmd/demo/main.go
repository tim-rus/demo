package main

import (
	"context"
	"demo/interal/proto/apigen/apigenconnect"
	"demo/interal/proto/tggen/tggenconnect"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SladkyCitron/slogcolor"
)

//

func main() {

	// config

	opts := slogcolor.DefaultOptions

	slog.SetDefault(slog.New(
		slogcolor.NewHandler(os.Stderr, opts),
	))

	ctx, sysCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer sysCancel()

	if err := run(ctx); err != nil {
		slog.Error("app error", "err", err)
		os.Exit(1)
	}

	slog.Info("app exited")
}

//

func run(ctx context.Context) error {

	// services

	tgSvc := NewTGServce()
	apiSvc := NewAPIServce()

	// router

	mux := http.NewServeMux()

	tgPath, tgHandler := tggenconnect.NewPingServiceHandler(tgSvc)
	mux.Handle(tgPath, tgHandler)

	apiPath, apiHandler := apigenconnect.NewPingServiceHandler(apiSvc)
	mux.Handle(apiPath, apiHandler)

	// server

	handlerWithCORS := corsAllowAll(mux)

	addr := net.JoinHostPort(os.Getenv("HOST"), os.Getenv("PORT"))

	protocols := &http.Protocols{}
	protocols.SetUnencryptedHTTP2(true)
	protocols.SetHTTP1(true)

	rpcServer := &http.Server{
		Addr:              addr,
		Handler:           handlerWithCORS,
		Protocols:         protocols,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       0,
		WriteTimeout:      0,
		IdleTimeout:       0,
	}

	// run

	rpcErrCh := make(chan error, 1)

	go func() {
		slog.Info("rpc server started", "at", rpcServer.Addr)
		if err := rpcServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			rpcErrCh <- err
		}
	}()

	// shutdown

	var unexpectedErr error

	select {
	case unexpectedErr = <-rpcErrCh:
		slog.Error("rpc server failed", "error", unexpectedErr)
	case <-ctx.Done():
		slog.Info("interrupt signal received")
	}

	slog.Info("shutting down")

	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		time.Duration(10)*time.Second, // TODO: configure
	)
	defer shutdownCancel()

	if err := rpcServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("failed to gracefully shutdown rpc server", "error", err)
	}

	return unexpectedErr
}
