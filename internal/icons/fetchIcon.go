package icons

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"golang.org/x/net/html"
)

type icon struct {
	Href string
	Rel  string
	Size string
}

func GetIcon(url string) (string, error) {
	body, err := fetchHTML(url)
	if err != nil {
		return "", err
	}
	icons, err := parseIcons(body)
	if err != nil {
		return "", err
	}
	bestIcon := filterBestIcon(icons)
	if bestIcon == nil {
		return "", err
	}
	return bestIcon.Href, nil
}

func fetchHTML(url string) (io.Reader, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set headers to mimic a web browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Connection", "keep-alive")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch URL: %s, status code: %d", url, resp.StatusCode)
	}

	return resp.Body, nil
}

// parseIcons parses the HTML content and extracts potential icon links
func parseIcons(body io.Reader) (map[IconType][]icon, error) {
	var icons = make(map[IconType][]icon)

	doc, err := html.Parse(body)
	if err != nil {
		return nil, err
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "link" {
			var href, rel, size string
			for _, attr := range n.Attr {
				switch strings.ToLower(attr.Key) {
				case "href":
					href = attr.Val
				case "rel":
					rel = attr.Val
				case "sizes":
					size = attr.Val
				}

			}
			if href != "" {
				switch IconType(rel) {
				case appleIcon:
					icons[appleIcon] = append(icons[appleIcon], icon{Href: href, Rel: rel, Size: size})
				case shortcutIcon:
					icons[shortcutIcon] = append(icons[shortcutIcon], icon{Href: href, Rel: rel})
				case defaultIcon:
					icons[defaultIcon] = append(icons[defaultIcon], icon{Href: href, Rel: rel, Size: size})
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return icons, nil
}

type IconType string

const (
	appleIcon    = IconType("apple-touch-icon")
	shortcutIcon = IconType("shortcut icon")
	defaultIcon  = IconType("icon")
)

// filterBestIcon filters and returns the best icon based on common practices
func filterBestIcon(icons map[IconType][]icon) *icon {

	if icon, ok := icons[appleIcon]; ok {
		if len(icon) > 0 {
			return &icon[0]
		}
	}

	if icon, ok := icons[defaultIcon]; ok {
		sort.Slice(icon, func(i, j int) bool {
			return icon[i].Size < icon[j].Size
		})
		return &icon[0]
	}

	if icon, ok := icons[shortcutIcon]; ok {
		if len(icon) > 0 {
			return &icon[0]
		}
	}
	return nil
}
