package server

import (
	"context"
	"crypto/tls"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/skdiver33/gophkeeper/internal/auth"
	"github.com/skdiver33/gophkeeper/internal/datamanager"
	"github.com/skdiver33/gophkeeper/internal/server/handler"
	"github.com/skdiver33/gophkeeper/internal/server/middleware"
	"github.com/skdiver33/gophkeeper/internal/usermanager"
	"github.com/skdiver33/gophkeeper/storage"
)

type AuthInterface interface {
	GetBaseToken() *jwtauth.JWTAuth
	CreateUserToken(userID int) (string, error)
	GetUserIDFromClaims(ctx context.Context) (int, error)
}

type StorageInterface interface {
	CloseStorage()
}

type KeeperServer struct {
	HandlersRouter http.Handler
	Auth           AuthInterface
	Storage        StorageInterface
}

func NewKeeperServer(config *KeeperServerConfig) (*KeeperServer, error) {

	newServer := &KeeperServer{}

	au := auth.NewAuth(config.SignKey)

	store, err := storage.NewSQLStorage(config.DBAddress)
	if err != nil {
		return nil, err
	}

	um := usermanager.NewUserManager(store, au)
	dm := datamanager.NewDataManager(store)

	newServer.Storage = store

	serverHandler := handler.NewServerHandler(um, dm)
	newRouter := chi.NewRouter()
	newRouter.Use(middleware.RequestLogger)
	newRouter.Use(middleware.GzipHandle)
	newRouter.Group(func(r chi.Router) {
		r.Route("/api/user/", func(r chi.Router) {
			r.Post("/register", serverHandler.UserRegisterHandler)
			r.Post("/login", serverHandler.UserLoginHandler)
		})
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(au.GetBaseToken()))
			r.Use(jwtauth.Authenticator(au.GetBaseToken()))
			r.Post("/data", serverHandler.LoadDataHandler)
			r.Get("/data", serverHandler.GetDataHandler)
			r.Get("/alldata", serverHandler.GetAllDataHandler)
			r.Delete("/data", serverHandler.DeleteDataHandler)
		})
	})
	newServer.HandlersRouter = newRouter
	return newServer, nil
}

func Run() {

	retCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer stop()

	config, err := NewKeeperServerConfig()
	if err != nil {
		slog.Error("read server config", "error", err.Error())
		return
	}

	ks, err := NewKeeperServer(config)
	if err != nil {
		slog.Error("create keeper server", "error", err.Error())
		return
	}

	server := &http.Server{
		Addr:      config.ListenAddr,
		TLSConfig: &tls.Config{},
		Handler:   ks.HandlersRouter,
	}

	go func() {

		slog.Info("Starting server", "address", config.ListenAddr)
		if err := server.ListenAndServeTLS(config.CertFile, config.KeyPath); err != nil && err != http.ErrServerClosed {
			slog.Error("start https server", "error", err.Error())
			stop()
		}
	}()
	<-retCtx.Done()
	stop()
	slog.Info("Server shutdowning....")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()
	server.Shutdown(shutdownCtx)
	ks.CloseStorage()
	slog.Info("Server shutdown.")
}

func (server *KeeperServer) CloseStorage() {
	server.Storage.CloseStorage()
}
