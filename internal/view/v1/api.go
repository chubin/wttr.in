//nolint:forbidigo,funlen,nestif,goerr113,gocognit,cyclop
package v1

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

//nolint:tagliatelle
type cond struct {
	ChanceOfRain   string  `json:"chanceofrain"`
	FeelsLikeC     int     `json:",string"`
	PrecipMM       float32 `json:"precipMM,string"`
	TempC          int     `json:"tempC,string"`
	TempC2         int     `json:"temp_C,string"`
	Time           int     `json:"time,string"`
	VisibleDistKM  int     `json:"visibility,string"`
	WeatherCode    int     `json:"weatherCode,string"`
	WeatherDesc    []struct{ Value string }
	WindGustKmph   int `json:",string"`
	Winddir16Point string
	WindspeedKmph  int `json:"windspeedKmph,string"`
}

type astro struct {
	Moonrise string
	Moonset  string
	Sunrise  string
	Sunset   string
}

type weather struct {
	Astronomy []astro
	Date      string
	Hourly    []cond
	MaxtempC  int `json:"maxtempC,string"`
	MintempC  int `json:"mintempC,string"`
}

type loc struct {
	Query string `json:"query"`
	Type  string `json:"type"`
}

//nolint:tagliatelle
type resp struct {
	Data struct {
		Cur     []cond                 `json:"current_condition"`
		Err     []struct{ Msg string } `json:"error"`
		Req     []loc                  `json:"request"`
		Weather []weather              `json:"weather"`
	} `json:"data"`
}

func (g *global) getDataFromAPI() (*resp, error) {
	var (
		ret    resp
		params []string
	)

	if len(g.config.APIKey) == 0 {
		return nil, fmt.Errorf("no API key specified. Setup instructions are in the README")
	}
	params = append(params, "key="+g.config.APIKey)

	// non-flag shortcut arguments will overwrite possible flag arguments
	for _, arg := range flag.Args() {
		if v, err := strconv.Atoi(arg); err == nil && len(arg) == 1 {
			g.config.Numdays = v
		} else {
			g.config.City = arg
		}
	}

	if len(g.config.City) > 0 {
		params = append(params, "q="+url.QueryEscape(g.config.City))
	}
	params = append(params, "format=json", "num_of_days="+strconv.Itoa(g.config.Numdays), "tp=3")
	if g.config.Lang != "" {
		params = append(params, "lang="+g.config.Lang)
	}

	if g.debug {
		fmt.Fprintln(os.Stderr, params)
	}

	res, err := http.Get(wuri + strings.Join(params, "&"))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if g.debug {
		var out bytes.Buffer

		err := json.Indent(&out, body, "", "  ")
		if err != nil {
			return nil, err
		}

		_, err = out.WriteTo(os.Stderr)
		if err != nil {
			return nil, err
		}

		fmt.Print("\n\n")
	}

	if g.config.Lang == "" {
		if err = json.Unmarshal(body, &ret); err != nil {
			return nil, err
		}
	} else {
		if err = g.unmarshalLang(body, &ret); err != nil {
			return nil, err
		}
	}

	return &ret, nil
}

func (g *global) unmarshalLang(body []byte, r *resp) error {
	var rv map[string]interface{}
	if err := json.Unmarshal(body, &rv); err != nil {
		return err
	}
	if data, ok := rv["data"].(map[string]interface{}); ok {
		if ccs, ok := data["current_condition"].([]interface{}); ok {
			for _, cci := range ccs {
				cc, ok := cci.(map[string]interface{})
				if !ok {
					continue
				}
				langs, ok := cc["lang_"+g.config.Lang].([]interface{})
				if !ok || len(langs) == 0 {
					continue
				}
				weatherDesc, ok := cc["weatherDesc"].([]interface{})
				if !ok || len(weatherDesc) == 0 {
					continue
				}
				weatherDesc[0] = langs[0]
			}
		}
		if ws, ok := data["weather"].([]interface{}); ok {
			for _, wi := range ws {
				w, ok := wi.(map[string]interface{})
				if !ok {
					continue
				}
				if hs, ok := w["hourly"].([]interface{}); ok {
					for _, hi := range hs {
						h, ok := hi.(map[string]interface{})
						if !ok {
							continue
						}
						langs, ok := h["lang_"+g.config.Lang].([]interface{})
						if !ok || len(langs) == 0 {
							continue
						}
						weatherDesc, ok := h["weatherDesc"].([]interface{})
						if !ok || len(weatherDesc) == 0 {
							continue
						}
						weatherDesc[0] = langs[0]
					}
				}
			}
		}
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(rv); err != nil {
		return err
	}
	if err := json.NewDecoder(&buf).Decode(r); err != nil {
		return err
	}

	return nil
}
