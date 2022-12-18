package location

import (
	"encoding/json"
	"log"
)

type Location struct {
	Name     string `db:"name,key"`
	Lat      string `db:"lat"`
	Lon      string `db:"lon"`
	Timezone string `db:"timezone"`
	//nolint:tagliatelle
	Fullname string `db:"displayName" json:"display_name"`
}

// String returns string representation of location.
func (l *Location) String() string {
	bytes, err := json.Marshal(l)
	if err != nil {
		// should never happen
		log.Fatalln(err)
	}

	return string(bytes)
}
