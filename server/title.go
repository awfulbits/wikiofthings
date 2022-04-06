package server

import (
	"html/template"

	"github.com/awfulbits/wikiofthings/database"
	"github.com/awfulbits/wikiofthings/pages"
)

type Title struct {
	MD   []byte
	Tags []string
	Page template.HTML
}

func loadTitle(title string, db *database.DB) (*Title, error) {
	page := pages.New(title)
	err := page.Load(db)
	if err != nil {
		return nil, err
	}
	htmlPage := template.HTML(page.HTML)
	return &Title{Page: htmlPage}, nil
}
