package server

import (
	"html/template"

	"github.com/awfulbits/wikiofthings/database"
	"github.com/awfulbits/wikiofthings/pages"
)

type Title struct {
	MD           []byte
	Tags         []string
	Page         *pages.Page
	PageTemplate template.HTML
}

func loadTitle(title string, db *database.DB) (titlePage *Title, err error) {
	page, err := pages.New(title, db)
	if err != nil {
		return
	}
	err = page.Load(db)
	if err != nil {
		return
	}
	titlePage = &Title{}
	titlePage.PageTemplate = template.HTML(page.HTML)
	titlePage.Page = page
	return
}
