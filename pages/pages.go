package pages

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

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

func New(escapedTitle string, db *database.DB) (page *Page, err error) {
	page = &Page{}
	page.Title, err = page.Unescape(escapedTitle)
	if err != nil {
		return
	}
	id, err := db.Read("pageidindex", []byte(page.Title))
	if err != nil {
		return
	}
	page.ID = db.BytesToId(id)
	return
}

func (p *Page) Create(db *database.DB) error {
	if p.Title == "" {
		return fmt.Errorf("cannot create page: empty title: \"%s\"", p.Title)
	}
	p.ID = db.CreateID("pages")
	if p.ID <= 0 {
		return fmt.Errorf("cannot create page: invalid id: \"%v\" - things are not well with the database if you are seeing this message, check page id and page title in pages bucket against key-value pairs in pageidindex bucket", p.ID)
	}
	pageBuf, err := json.Marshal(p)
	if err != nil {
		return err
	}
	err = db.Create("pages", db.IdToBytes(p.ID), pageBuf)
	if err != nil {
		return err
	}
	err = db.Create("pageidindex", []byte(p.Title), db.IdToBytes(p.ID))
	if err != nil {
		return err
	}
	return nil
}

func (p *Page) Load(db *database.DB) error {
	pageJSONBytes, err := db.Read("pages", db.IdToBytes(p.ID))
	if err != nil {
		return err
	}
	err = json.Unmarshal(pageJSONBytes, &p)
	if err != nil {
		return err
	}
	p.MdToHtml()
	return err
}

func (p *Page) MdToHtml() {
	if p.MD != nil {
		html := markdown.ToHTML(p.MD, mdParser(), nil)
		p.HTML = sanitize(html)
	}
}

func (p *Page) Escape() string {
	query := url.QueryEscape(p.Title)
	return strings.ReplaceAll(query, "+", "_")
}

func (p *Page) Unescape(escapedTitle string) (title string, err error) {
	query := strings.ReplaceAll(escapedTitle, "_", "+")
	title, err = url.QueryUnescape(query)
	return
}
