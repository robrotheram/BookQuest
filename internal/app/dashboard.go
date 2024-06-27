package app

import (
	"BookQuest/internal/auth"
	"BookQuest/internal/models"
	"net/http"
	"sort"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type DashboadPage struct {
	Teams               []models.Team
	TopLinks            []models.LinkMeta
	TotalLinksShared    int
	TotalLinksClicked   int
	TotalLinksFavorited int
}

func (app *App) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.GetUser(r)

	teams, _ := models.GetTeamsByUser(app.db, user.Id)
	favs, _ := models.UserGetFavorites(app.db, user.Id)
	userLinks, _ := models.UserGetLinks(app.db, user.Id)
	meta, _ := models.GetUserLinksMeta(app.db, user.Id, -1)

	linksClicked := 0
	for _, m := range meta {
		linksClicked += m.Clicked
	}

	if len(meta) > 3 {
		meta = meta[:3]
	}

	app.RenderPage(w, "dashboard", user, DashboadPage{
		Teams:               teams,
		TopLinks:            meta,
		TotalLinksShared:    len(favs),
		TotalLinksFavorited: len(userLinks),
		TotalLinksClicked:   linksClicked,
	})
}

type LinkDashboardPage struct {
	Links []models.Link
}

func (app *App) HandleLinkDashboard(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.GetUser(r)

	links, _ := models.UserGetLinks(app.db, user.Id)
	page := LinkDashboardPage{
		Links: links,
	}
	app.RenderPage(w, "link_dashboard", user, page)
}

func (app *App) HandleLinkFilter(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.GetUser(r)

	links, _ := models.UserGetLinks(app.db, user.Id)
	query := r.FormValue("search")
	page := LinkDashboardPage{
		Links: models.FilterLinks(links, query),
	}
	app.RenderComponent(w, "link_table", page)
}

func (app *App) HandleLinkDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	models.DeleteLink(app.db, id)
	w.WriteHeader(http.StatusOK)
}

func (app *App) HandleTeamDashboard(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.GetUser(r)
	teams, _ := getTeams(app.db, user.Id)
	app.RenderPage(w, "teams_dashboard", user, teams)
}

func (app *App) HandleTeamFilter(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.GetUser(r)
	query := r.FormValue("search")
	teams, _ := getTeams(app.db, user.Id)
	app.RenderComponent(w, "team_cards", models.FilterTeams(teams, query))
}

func getTeams(db *bun.DB, id uuid.UUID) ([]models.Team, error) {
	teams, _ := models.GetTeamsForUser(db, id)
	sort.Slice(teams, func(i, j int) bool {
		return teams[i].Name < teams[j].Name
	})
	return teams, nil
}
