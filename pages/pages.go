package pages

import (
	"encoding/json"

	"github.com/awfulbits/wikiofthings/database"
	"github.com/gomarkdown/markdown"
)

type Page struct {
	ID    int      `json:"id"`
	Title string   `json:"title"`
	MD    []byte   `json:"md"`
	Tags  []string `json:"tags"`
	JSON  []byte
	HTML  []byte
}

func New(title string) *Page {
	return &Page{Title: title}
}

func (p *Page) Load(db *database.DB) error {
	pageJSONBytes := db.Get("pages", p.Title)
	err := json.Unmarshal(pageJSONBytes, &p)
	if err != nil {
		return err
	}
	p.MdToHtml()
	return err
}

func (p *Page) MdToHtml() {
	p.HTML = sanitize(markdown.ToHTML(p.MD, mdParser(), mdRenderer()))
}
