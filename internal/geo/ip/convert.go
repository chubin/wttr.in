package ip

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/samonzeweb/godb"
	"github.com/samonzeweb/godb/adapters/sqlite"
)

func (c *Cache) ConvertCache() error {
	dbfile := c.config.Geo.IPCacheDB

	err := removeDBIfExists(dbfile)
	if err != nil {
		return err
	}

	db, err := godb.Open(sqlite.Adapter, dbfile)
	if err != nil {
		return err
	}

	err = createTable(db, "Location")
	if err != nil {
		return err
	}

	log.Println("listing cache entries...")
	files, err := filepath.Glob(filepath.Join(c.config.Geo.IPCache, "*"))
	if err != nil {
		return err
	}

	log.Printf("going to convert %d entries\n", len(files))

	block := []Location{}
	for i, file := range files {
		ip := filepath.Base(file)
		loc, err := c.Read(ip)
		if err != nil {
			log.Println("invalid entry for", ip)
			continue
		}

		block = append(block, *loc)
		if i%1000 != 0 || i == 0 {
			continue
		}

		err = db.BulkInsert(&block).Do()
		if err != nil {
			return err
		}
		block = []Location{}
		log.Println("converted", i+1, "entries")
	}

	// inserting the rest.
	err = db.BulkInsert(&block).Do()
	if err != nil {
		return err
	}

	log.Println("converted", len(files), "entries")

	return nil
}

func createTable(db *godb.DB, tableName string) error {
	createTable := fmt.Sprintf(
		`create table %s (
	    ip             text not null primary key,
        countryCode    text not null,
        country        text not null,
        region         text not null,
        city           text not null,
        latitude       text not null,
        longitude      text not null);
	`, tableName)

	_, err := db.CurrentDB().Exec(createTable)
	return err
}

func removeDBIfExists(filename string) error {
	_, err := os.Stat(filename)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		// no db file
		return nil
	}

	return os.Remove(filename)
}
