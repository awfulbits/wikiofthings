package pages

import (
	"github.com/gomarkdown/markdown/html"
)

func mdRenderer() *html.Renderer {
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)
	renderer = nil
	return renderer
}
