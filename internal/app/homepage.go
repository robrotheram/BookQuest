package app

import (
	"BookQuest/internal/auth"
	"BookQuest/internal/models"
	"net/http"
)

type HomePage struct {
	MyLinks  []SearchResult
	TopLinks []SearchResult
}

func (app *App) HandleHomepage(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.GetUser(r)

	top := models.GetTopLinks(app.db)
	myLinks := models.GetUserTopLinks(app.db, user.Id)
	app.Render(w, "index", HomePage{
		MyLinks:  ConvertLinkToSearchResult(myLinks),
		TopLinks: ConvertLinkToSearchResult(top),
	})
}
