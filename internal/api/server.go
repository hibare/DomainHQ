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
	"github.com/hibare/DomainHQ/internal/api/handler"
	"github.com/hibare/DomainHQ/internal/config"
	"github.com/hibare/DomainHQ/internal/models"
	commonHandler "github.com/hibare/GoCommon/v2/pkg/http/handler"
	commonMiddleware "github.com/hibare/GoCommon/v2/pkg/http/middleware"
	"gorm.io/gorm"
)

type App struct {
	Router *chi.Mux
	DB     *gorm.DB
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Good to see you")
}

func (a *App) Init() {
	db, err := models.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	a.DB = db

	a.Router = chi.NewRouter()
	a.Router.Use(middleware.RequestID)
	a.Router.Use(middleware.RealIP)
	a.Router.Use(middleware.Logger)
	a.Router.Use(middleware.Recoverer)
	a.Router.Use(middleware.Timeout(60 * time.Second))
	a.Router.Use(middleware.StripSlashes)
	a.Router.Use(middleware.CleanPath)

	a.Router.Get("/", home)
	a.Router.Get("/ping", commonHandler.HealthCheck)
	a.Router.With(httpin.NewInput(handler.WebFingerParams{})).Get("/.well-known/webfinger", handler.WebFinger)
	a.Router.Route("/pks", func(r chi.Router) {
		r.With(httpin.NewInput(handler.GPGLookupParams{})).Get("/lookup", func(w http.ResponseWriter, r *http.Request) {
			handler.GPGPubKeyLookup(a.DB, w, r)
		})
		r.Group(func(r chi.Router) {
			r.Use(func(h http.Handler) http.Handler {
				return commonMiddleware.TokenAuth(h, config.Current.APIConfig.APIKeys)
			})
			r.With(httpin.NewInput(handler.GPGKeyAddParams{})).Post("/add", func(w http.ResponseWriter, r *http.Request) {
				handler.GPGPubKeyAdd(a.DB, w, r)
			})
		})
	})
}

func (a *App) Serve() {
	wait := time.Second * 15
	addr := fmt.Sprintf("%s:%d", config.Current.Server.ListenAddr, config.Current.Server.ListenPort)

	srv := &http.Server{
		Handler:      a.Router,
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
