package app

import (
	"BookQuest/internal/auth"
	"BookQuest/internal/models"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/zitadel/schema"
)

func (app *App) HandleLinkedUsed(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.GetUser(r)

	id := chi.URLParam(r, "id")
	link, _ := models.GetLink(app.db, id)
	models.UpdateUserLinkMeta(app.db, id, user.Id)
	http.Redirect(w, r, link.Url, http.StatusTemporaryRedirect)
}

type LinkEditData struct {
	Link          models.Link
	Teams         []models.Team
	SelectedTeams []models.Team
	Shareing      []models.ShareSettings
}

func (app *App) HandleLinkCreateDashboard(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.GetUser(r)

	teams, _ := models.GetTeamsByUser(app.db, user.Id)
	app.Render(w, "link_create_dashboard", LinkEditData{
		Link:  models.Link{},
		Teams: teams,
		Shareing: []models.ShareSettings{
			models.PUBLIC,
			models.PRIVATE,
			models.TEAM,
		},
	})
}

func (app *App) HandleLinkEditDashboard(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.GetUser(r)

	teams, _ := models.GetTeamsByUser(app.db, user.Id)
	id := chi.URLParam(r, "id")
	link, _ := models.GetLink(app.db, id)

	data := LinkEditData{
		Link:  link,
		Teams: teams,
		Shareing: []models.ShareSettings{
			models.PUBLIC,
			models.PRIVATE,
			models.TEAM,
		},
	}

	if link.Sharing == models.TEAM {
		linkTeams, _ := models.GetTeamsByLink(app.db, id)
		data.SelectedTeams = linkTeams
	}
	log.Println(app.Render(w, "link_edit_dashboard", data))
}

var decoder = schema.NewDecoder()

type LinkCreation struct {
	*models.Link
	Teams []string `schema:"team[]"`
}

func (app *App) HandleLinkCreation(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.GetUser(r)

	create := LinkCreation{
		Link: &models.Link{
			Id:      uuid.New(),
			Updated: time.Now(),
		},
	}
	if r.ParseForm() != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if decoder.Decode(&create, r.PostForm) != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if models.CreateLink(app.db, *create.Link, user.Id) != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if create.Sharing == models.TEAM {
		for _, team := range create.Teams {
			if models.AddLinkToTeam(app.db, create.Link.Id.String(), team) != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
	}
	w.Header().Set("HX-Redirect", "/dashboard/link/"+create.Id.String())
	w.WriteHeader(http.StatusOK)
}

func (app *App) HandleLinkEdit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	link, err := models.GetLink(app.db, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	create := LinkCreation{
		Link: &models.Link{
			Updated: time.Now(),
		},
	}
	if r.ParseForm() != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if decoder.Decode(&create, r.PostForm) != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//remove all team links
	if link.Sharing == models.TEAM {
		teams, err := models.GetTeamsByLink(app.db, link.Id.String())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		for _, team := range teams {
			models.RemoveLinksToTeam(app.db, link.Id.String(), team.Id.String())
		}
	}

	//Update Link
	link.Update(*create.Link)
	if create.Sharing == models.TEAM {
		for _, team := range create.Teams {
			if models.AddLinkToTeam(app.db, link.Id.String(), team) != nil {
				w.WriteHeader(http.StatusBadRequest)
				app.Render(w, "edit_error_alert", nil)
				return
			}
		}
	}

	if models.UpdateLink(app.db, link) != nil {
		w.WriteHeader(http.StatusBadRequest)
		app.Render(w, "edit_error_alert", nil)
		return
	}
	w.WriteHeader(http.StatusOK)
	app.Render(w, "edit_success_alert", nil)
}
