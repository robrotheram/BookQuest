package app

import (
	"BookQuest/internal/auth"
	"BookQuest/internal/models"
	"fmt"
	"net/http"
)

type HomePage struct {
	Form     FormData
	MyLinks  []SearchResult
	TopLinks []SearchResult
}

func (app *App) HandleHomepage(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.GetUser(r)

	top := models.GetTopLinks(app.db)
	myLinks := models.GetUserTopLinks(app.db, user.Id)
	err := app.RenderPage(w, "index", user, HomePage{
		MyLinks:  ConvertLinkToSearchResult(myLinks),
		TopLinks: ConvertLinkToSearchResult(top),
	})
	fmt.Println(err)
}
