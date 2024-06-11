package app

import (
	"BookQuest/internal/models"
	"BookQuest/internal/render"
	"io"

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

func (app *App) RenderComponent(w io.Writer, name string, data any) error {
	return app.template.Render(w, name, data)
}

type PageData struct {
	User models.User
	Data any
}

func (app *App) RenderPage(w io.Writer, name string, user models.User, data any) error {
	return app.template.Render(w, name, PageData{
		User: user,
		Data: data,
	})
}
