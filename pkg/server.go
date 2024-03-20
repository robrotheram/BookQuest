package BookQuest

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/rs/cors"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/zitadel-go/v3/pkg/authentication"
	openid "github.com/zitadel/zitadel-go/v3/pkg/authentication/oidc"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
)

type Server struct {
	Templates *template.Template
	auth      *authentication.Authenticator[*openid.UserInfoContext[*oidc.IDTokenClaims, *oidc.UserInfo]]
	oidc      *authentication.Interceptor[*openid.UserInfoContext[*oidc.IDTokenClaims, *oidc.UserInfo]]
	*http.ServeMux
	*Store
}

func (server *Server) Render(w io.Writer, name string, data any) error {
	return server.Templates.ExecuteTemplate(w, name, data)
}

func (server *Server) HandleSearch(w http.ResponseWriter, r *http.Request) {
	profile, _ := server.getUserProfile(r)

	query := r.URL.Query().Get("q")
	results, _ := server.Search(query)
	if len(results) == 1 {
		http.Redirect(w, r, results[0].Url, http.StatusTemporaryRedirect)
		return
	}
	search := r.URL.Query().Get("search")
	if len(search) > 0 {
		query = search
	}
	results, _ = server.Search(query)
	for i, result := range results {
		results[i].FavouritedByUser = result.isFavoirte(profile)
	}

	server.Render(w, "search", TemplateData{
		Form: FormData{
			Value: query,
		},
		Results: results,
	})
}

func (server *Server) HandleEdit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	profile, _ := server.getUserProfile(r)
	page, _ := server.Store.Get(id)
	if r.Method == "POST" {
		name := r.FormValue("Name")
		description := r.FormValue("Description")
		tag := r.FormValue("Tags")
		tags := []string{}
		for _, v := range strings.Split(tag, ",") {
			v = strings.ReplaceAll(v, " ", "")
			if len(v) > 1 {
				tags = append(tags, v)
			}
		}
		page.Name = name
		page.Description = description
		page.Tags = tags
		server.Save(page, profile)
		server.Render(w, "edit_success_alert", page)
		return
	}
	server.Render(w, "edit", page)
}

func (server *Server) HandleHomepage(w http.ResponseWriter, r *http.Request) {
	server.Render(w, "index", NewTemplateData())
}

func (server *Server) HandleCreate(w http.ResponseWriter, r *http.Request) {
	page := Page{}
	profile, _ := server.getUserProfile(r)

	if r.Method == "POST" {
		page = NewPage(r.FormValue("Url"), "")
		name := r.FormValue("Name")
		description := r.FormValue("Description")
		tag := r.FormValue("Tags")
		tags := []string{}
		for _, v := range strings.Split(tag, ",") {
			v = strings.ReplaceAll(v, " ", "")
			if len(v) > 1 {
				tags = append(tags, v)
			}
		}
		page.Name = name
		page.Description = description
		page.Tags = tags
		server.Save(page, profile)
		http.Redirect(w, r, "/edit/"+page.ID, http.StatusTemporaryRedirect)
		return
	}
	server.Render(w, "edit", page)
}

func (server *Server) HandleDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	server.Delete(id)
	server.Render(w, "delete_success_alert", nil)
}

func (server *Server) HandleFavoirte(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	profile, _ := server.getUserProfile(r)
	page, _ := server.Get(id)
	page.ToggleFavoirte(profile)
	server.update(page)
	if page.isFavoirte(profile) {
		page.FavouritedByUser = true
	}
	server.Render(w, "favourite_btn", page)
}

func (server *Server) HandleAdd(w http.ResponseWriter, r *http.Request) {
	profile, _ := server.getUserProfile(r)
	type req struct {
		Url  string `json:"url"`
		Html string `json:"html"`
	}
	var request req
	json.NewDecoder(r.Body).Decode(&request)
	page := NewPage(request.Url, profile.Id)
	page.parseHTML(request.Html)
	server.Save(page, profile)
	w.Write([]byte(page.ID))
}

func (server *Server) getUserProfile(r *http.Request) (Profile, error) {
	authCtx := server.oidc.Context(r.Context())
	if authCtx == nil {
		return Profile{}, fmt.Errorf("Profile Not found")
	}
	profile := Profile{
		Id:   authCtx.UserInfo.Subject,
		Name: authCtx.UserInfo.Name,
	}
	return profile, nil
}

func (server *Server) HandleProfile(w http.ResponseWriter, r *http.Request) {
	profile, _ := server.getUserProfile(r)
	server.UpdateStats(&profile)
	server.Render(w, "profile", profile)
}

func (server *Server) HandleAuthFunc(path string, handler func(w http.ResponseWriter, r *http.Request), check bool) {
	if check {
		server.Handle(path, server.oidc.CheckAuthentication()(http.HandlerFunc(handler)))
	} else {
		server.Handle(path, server.oidc.RequireAuthentication()(http.HandlerFunc(handler)))
	}
}

func (server *Server) Init(static embed.FS) {
	fSys, _ := fs.Sub(static, "static")
	server.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(fSys))))
	server.HandleAuthFunc("/delete/{id}", server.HandleDelete, true)
	server.HandleAuthFunc("/edit/{id}", server.HandleEdit, true)
	server.HandleAuthFunc("/fav/{id}", server.HandleFavoirte, true)
	server.HandleAuthFunc("/search", server.HandleSearch, false)
	server.HandleAuthFunc("/create", server.HandleCreate, false)
	server.HandleAuthFunc("/add", server.HandleAdd, true)
	server.HandleAuthFunc("/profile", server.HandleProfile, false)
	server.HandleAuthFunc("/", server.HandleHomepage, false)
	server.Handle("/auth/", server.auth)
}

func (server *Server) Start() {
	fmt.Println("Serving on port http://localhost:8090")
	handler := cors.Default().Handler(server)
	log.Fatal(http.ListenAndServe(":8090", handler))
}

func NewServer(store *Store, tmplateFS, staticFS embed.FS) *Server {
	ctx := context.Background()
	authN, err := authentication.New(ctx, zitadel.New(Configuration.OIDCServer), Configuration.SecretKey,
		openid.DefaultAuthentication(Configuration.ClientID, Configuration.RedirectURI, Configuration.SecretKey),
	)
	if err != nil {
		slog.Error("zitadel sdk could not initialize", "error", err)
		os.Exit(1)
	}

	router := http.NewServeMux()
	tmpl := template.Must(template.New("layout").ParseFS(tmplateFS, "views/*.html"))
	server := Server{
		ServeMux:  router,
		Templates: tmpl,
		Store:     store,
		auth:      authN,
		oidc:      authentication.Middleware(authN),
	}
	server.Init(staticFS)
	return &server
}
