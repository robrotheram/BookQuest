package BookQuest

import (
	"BookQuest/internal/app"
	"BookQuest/internal/auth"
	"BookQuest/internal/icons"
	"BookQuest/internal/migrations"
	"BookQuest/internal/models"
	"BookQuest/internal/render"
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/zitadel/zitadel-go/v3/pkg/authentication"
	openid "github.com/zitadel/zitadel-go/v3/pkg/authentication/oidc"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
)

func CacheControlMiddleware(next http.Handler, maxAge time.Duration) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age="+strconv.Itoa(int(maxAge.Seconds())))
		next.ServeHTTP(w, r)
	})
}

func DB() *bun.DB {
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(Configuration.Database)))
	db := bun.NewDB(pgdb, pgdialect.New())

	if Configuration.Development {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}
	models.RegisterModels(db)
	return db
}

func NewServer(static, tmplateFS fs.FS) *chi.Mux {
	db := DB()
	migrations.Create(db)
	migrations.Migrate(db)

	authN, err := authentication.New(context.Background(), zitadel.New(Configuration.OIDCServer), Configuration.SecretKey,
		openid.DefaultAuthentication(Configuration.ClientID, Configuration.RedirectURI, Configuration.SecretKey),
	)
	if err != nil {
		log.Fatalf("Failed setting up Authentication: %v", err)
	}
	mw := auth.AuthMiddleware(db, authN)

	fSys, _ := fs.Sub(static, "static")
	cacheMaxAge := 24 * time.Hour // Example cache duration of 24 hours
	cacheControlFileServer := CacheControlMiddleware(http.FileServer(http.FS(fSys)), cacheMaxAge)

	r := chi.NewRouter()

	// Add middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))

	r.Use(cors.Handler(cors.Options{
		AllowOriginFunc: func(r *http.Request, origin string) bool { return true },
		AllowedMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:  []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	}))

	app := app.NewApp(db, render.NewRender(r, tmplateFS, Configuration.Development))

	// Protected route
	r.Group(func(r chi.Router) {
		r.Use(mw)
		r.Get("/", app.HandleHomepage)
		r.Get("/search", app.HandleSearch)

		r.Get("/link/{id}", app.HandleLinkedUsed)
		r.Post("/link/{id}/fav", app.HandleFavourite)
		r.Post("/link/{id}", app.HandleFavourite)

		r.Get("/dashboard", app.HandleDashboard)
		r.Get("/dashboard/favourites", app.HandleFavDashboard)
		r.Post("/dashboard/favourites", app.HandleFavouriteFilter)
		r.Delete("/dashboard/favourites/{id}", app.HandleFavouriteDelete)

		r.Get("/dashboard/links", app.HandleLinkDashboard)
		r.Post("/dashboard/links", app.HandleLinkFilter)

		r.Get("/dashboard/link", app.HandleLinkCreateDashboard)
		r.Post("/dashboard/link", app.HandleLinkCreation)
		r.Post("/dashboard/link/icon", app.HandleGetIcon)
		r.Get("/dashboard/link/{id}", app.HandleLinkEditDashboard)
		r.Post("/dashboard/link/{id}", app.HandleLinkEdit)
		r.Delete("/dashboard/link/{id}", app.HandleLinkDelete)

		r.Get("/dashboard/teams", app.HandleTeamDashboard)
		r.Post("/dashboard/teams", app.HandleTeamFilter)
		r.Get("/dashboard/team", app.HandleCreateTeamPage)
		r.Put("/dashboard/team", app.HandleTeamCreate)
		r.Get("/dashboard/team/{id}", app.HandleTeamPage)
		r.Post("/dashboard/team/{id}", app.HandleTeamEdit)
		r.Put("/dashboard/team/{id}/members", app.HandleTeamMemberAdd)
		r.Delete("/dashboard/team/{id}/members", app.HandleTeamMemberRemove)
		r.Post("/dashboard/team/{id}/members", app.HandleTeamMemberEdit)
	})
	// Public route
	r.Handle("/static/*", http.StripPrefix("/static/", cacheControlFileServer))
	r.Get("/icon/{name}", icons.HandleIcon)
	r.Get("/auth/callback", authN.Callback)
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})
	return r
}
