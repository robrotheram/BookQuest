package app

import (
	"BookQuest/internal/auth"
	"BookQuest/internal/models"
	"net/http"
)

type FormData struct {
	Value string
}

type SearchPage struct {
	Form    FormData
	Results []SearchResult
}

func (app *App) HandleSearch(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.GetUser(r)
	query := r.URL.Query().Get("q")
	search := r.URL.Query().Get("search")
	if len(search) > 0 {
		query = search
	}
	links, _ := QueryLinks(app.db, query, user.Id.String())
	if len(links) == 1 {
		go models.UpdateUserLinkMeta(app.db, links[0].Id.String(), user.Id)
		http.Redirect(w, r, links[0].Url, http.StatusTemporaryRedirect)
		return
	}
	app.RenderPage(w, "search", user, SearchPage{
		Form: FormData{
			Value: query,
		},
		Results: links,
	})
}
