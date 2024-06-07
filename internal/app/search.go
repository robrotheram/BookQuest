package app

import (
	"BookQuest/internal/auth"
	"fmt"
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
	// results, _ := server.Search(query)
	// if len(results) == 1 {
	// 	http.Redirect(w, r, results[0].Url, http.StatusTemporaryRedirect)
	// 	return
	// }
	search := r.URL.Query().Get("search")
	if len(search) > 0 {
		query = search
	}

	links, _ := QueryLinks(app.db, query, user.Id)

	// results, _ = server.Search(query)
	// for i, result := range results {
	// 	results[i].FavouritedByUser = result.isFavoirte(profile)
	// }
	err := app.Render(w, "search", SearchPage{
		Form: FormData{
			Value: query,
		},
		Results: links,
	})

	fmt.Printf("%v \n", err)
}
