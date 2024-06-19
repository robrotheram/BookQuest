package app

import (
	"BookQuest/internal/models"
	"BookQuest/internal/render"
	"net/http"

	"github.com/uptrace/bun"
)

type App struct {
	db       *bun.DB
	template *render.Render
}

func NewApp(db *bun.DB, template *render.Render) *App {
	return &App{
		db:       db,
		template: template,
	}
}

func (app *App) RenderComponent(w http.ResponseWriter, name string, data any) error {
	w.Header().Add("Content-Type", "text/html")
	w.Header().Add("Content-Encoding", "gzip")
	
	return app.template.Render(w, name, data)
}

type PageData struct {
	User       models.User
	Data       any
	LiveReload bool
}

func (app *App) RenderPage(w http.ResponseWriter, name string, user models.User, data any) error {
	w.Header().Add("Content-Type", "text/html")
	w.Header().Add("Content-Encoding", "gzip")

	return app.template.Render(w, name, PageData{
		User:       user,
		Data:       data,
		LiveReload: (app.template.Reload != nil),
	})
}
