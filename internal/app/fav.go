package app

import (
	"BookQuest/internal/auth"
	"BookQuest/internal/models"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (app *App) HandleFavourite(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.GetUser(r)
	linkId := chi.URLParam(r, "id")

	if models.FavouriteExists(app.db, linkId, user.Id) {
		models.RemoveFavourite(app.db, linkId, user.Id)
		app.RenderComponent(w, "fav_btn", SearchResult{
			Link: models.Link{
				Id: uuid.MustParse(linkId),
			},
			Favorited: false,
		})
	} else {
		models.AddFavourite(app.db, linkId, user.Id)
		app.RenderComponent(w, "fav_btn", SearchResult{
			Link: models.Link{
				Id: uuid.MustParse(linkId),
			},
			Favorited: true,
		})
	}
}

func (app *App) HandleFavouriteFilter(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.GetUser(r)

	links, _ := models.UserGetFavorites(app.db, user.Id)
	query := r.FormValue("search")
	page := LinkDashboardPage{
		Links: models.FilterLinks(links, query),
	}
	app.RenderComponent(w, "link_fav_table", page)
}

func (app *App) HandleFavouriteDelete(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.GetUser(r)

	id := chi.URLParam(r, "id")
	models.RemoveFavourite(app.db, id, user.Id)
	w.WriteHeader(http.StatusOK)
}

func (app *App) HandleFavDashboard(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.GetUser(r)

	links, _ := models.UserGetFavorites(app.db, user.Id)
	page := LinkDashboardPage{
		Links: links,
	}
	app.RenderPage(w, "fav_dashboard", user, page)
}
