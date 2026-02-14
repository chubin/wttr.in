package location

import (
	"encoding/json"
	"log"
)

type Location struct {
	Name     string `db:"name,key" json:"name"`
	Lat      string `db:"lat" json:"latitude"`
	Lon      string `db:"lon" json:"longitude"`
	Timezone string `db:"timezone" json:"timezone"`
	Fullname string `db:"displayName" json:"address"`
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
