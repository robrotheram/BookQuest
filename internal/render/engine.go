package render

import (
	"BookQuest/internal/models"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/fsnotify/fsnotify"
	"github.com/go-chi/chi/v5"
)

type Render struct {
	watcher    *fsnotify.Watcher
	Templates  *template.Template
	templateFS fs.FS
	Reload     *LiveRealod
}

func hashFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return nil, err
	}
	return hasher.Sum(nil), nil
}

func (render *Render) applyFunctions() {
	funcMap := template.FuncMap{
		"ToUpper":        strings.ToUpper,
		"ToLower":        strings.ToLower,
		"ShareingFormat": func(s models.ShareSettings) string { return strings.ToTitle(string(s)) },
		"IncludesTeam":   isSelectedTeam,
	}
	render.Templates = render.Templates.Funcs(funcMap)
}
func (render *Render) BuildTemplates() (string, error) {
	render.Templates = template.New("layout")
	render.applyFunctions()
	files := []string{}

	err := fs.WalkDir(render.templateFS, "views", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".html") {
			files = append(files, path)
			tmplData, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			_, err = render.Templates.Parse(string(tmplData))
			if err != nil {
				return err
			}
		}
		return nil
	})
	// Ensure consistent order
	sort.Strings(files)
	hasher := sha256.New()
	for _, file := range files {
		fileHash, err := hashFile(file)
		if err != nil {
			break
		}
		hasher.Write(fileHash)
	}
	return hex.EncodeToString(hasher.Sum(nil)), err
}

func (render *Render) watch() {
	fs.WalkDir(render.templateFS, "views", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".html") {
			render.watcher.Add(path)
		}
		return nil
	})

	go func() {
		for {
			select {
			case event := <-render.watcher.Events:
				slog.Info("Change: " + event.Name)
				if hash, err := render.BuildTemplates(); err == nil {
					render.Reload.Broadcast(hash)
				}
			case err := <-render.watcher.Errors:
				slog.Warn("ERROR", err)
			}
		}
	}()
}

func (render *Render) Render(w io.Writer, name string, data any) error {
	return render.Templates.ExecuteTemplate(w, name, data)
}

func NewRender(mux *chi.Mux, templateFS fs.FS, liveReload bool) *Render {
	render := &Render{
		templateFS: templateFS,
	}
	if hash, err := render.BuildTemplates(); err == nil && liveReload {
		watcher, _ := fsnotify.NewWatcher()
		render.watcher = watcher
		render.Reload = NewLiveReloadSever(hash)
		mux.HandleFunc("/live-reload", render.Reload.CreateLiveReloadHandler())
		render.watch()
	}
	return render
}
