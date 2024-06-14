package app

import (
	"BookQuest/internal/auth"
	"BookQuest/internal/icons"
	"BookQuest/internal/models"
	"encoding/base64"
	"fmt"
	"log/slog"
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

	data, _ := base64.StdEncoding.DecodeString(r.URL.Query().Get("data"))
	process(string(data))

	link := models.Link{
		Title:       r.URL.Query().Get("title"),
		Icon:        r.URL.Query().Get("icon"),
		Description: r.URL.Query().Get("description"),
		Url:         r.URL.Query().Get("url"),
	}
	app.RenderPage(w, "link_create_dashboard", user, LinkEditData{
		Link:  link,
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
	app.RenderPage(w, "link_edit_dashboard", user, data)
}

var decoder = schema.NewDecoder()

type LinkCreation struct {
	*models.Link
	Teams []string `schema:"team[]"`
}

func (app *App) HandleGetIcon(w http.ResponseWriter, r *http.Request) {
	if r.ParseForm() != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	url := r.PostForm["url"][0]
	icon, err := icons.GetIcon(url)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	template := `
	<input id="icon" autocomplete="icon" required=""
                class="block w-full appearance-none rounded-md border border-gray-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 "
                type="text" name="icon" value="%s" />
	`
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(template, icon)))
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
		slog.Warn("unable to  form data for user: %s", "user", user.Username)
		return
	}
	if decoder.Decode(&create, r.PostForm) != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Warn("unable to decode form data", "user", user.Username)
		return
	}

	if err := models.CreateLink(app.db, *create.Link, user.Id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Warn(err.Error(), "user", user.Username)
		return
	}
	if create.Sharing == models.TEAM {
		for _, team := range create.Teams {
			if err := models.AddLinkToTeam(app.db, create.Link.Id.String(), team); err != nil {
				slog.Warn(err.Error(), "user", user.Username)
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
	user, _ := auth.GetUser(r)
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
		slog.Warn("unable to decode form data", "user", user.Username)
		return
	}
	if decoder.Decode(&create, r.PostForm) != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Warn("unable to decode form data", "user", user.Username)
		return
	}

	//remove all team links
	if link.Sharing == models.TEAM {
		teams, err := models.GetTeamsByLink(app.db, link.Id.String())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			slog.Warn("unable to remove links", "user", user.Username)
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
			if err := models.AddLinkToTeam(app.db, link.Id.String(), team); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				slog.Warn(err.Error(), "user", user.Username)
				app.RenderComponent(w, "edit_error_alert", nil)
				return
			}
		}
	}

	if err := models.UpdateLink(app.db, link); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Warn(err.Error(), "user", user.Username)
		app.RenderComponent(w, "edit_error_alert", nil)
		return
	}
	w.WriteHeader(http.StatusOK)
	app.RenderComponent(w, "edit_success_alert", nil)
}
