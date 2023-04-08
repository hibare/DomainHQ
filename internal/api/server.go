package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hibare/DomainHQ/internal/api/handlers"
	"github.com/hibare/DomainHQ/internal/config"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Good to see you")
}

func Serve() {
	wait := time.Second * 15
	addr := fmt.Sprintf("%s:%d", config.Current.Server.ListenAddr, config.Current.Server.ListenPort)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.StripSlashes)
	r.Use(middleware.CleanPath)

	r.Get("/", home)
	r.Get("/ping", handlers.HealthCheck)
	r.With(httpin.NewInput(handlers.WebFingerParams{})).Get("/.well-known/webfinger", handlers.WebFinger)

	srv := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  time.Second * 60,
	}

	log.Printf("Listening for address %s on port %d\n", config.Current.Server.ListenAddr, config.Current.Server.ListenPort)

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	srv.Shutdown(ctx)
}

func init() {
	// Register a directive named "path" to retrieve values from `chi.URLParam`,
	httpin.UseGochiURLParam("path", chi.URLParam)
}
