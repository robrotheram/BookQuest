package icons

import (
	"crypto/md5"
	"image/color"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func shortText(text string) string {
	split := strings.Split(text, " ")
	if len(split) >= 2 {
		return string([]rune(split[0])[0]) + string([]rune(split[1])[0])
	}
	return string([]rune(text)[0]) + string([]rune(text)[1])
}

func colorFromText(text string) color.RGBA {
	hash := md5.Sum([]byte(text))
	return color.RGBA{
		R: hash[0],
		G: hash[1],
		B: hash[2],
		A: 255,
	}
}

func HandleIcon(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	w.Header().Set("Content-Type", "image/svg+xml")
	RenderSVG(strings.ToUpper(shortText(name)), colorFromText(name), w)
}
