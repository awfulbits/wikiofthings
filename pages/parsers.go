package pages

import (
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
)

func mdParser() *parser.Parser {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)
	return parser
}

func sanitize(html []byte) []byte {
	return bluemonday.UGCPolicy().SanitizeBytes(html)
}
