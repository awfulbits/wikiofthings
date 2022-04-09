package testpage

import (
	"fmt"
	"log"
	"os"

	"github.com/awfulbits/wikiofthings/database"
	"github.com/awfulbits/wikiofthings/pages"
)

func RunTest(db *database.DB) error {
	page := &pages.Page{}
	title := "Hello friend"
	page.Title = title
	escapedTitle := page.Escape()
	page, err := pages.New(escapedTitle, db)
	if err != nil {
		return fmt.Errorf("error fetching test page id: %s", err)
	}
	err = page.Load(db)
	if err != nil {
		return fmt.Errorf("error loading test page: %s", err)
	}
	if page.Title != title || page.ID <= 0 {
		err = createTestPage(page, db)
		if err != nil {
			return err
		}
		log.Print("Test page created, you should only see this message on first run or if there is no database on this system")
		return nil
	} else {
		log.Print("Test page id should be 1: " + fmt.Sprint(page.ID))
	}
	return nil
}

func createTestPage(page *pages.Page, db *database.DB) error {
	mdBytes, err := os.ReadFile("testpage/test.md")
	if err != nil {
		return fmt.Errorf("error reading test file: %s", err)
	}
	page.MD = mdBytes
	err = page.Create(db)
	if err != nil {
		return fmt.Errorf("error creating test page: %s", err)
	}
	return nil
}
