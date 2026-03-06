package location

import (
	"fmt"
	"os"

	"github.com/samonzeweb/godb"
	"github.com/samonzeweb/godb/adapters/sqlite"
)

// ConvertCache converts files-based cache into the DB-based cache.
// If reset is true, the DB cache is created from scratch.
//
//nolint:funlen,cyclop
func (c *Cache) InitDB(reset bool) error {
	var (
		dbfile    = c.config.LocationCacheDB
		tableName = "Location"
	)

	if !reset && fileExists(dbfile) {
		return nil
	}

	if reset {
		err := removeDBIfExists(dbfile)
		if err != nil {
			return err
		}
	}

	db, err := godb.Open(sqlite.Adapter, dbfile)
	if err != nil {
		return err
	}

	err = createTable(db, tableName)
	if err != nil {
		return err
	}

	return nil
}

func createTable(db *godb.DB, tableName string) error {
	createTable := fmt.Sprintf(
		`create table %s (
	    name           text not null primary key,
        displayName    text not null,
        lat            text not null,
        lon            text not null,
		timezone       text not null);
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

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return true
	}
	return true
}
