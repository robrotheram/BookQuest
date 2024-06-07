package app

import (
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

func (app *App) Render(w io.Writer, name string, data any) error {
	return app.template.Render(w, name, data)
}
