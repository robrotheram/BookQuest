package app

import (
	"BookQuest/internal/auth"
	"BookQuest/internal/models"
	"io"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type TeamEditData struct {
	Team        models.Team
	Permissions []models.TeamPermission
	Memebers    []models.UserToTeam
}

func (app *App) HandleTeamPage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, _ := auth.GetUser(r)

	team, _ := models.GetTeam(app.db, id)
	users, _ := models.GetTeamPermissions(app.db, id)
	app.RenderPage(w, "edit_team_dashboard", user, TeamEditData{
		Team:     team,
		Memebers: users,
		Permissions: []models.TeamPermission{
			models.OWNER,
			models.MEMBER,
		},
	})
}

func (app *App) HandleTeamMemberAdd(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	team, _ := models.GetTeam(app.db, id)

	if r.ParseForm() != nil {
		w.WriteHeader(http.StatusBadRequest)
		app.RenderComponent(w, "edit_error_alert", nil)
		return
	}
	email := r.PostForm.Get("email")
	permission := models.TeamPermission(r.PostForm.Get("permission"))
	//TODO: validate email and permission

	user, err := models.GetUserByEmail(app.db, email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		app.RenderComponent(w, "edit_error_alert", nil)
		return
	}

	if models.AddUserToTeam(app.db, user.Id, team.Id.String(), permission) != nil {
		w.WriteHeader(http.StatusBadRequest)
		app.RenderComponent(w, "edit_error_alert", nil)
		return
	}

	users, _ := models.GetTeamPermissions(app.db, team.Id.String())
	app.RenderComponent(w, "team_members", TeamEditData{
		Memebers: users,
		Permissions: []models.TeamPermission{
			models.OWNER,
			models.MEMBER,
		},
	})
}

func (app *App) HandleTeamMemberRemove(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	team, _ := models.GetTeam(app.db, id)

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	form, err := url.ParseQuery(string(body))
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	email := form.Get("email")
	user, err := models.GetUserByEmail(app.db, email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		app.RenderComponent(w, "edit_error_alert", nil)
		return
	}

	if models.RemoveUserToTeam(app.db, user.Id.String(), team.Id.String()) != nil {
		w.WriteHeader(http.StatusBadRequest)
		app.RenderComponent(w, "edit_error_alert", nil)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *App) HandleTeamMemberEdit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	team, _ := models.GetTeam(app.db, id)

	if r.ParseForm() != nil {
		w.WriteHeader(http.StatusBadRequest)
		app.RenderComponent(w, "edit_error_alert", nil)
		return
	}

	email := r.PostForm.Get("email")
	permission := models.TeamPermission(r.PostForm.Get("permission"))

	user, err := models.GetUserByEmail(app.db, email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		app.RenderComponent(w, "edit_error_alert", nil)
		return
	}

	if models.ModifyUserToTeam(app.db, user.Id.String(), team.Id.String(), permission) != nil {
		w.WriteHeader(http.StatusBadRequest)
		app.RenderComponent(w, "edit_error_alert", nil)
		return
	}

	users, _ := models.GetTeamPermissions(app.db, team.Id.String())
	app.RenderComponent(w, "team_members", TeamEditData{
		Team:     team,
		Memebers: users,
		Permissions: []models.TeamPermission{
			models.OWNER,
			models.MEMBER,
		},
	})
}

func (app *App) HandleCreateTeamPage(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.GetUser(r)
	app.RenderPage(w, "create_team_dashboard", user, nil)
}

func (app *App) HandleTeamCreate(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.GetUser(r)

	var team = models.Team{
		Id: uuid.New(),
	}
	if r.ParseForm() != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if decoder.Decode(&team, r.PostForm) != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if models.CreateTeam(app.db, team) != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if models.AddUserToTeam(app.db, user.Id, team.Id.String(), models.OWNER) != nil {
		w.WriteHeader(http.StatusBadRequest)
		app.RenderComponent(w, "edit_error_alert", nil)
		return
	}
	w.Header().Set("HX-Redirect", "/dashboard/team/"+team.Id.String())
	w.WriteHeader(http.StatusOK)
}

func (app *App) HandleTeamEdit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	team, err := models.GetTeam(app.db, id)
	var editTeam = models.Team{}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		app.RenderComponent(w, "edit_error_alert", nil)
		return
	}
	if r.ParseForm() != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if decoder.Decode(&editTeam, r.PostForm) != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	team.Update(editTeam)
	if models.UpdateTeam(app.db, team) != nil {
		w.WriteHeader(http.StatusBadRequest)
		app.RenderComponent(w, "edit_error_alert", nil)
		return
	}
	w.WriteHeader(http.StatusOK)
	app.RenderComponent(w, "edit_success_alert", nil)
}
