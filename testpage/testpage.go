package testpage

import (
	"fmt"
	"log"
	"os"

	"github.com/awfulbits/wikiofthings/database"
	"github.com/awfulbits/wikiofthings/pages"
)

func RunTest(db *database.DB) error {
	title := "Hello friend"
	page, err := getTestPage(title, db)
	if err != nil {
		log.Printf("error fetching test page id: %s - it is likely not present yet", err)
		if page.Title != title || page.ID <= 0 {
			err = createTestPage(page, db)
			if err != nil {
				return err
			}
			log.Print("Test page created - you should only see this message on first run or if there is no database on this system")
		}
	}
	page, err = getTestPage(title, db)
	if err != nil {
		return fmt.Errorf("error fetching test page id: %s", err)
	}
	log.Print("Test page id should be 1: " + fmt.Sprint(page.ID))
	err = page.Load(db)
	if err != nil {
		return fmt.Errorf("error loading test page: %s", err)
	}
	return nil
}

func getTestPage(title string, db *database.DB) (*pages.Page, error) {
	page := pages.Page{}
	page.Title = title
	escapedTitle := page.Escape()
	return pages.New(escapedTitle, db)
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
