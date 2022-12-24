// This code represents wttr.in view v1.
// It is based on wego (github.com/schachmat/wego) from which it diverged back in 2016.

package v1

import (
	_ "crypto/sha512"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"regexp"
	"strings"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-runewidth"
)

type Configuration struct {
	APIKey       string
	City         string
	Numdays      int
	Imperial     bool
	WindUnit     bool
	Inverse      bool
	Lang         string
	Narrow       bool
	LocationName string
	WindMS       bool
	RightToLeft  bool
}

type global struct {
	ansiEsc    *regexp.Regexp
	config     Configuration
	configpath string
	debug      bool
}

const (
	wuri      = "http://127.0.0.1:5001/premium/v1/weather.ashx?"
	suri      = "http://127.0.0.1:5001/premium/v1/search.ashx?"
	slotcount = 4
)

func (g *global) configload() error {
	b, err := ioutil.ReadFile(g.configpath)
	if err == nil {
		return json.Unmarshal(b, &g.config)
	}

	return err
}

func (g *global) configsave() error {
	j, err := json.MarshalIndent(g.config, "", "\t")
	if err == nil {
		return ioutil.WriteFile(g.configpath, j, 0o600)
	}

	return err
}

func (g *global) init() {
	flag.IntVar(&g.config.Numdays, "days", 3, "Number of days of weather forecast to be displayed")
	flag.StringVar(&g.config.Lang, "lang", "en", "Language of the report")
	flag.StringVar(&g.config.City, "city", "New York", "City to be queried")
	flag.BoolVar(&g.debug, "debug", false, "Print out raw json response for debugging purposes")
	flag.BoolVar(&g.config.Imperial, "imperial", false, "Use imperial units")
	flag.BoolVar(&g.config.Inverse, "inverse", false, "Use inverted colors")
	flag.BoolVar(&g.config.Narrow, "narrow", false, "Narrow output (two columns)")
	flag.StringVar(&g.config.LocationName, "location_name", "", "Location name (used in the caption)")
	flag.BoolVar(&g.config.WindMS, "wind_in_ms", false, "Show wind speed in m/s")
	flag.BoolVar(&g.config.RightToLeft, "right_to_left", false, "Right to left script")
	g.configpath = os.Getenv("WEGORC")
	if g.configpath == "" {
		usr, err := user.Current()
		if err != nil {
			log.Fatalf("%v\nYou can set the environment variable WEGORC to point to your config file as a workaround.", err)
		}
		g.configpath = path.Join(usr.HomeDir, ".wegorc")
	}
	g.config.APIKey = ""
	g.config.Imperial = false
	g.config.Lang = "en"
	err := g.configload()
	if _, ok := err.(*os.PathError); ok {
		log.Printf("No config file found. Creating %s ...", g.configpath)
		if err2 := g.configsave(); err2 != nil {
			log.Fatal(err2)
		}
	} else if err != nil {
		log.Fatalf("could not parse %v: %v", g.configpath, err)
	}

	g.ansiEsc = regexp.MustCompile("\033.*?m")
}

func Cmd() error {
	g := global{}
	g.init()

	flag.Parse()

	r, err := g.getDataFromAPI()
	if err != nil {
		return err
	}

	if r.Data.Req == nil || len(r.Data.Req) < 1 {
		if r.Data.Err != nil && len(r.Data.Err) >= 1 {
			log.Fatal(r.Data.Err[0].Msg)
		}
		log.Fatal("Malformed response.")
	}
	locationName := r.Data.Req[0].Query
	if g.config.LocationName != "" {
		locationName = g.config.LocationName
	}
	if g.config.Lang == "he" || g.config.Lang == "ar" || g.config.Lang == "fa" {
		g.config.RightToLeft = true
	}
	if caption, ok := localizedCaption()[g.config.Lang]; !ok {
		fmt.Printf("Weather report: %s\n\n", locationName)
	} else {
		if g.config.RightToLeft {
			caption = locationName + " " + caption
			space := strings.Repeat(" ", 125-runewidth.StringWidth(caption))
			fmt.Printf("%s%s\n\n", space, caption)
		} else {
			fmt.Printf("%s %s\n\n", caption, locationName)
		}
	}
	stdout := colorable.NewColorableStdout()

	if r.Data.Cur == nil || len(r.Data.Cur) < 1 {
		log.Fatal("No weather data available.")
	}
	out := g.formatCond(make([]string, 5), r.Data.Cur[0], true)
	for _, val := range out {
		if g.config.RightToLeft {
			fmt.Fprint(stdout, strings.Repeat(" ", 94))
		} else {
			fmt.Fprint(stdout, " ")
		}
		fmt.Fprintln(stdout, val)
	}

	if g.config.Numdays == 0 {
		return nil
	}
	if r.Data.Weather == nil {
		log.Fatal("No detailed weather forecast available.")
	}
	for _, d := range r.Data.Weather {
		lines, err := g.printDay(d)
		if err != nil {
			return err
		}
		for _, val := range lines {
			fmt.Fprintln(stdout, val)
		}
	}

	return nil
}
