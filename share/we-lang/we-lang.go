package main

import (
	"bytes"
	_ "crypto/sha512"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/klauspost/lctime"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-runewidth"
)

type configuration struct {
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

type resp struct {
	Data struct {
		Cur     []cond                 `json:"current_condition"`
		Err     []struct{ Msg string } `json:"error"`
		Req     []loc                  `json:"request"`
		Weather []weather              `json:"weather"`
	} `json:"data"`
}

var (
	ansiEsc    *regexp.Regexp
	config     configuration
	configpath string
	debug      bool
	windDir    = map[string]string{
		"N":   "\033[1m↓\033[0m",
		"NNE": "\033[1m↓\033[0m",
		"NE":  "\033[1m↙\033[0m",
		"ENE": "\033[1m↙\033[0m",
		"E":   "\033[1m←\033[0m",
		"ESE": "\033[1m←\033[0m",
		"SE":  "\033[1m↖\033[0m",
		"SSE": "\033[1m↖\033[0m",
		"S":   "\033[1m↑\033[0m",
		"SSW": "\033[1m↑\033[0m",
		"SW":  "\033[1m↗\033[0m",
		"WSW": "\033[1m↗\033[0m",
		"W":   "\033[1m→\033[0m",
		"WNW": "\033[1m→\033[0m",
		"NW":  "\033[1m↘\033[0m",
		"NNW": "\033[1m↘\033[0m",
	}
	unitRain = map[bool]string{
		false: "mm",
		true:  "in",
	}
	unitTemp = map[bool]string{
		false: "C",
		true:  "F",
	}
	unitVis = map[bool]string{
		false: "km",
		true:  "mi",
	}
	unitWind = map[int]string{
		0: "km/h",
		1: "mph",
		2: "m/s",
	}
	slotTimes = [slotcount]int{9 * 60, 12 * 60, 18 * 60, 22 * 60}
	codes     = map[int][]string{
		113: iconSunny,
		116: iconPartlyCloudy,
		119: iconCloudy,
		122: iconVeryCloudy,
		143: iconFog,
		176: iconLightShowers,
		179: iconLightSleetShowers,
		182: iconLightSleet,
		185: iconLightSleet,
		200: iconThunderyShowers,
		227: iconLightSnow,
		230: iconHeavySnow,
		248: iconFog,
		260: iconFog,
		263: iconLightShowers,
		266: iconLightRain,
		281: iconLightSleet,
		284: iconLightSleet,
		293: iconLightRain,
		296: iconLightRain,
		299: iconHeavyShowers,
		302: iconHeavyRain,
		305: iconHeavyShowers,
		308: iconHeavyRain,
		311: iconLightSleet,
		314: iconLightSleet,
		317: iconLightSleet,
		320: iconLightSnow,
		323: iconLightSnowShowers,
		326: iconLightSnowShowers,
		329: iconHeavySnow,
		332: iconHeavySnow,
		335: iconHeavySnowShowers,
		338: iconHeavySnow,
		350: iconLightSleet,
		353: iconLightShowers,
		356: iconHeavyShowers,
		359: iconHeavyRain,
		362: iconLightSleetShowers,
		365: iconLightSleetShowers,
		368: iconLightSnowShowers,
		371: iconHeavySnowShowers,
		374: iconLightSleetShowers,
		377: iconLightSleet,
		386: iconThunderyShowers,
		389: iconThunderyHeavyRain,
		392: iconThunderySnowShowers,
		395: iconHeavySnowShowers,
	}

	iconUnknown = []string{
		"    .-.      ",
		"     __)     ",
		"    (        ",
		"     `-’     ",
		"      •      "}
	iconSunny = []string{
		"\033[38;5;226m    \\   /    \033[0m",
		"\033[38;5;226m     .-.     \033[0m",
		"\033[38;5;226m  ― (   ) ―  \033[0m",
		"\033[38;5;226m     `-’     \033[0m",
		"\033[38;5;226m    /   \\    \033[0m"}
	iconPartlyCloudy = []string{
		"\033[38;5;226m   \\  /\033[0m      ",
		"\033[38;5;226m _ /\"\"\033[38;5;250m.-.    \033[0m",
		"\033[38;5;226m   \\_\033[38;5;250m(   ).  \033[0m",
		"\033[38;5;226m   /\033[38;5;250m(___(__) \033[0m",
		"             "}
	iconCloudy = []string{
		"             ",
		"\033[38;5;250m     .--.    \033[0m",
		"\033[38;5;250m  .-(    ).  \033[0m",
		"\033[38;5;250m (___.__)__) \033[0m",
		"             "}
	iconVeryCloudy = []string{
		"             ",
		"\033[38;5;240;1m     .--.    \033[0m",
		"\033[38;5;240;1m  .-(    ).  \033[0m",
		"\033[38;5;240;1m (___.__)__) \033[0m",
		"             "}
	iconLightShowers = []string{
		"\033[38;5;226m _`/\"\"\033[38;5;250m.-.    \033[0m",
		"\033[38;5;226m  ,\\_\033[38;5;250m(   ).  \033[0m",
		"\033[38;5;226m   /\033[38;5;250m(___(__) \033[0m",
		"\033[38;5;111m     ‘ ‘ ‘ ‘ \033[0m",
		"\033[38;5;111m    ‘ ‘ ‘ ‘  \033[0m"}
	iconHeavyShowers = []string{
		"\033[38;5;226m _`/\"\"\033[38;5;240;1m.-.    \033[0m",
		"\033[38;5;226m  ,\\_\033[38;5;240;1m(   ).  \033[0m",
		"\033[38;5;226m   /\033[38;5;240;1m(___(__) \033[0m",
		"\033[38;5;21;1m   ‚‘‚‘‚‘‚‘  \033[0m",
		"\033[38;5;21;1m   ‚’‚’‚’‚’  \033[0m"}
	iconLightSnowShowers = []string{
		"\033[38;5;226m _`/\"\"\033[38;5;250m.-.    \033[0m",
		"\033[38;5;226m  ,\\_\033[38;5;250m(   ).  \033[0m",
		"\033[38;5;226m   /\033[38;5;250m(___(__) \033[0m",
		"\033[38;5;255m     *  *  * \033[0m",
		"\033[38;5;255m    *  *  *  \033[0m"}
	iconHeavySnowShowers = []string{
		"\033[38;5;226m _`/\"\"\033[38;5;240;1m.-.    \033[0m",
		"\033[38;5;226m  ,\\_\033[38;5;240;1m(   ).  \033[0m",
		"\033[38;5;226m   /\033[38;5;240;1m(___(__) \033[0m",
		"\033[38;5;255;1m    * * * *  \033[0m",
		"\033[38;5;255;1m   * * * *   \033[0m"}
	iconLightSleetShowers = []string{
		"\033[38;5;226m _`/\"\"\033[38;5;250m.-.    \033[0m",
		"\033[38;5;226m  ,\\_\033[38;5;250m(   ).  \033[0m",
		"\033[38;5;226m   /\033[38;5;250m(___(__) \033[0m",
		"\033[38;5;111m     ‘ \033[38;5;255m*\033[38;5;111m ‘ \033[38;5;255m* \033[0m",
		"\033[38;5;255m    *\033[38;5;111m ‘ \033[38;5;255m*\033[38;5;111m ‘  \033[0m"}
	iconThunderyShowers = []string{
		"\033[38;5;226m _`/\"\"\033[38;5;250m.-.    \033[0m",
		"\033[38;5;226m  ,\\_\033[38;5;250m(   ).  \033[0m",
		"\033[38;5;226m   /\033[38;5;250m(___(__) \033[0m",
		"\033[38;5;228;5m    ⚡\033[38;5;111;25m‘‘\033[38;5;228;5m⚡\033[38;5;111;25m‘‘ \033[0m",
		"\033[38;5;111m    ‘ ‘ ‘ ‘  \033[0m"}
	iconThunderyHeavyRain = []string{
		"\033[38;5;240;1m     .-.     \033[0m",
		"\033[38;5;240;1m    (   ).   \033[0m",
		"\033[38;5;240;1m   (___(__)  \033[0m",
		"\033[38;5;21;1m  ‚‘\033[38;5;228;5m⚡\033[38;5;21;25m‘‚\033[38;5;228;5m⚡\033[38;5;21;25m‚‘ \033[0m",
		"\033[38;5;21;1m  ‚’‚’\033[38;5;228;5m⚡\033[38;5;21;25m’‚’  \033[0m"}
	iconThunderySnowShowers = []string{
		"\033[38;5;226m _`/\"\"\033[38;5;250m.-.    \033[0m",
		"\033[38;5;226m  ,\\_\033[38;5;250m(   ).  \033[0m",
		"\033[38;5;226m   /\033[38;5;250m(___(__) \033[0m",
		"\033[38;5;255m     *\033[38;5;228;5m⚡\033[38;5;255;25m*\033[38;5;228;5m⚡\033[38;5;255;25m* \033[0m",
		"\033[38;5;255m    *  *  *  \033[0m"}
	iconLightRain = []string{
		"\033[38;5;250m     .-.     \033[0m",
		"\033[38;5;250m    (   ).   \033[0m",
		"\033[38;5;250m   (___(__)  \033[0m",
		"\033[38;5;111m    ‘ ‘ ‘ ‘  \033[0m",
		"\033[38;5;111m   ‘ ‘ ‘ ‘   \033[0m"}
	iconHeavyRain = []string{
		"\033[38;5;240;1m     .-.     \033[0m",
		"\033[38;5;240;1m    (   ).   \033[0m",
		"\033[38;5;240;1m   (___(__)  \033[0m",
		"\033[38;5;21;1m  ‚‘‚‘‚‘‚‘   \033[0m",
		"\033[38;5;21;1m  ‚’‚’‚’‚’   \033[0m"}
	iconLightSnow = []string{
		"\033[38;5;250m     .-.     \033[0m",
		"\033[38;5;250m    (   ).   \033[0m",
		"\033[38;5;250m   (___(__)  \033[0m",
		"\033[38;5;255m    *  *  *  \033[0m",
		"\033[38;5;255m   *  *  *   \033[0m"}
	iconHeavySnow = []string{
		"\033[38;5;240;1m     .-.     \033[0m",
		"\033[38;5;240;1m    (   ).   \033[0m",
		"\033[38;5;240;1m   (___(__)  \033[0m",
		"\033[38;5;255;1m   * * * *   \033[0m",
		"\033[38;5;255;1m  * * * *    \033[0m"}
	iconLightSleet = []string{
		"\033[38;5;250m     .-.     \033[0m",
		"\033[38;5;250m    (   ).   \033[0m",
		"\033[38;5;250m   (___(__)  \033[0m",
		"\033[38;5;111m    ‘ \033[38;5;255m*\033[38;5;111m ‘ \033[38;5;255m*  \033[0m",
		"\033[38;5;255m   *\033[38;5;111m ‘ \033[38;5;255m*\033[38;5;111m ‘   \033[0m"}
	iconFog = []string{
		"             ",
		"\033[38;5;251m _ - _ - _ - \033[0m",
		"\033[38;5;251m  _ - _ - _  \033[0m",
		"\033[38;5;251m _ - _ - _ - \033[0m",
		"             "}

	locale = map[string]string{
		"af":     "af_ZA",
		"am":     "am_ET",
		"ar":     "ar_TN",
		"az":     "az_AZ",
		"be":     "be_BY",
		"bg":     "bg_BG",
		"bn":     "bn_IN",
		"bs":     "bs_BA",
		"ca":     "ca_ES",
		"cs":     "cs_CZ",
		"cy":     "cy_GB",
		"da":     "da_DK",
		"de":     "de_DE",
		"el":     "el_GR",
		"eo":     "eo",
		"es":     "es_ES",
		"et":     "et_EE",
		"eu":     "eu_ES",
		"fa":     "fa_IR",
		"fi":     "fi_FI",
		"fr":     "fr_FR",
		"fy":     "fy_NL",
		"ga":     "ga_IE",
		"he":     "he_IL",
		"hi":     "hi_IN",
		"hr":     "hr_HR",
		"hu":     "hu_HU",
		"hy":     "hy_AM",
		"ia":     "ia",
		"id":     "id_ID",
		"is":     "is_IS",
		"it":     "it_IT",
		"ja":     "ja_JP",
		"jv":     "en_US",
		"ka":     "ka_GE",
		"kk":     "kk_KZ",
		"ko":     "ko_KR",
		"ky":     "ky_KG",
		"lt":     "lt_LT",
		"lv":     "lv_LV",
		"mg":     "mg_MG",
		"mk":     "mk_MK",
		"ml":     "ml_IN",
		"nb":     "nb_NO",
		"nl":     "nl_NL",
		"nn":     "nn_NO",
		"oc":     "oc_FR",
		"pl":     "pl_PL",
		"pt-br":  "pt_BR",
		"pt":     "pt_PT",
		"ro":     "ro_RO",
		"ru":     "ru_RU",
		"sk":     "sk_SK",
		"sl":     "sl_SI",
		"sr-lat": "sr_RS@latin",
		"sr":     "sr_RS",
		"sv":     "sv_SE",
		"sw":     "sw_KE",
		"ta":     "ta_IN",
		"th":     "th_TH",
		"tr":     "tr_TR",
		"uk":     "uk_UA",
		"uz":     "uz_UZ",
		"vi":     "vi_VN",
		"zh-cn":  "zh_CN",
		"zh-tw":  "zh_TW",
		"zh":     "zh_CN",
		"zu":     "zu_ZA",
	}

	localizedCaption = map[string]string{
		"af":     "Weer verslag vir:",
		"am":     "የአየር ሁኔታ ዘገባ ለ ፥",
		"ar":     "تقرير حالة ألطقس",
		"az":     "Hava proqnozu:",
		"be":     "Прагноз надвор'я для:",
		"bg":     "Прогноза за времето в:",
		"bn":     "আবহাওয়া সঙ্ক্রান্ত তথ্য",
		"bs":     "Vremenske prognoze za:",
		"ca":     "Informe del temps per a:",
		"cs":     "Předpověď počasí pro:",
		"cy":     "Adroddiad tywydd ar gyfer:",
		"da":     "Vejret i:",
		"de":     "Wetterbericht für:",
		"el":     "Πρόγνωση καιρού για:",
		"eo":     "Veterprognozo por:",
		"es":     "El tiempo en:",
		"et":     "Ilmaprognoos:",
		"eu":     "Eguraldia:",
		"fa":     "اوه و بآ تیعضو شرازگ",
		"fi":     "Säätiedotus:",
		"fr":     "Prévisions météo pour:",
		"fy":     "Waarberjocht foar:",
		"ga":     "Réamhaisnéis na haimsire do:",
		"he":     ":ריוואה גזמ תיזחת",
		"hi":     "मौसम की जानकारी",
		"hr":     "Vremenska prognoza za:",
		"hu":     "Időjárás előrejelzés:",
		"hy":     "Եղանակի տեսություն:",
		"ia":     "Le tempore a:",
		"id":     "Prakiraan cuaca:",
		"it":     "Previsioni meteo:",
		"is":     "Veðurskýrsla fyrir:",
		"ja":     "天気予報：",
		"jv":     "Weather forecast for:",
		"ka":     "ამინდის პროგნოზი:",
		"kk":     "Ауа райы:",
		"ko":     "일기 예보：",
		"ky":     "Аба ырайы:",
		"lt":     "Orų prognozė:",
		"lv":     "Laika ziņas:",
		"mk":     "Прогноза за времето во:",
		"ml":     "കാലാവസ്ഥ റിപ്പോർട്ട്:",
		"mr":     "हवामानाचा अंदाज:",
		"nb":     "Værmelding for:",
		"nl":     "Weerbericht voor:",
		"nn":     "Vêrmelding for:",
		"oc":     "Previsions metèo per:",
		"pl":     "Pogoda w:",
		"pt":     "Previsão do tempo para:",
		"pt-br":  "Previsão do tempo para:",
		"ro":     "Prognoza meteo pentru:",
		"ru":     "Прогноз погоды:",
		"sk":     "Predpoveď počasia pre:",
		"sl":     "Vremenska napoved za",
		"sr":     "Временска прогноза за:",
		"sr-lat": "Vremenska prognoza za:",
		"sv":     "Väderleksprognos för:",
		"sw":     "Ripoti ya hali ya hewa, jiji la:",
		"ta":     "வானிலை அறிக்கை",
		"te":     "వాతావరణ సమాచారము:",
		"th":     "รายงานสภาพอากาศ:",
		"tr":     "Hava beklentisi:",
		"uk":     "Прогноз погоди для:",
		"uz":     "Ob-havo bashorati:",
		"vi":     "Báo cáo thời tiết:",
		"zu":     "Isimo sezulu:",
		"zh":     "天气预报：",
		"zh-cn":  "天气预报：",
		"zh-tw":  "天氣預報：",
		"mg":     "Vinavina toetr'andro hoan'ny:",
	}

	daytimeTranslation = map[string][]string{
		"af":     {"Oggend", "Middag", "Vroegaand", "Laatnag"},
		"am":     {"ጠዋት", "ከሰዓት በኋላ", "ምሽት", "ሌሊት"},
		"ar":     {"ﺎﻠﻠﻴﻟ", "ﺎﻠﻤﺳﺍﺀ", "ﺎﻠﻈﻫﺭ", "ﺎﻠﺼﺑﺎﺣ"},
		"az":     {"Səhər", "Gün", "Axşam", "Gecə"},
		"be":     {"Раніца", "Дзень", "Вечар", "Ноч"},
		"bg":     {"Сутрин", "Обяд", "Вечер", "Нощ"},
		"bn":     {"সকাল", "দুপুর", "সন্ধ্যা", "রাত্রি"},
		"bs":     {"Ujutro", "Dan", "Večer", "Noć"},
		"cs":     {"Ráno", "Ve dne", "Večer", "V noci"},
		"ca":     {"Matí", "Dia", "Tarda", "Nit"},
		"cy":     {"Bore", "Dydd", "Hwyr", "Nos"},
		"da":     {"Morgen", "Middag", "Aften", "Nat"},
		"de":     {"Früh", "Mittag", "Abend", "Nacht"},
		"el":     {"Πρωί", "Μεσημέρι", "Απόγευμα", "Βράδυ"},
		"en":     {"Morning", "Noon", "Evening", "Night"},
		"eo":     {"Mateno", "Tago", "Vespero", "Nokto"},
		"es":     {"Mañana", "Mediodía", "Tarde", "Noche"},
		"et":     {"Hommik", "Päev", "Õhtu", "Öösel"},
		"eu":     {"Goiza", "Eguerdia", "Arratsaldea", "Gaua"},
		"fa":     {"حبص", "رهظ", "رصع", "بش"},
		"fi":     {"Aamu", "Keskipäivä", "Ilta", "Yö"},
		"fr":     {"Matin", "Après-midi", "Soir", "Nuit"},
		"fy":     {"Moarns", "Middeis", "Jûns", "Nachts"},
		"ga":     {"Maidin", "Nóin", "Tráthnóna", "Oíche"},
		"he":     {"רקוב", "םוֹיְ", "ברֶעֶ", "הלָיְלַ"},
		"hi":     {"प्रातःकाल", "दोपहर", "सायंकाल", "रात"},
		"hr":     {"Jutro", "Dan", "Večer", "Noć"},
		"hu":     {"Reggel", "Dél", "Este", "Éjszaka"},
		"hy":     {"Առավոտ", "Կեսօր", "Երեկո", "Գիշեր"},
		"ia":     {"Matino", "Mediedie", "Vespere", "Nocte"},
		"id":     {"Pagi", "Hari", "Petang", "Malam"},
		"it":     {"Mattina", "Pomeriggio", "Sera", "Notte"},
		"is":     {"Morgunn", "Dagur", "Kvöld", "Nótt"},
		"ja":     {"朝", "昼", "夕", "夜"},
		"jv":     {"Morning", "Noon", "Evening", "Night"},
		"ka":     {"დილა", "დღე", "საღამო", "ღამე"},
		"kk":     {"Таң", "Күндіз", "Кеш", "Түн"},
		"ko":     {"아침", "낮", "저녁", "밤"},
		"ky":     {"Эртең", "Күн", "Кеч", "Түн"},
		"lt":     {"Rytas", "Diena", "Vakaras", "Naktis"},
		"lv":     {"Rīts", "Diena", "Vakars", "Nakts"},
		"mk":     {"Утро", "Пладне", "Вечер", "Ноќ"},
		"ml":     {"രാവിലെ", "മധ്യാഹ്നം", "വൈകുന്നേരം", "രാത്രി"},
		"mr":     {"सकाळ", "दुपार", "संध्याकाळ", "रात्र"},
		"nl":     {"'s Ochtends", "'s Middags", "'s Avonds", "'s Nachts"},
		"nb":     {"Morgen", "Middag", "Kveld", "Natt"},
		"nn":     {"Morgon", "Middag", "Kveld", "Natt"},
		"oc":     {"Matin", "Jorn", "Vèspre", "Nuèch"},
		"pl":     {"Ranek", "Dzień", "Wieczór", "Noc"},
		"pt":     {"Manhã", "Meio-dia", "Tarde", "Noite"},
		"pt-br":  {"Manhã", "Meio-dia", "Tarde", "Noite"},
		"ro":     {"Dimineaţă", "Amiază", "Seară", "Noapte"},
		"ru":     {"Утро", "День", "Вечер", "Ночь"},
		"sk":     {"Ráno", "Cez deň", "Večer", "V noci"},
		"sl":     {"Jutro", "Dan", "Večer", "Noč"},
		"sr":     {"Јутро", "Подне", "Вече", "Ноћ"},
		"sr-lat": {"Jutro", "Podne", "Veče", "Noć"},
		"sv":     {"Morgon", "Eftermiddag", "Kväll", "Natt"},
		"sw":     {"Asubuhi", "Adhuhuri", "Jioni", "Usiku"},
		"ta":     {"காலை", "நண்பகல்", "சாயங்காலம்", "இரவு"},
		"te":     {"ఉదయం", "రోజు", "సాయంత్రం", "రాత్రి"},
		"th":     {"เช้า", "วัน", "เย็น", "คืน"},
		"tr":     {"Sabah", "Öğle", "Akşam", "Gece"},
		"uk":     {"Ранок", "День", "Вечір", "Ніч"},
		"uz":     {"Ertalab", "Kunduzi", "Kechqurun", "Kecha"},
		"vi":     {"Sáng", "Trưa", "Chiều", "Tối"},
		"zh":     {"早上", "中午", "傍晚", "夜间"},
		"zh-cn":  {"早上", "中午", "傍晚", "夜间"},
		"zh-tw":  {"早上", "中午", "傍晚", "夜間"},
		"zu":     {"Morning", "Noon", "Evening", "Night"},
		"mg":     {"Maraina", "Tolakandro", "Ariva", "Alina"},
	}
)

// Add this languages:
// da tr hu sr jv zu
// More languages: https://developer.worldweatheronline.com/api/multilingual.aspx

// const (
// 	wuri      = "https://api.worldweatheronline.com/premium/v1/weather.ashx?"
// 	suri      = "https://api.worldweatheronline.com/premium/v1/search.ashx?"
// 	slotcount = 4
// )

const (
	wuri      = "http://127.0.0.1:5001/premium/v1/weather.ashx?"
	suri      = "http://127.0.0.1:5001/premium/v1/search.ashx?"
	slotcount = 4
)

func configload() error {
	b, err := ioutil.ReadFile(configpath)
	if err == nil {
		return json.Unmarshal(b, &config)
	}
	return err
}

func configsave() error {
	j, err := json.MarshalIndent(config, "", "\t")
	if err == nil {
		return ioutil.WriteFile(configpath, j, 0600)
	}
	return err
}

func pad(s string, mustLen int) (ret string) {
	ret = s
	realLen := utf8.RuneCountInString(ansiEsc.ReplaceAllLiteralString(s, ""))
	delta := mustLen - realLen
	if delta > 0 {
		if config.RightToLeft {
			ret = strings.Repeat(" ", delta) + ret + "\033[0m"
		} else {
			ret += "\033[0m" + strings.Repeat(" ", delta)
		}
	} else if delta < 0 {
		toks := ansiEsc.Split(s, 2)
		tokLen := utf8.RuneCountInString(toks[0])
		esc := ansiEsc.FindString(s)
		if tokLen > mustLen {
			ret = fmt.Sprintf("%.*s\033[0m", mustLen, toks[0])
		} else {
			ret = fmt.Sprintf("%s%s%s", toks[0], esc, pad(toks[1], mustLen-tokLen))
		}
	}
	return
}

func formatTemp(c cond) string {
	color := func(temp int, explicitPlus bool) string {
		var col = 0
		if !config.Inverse {
			// Extemely cold temperature must be shown with violet
			// because dark blue is too dark
			col = 165
			switch temp {
			case -15, -14, -13:
				col = 171
			case -12, -11, -10:
				col = 33
			case -9, -8, -7:
				col = 39
			case -6, -5, -4:
				col = 45
			case -3, -2, -1:
				col = 51
			case 0, 1:
				col = 50
			case 2, 3:
				col = 49
			case 4, 5:
				col = 48
			case 6, 7:
				col = 47
			case 8, 9:
				col = 46
			case 10, 11, 12:
				col = 82
			case 13, 14, 15:
				col = 118
			case 16, 17, 18:
				col = 154
			case 19, 20, 21:
				col = 190
			case 22, 23, 24:
				col = 226
			case 25, 26, 27:
				col = 220
			case 28, 29, 30:
				col = 214
			case 31, 32, 33:
				col = 208
			case 34, 35, 36:
				col = 202
			default:
				if temp > 0 {
					col = 196
				}
			}
		} else {
			col = 16
			switch temp {
			case -15, -14, -13:
				col = 17
			case -12, -11, -10:
				col = 18
			case -9, -8, -7:
				col = 19
			case -6, -5, -4:
				col = 20
			case -3, -2, -1:
				col = 21
			case 0, 1:
				col = 30
			case 2, 3:
				col = 28
			case 4, 5:
				col = 29
			case 6, 7:
				col = 30
			case 8, 9:
				col = 34
			case 10, 11, 12:
				col = 35
			case 13, 14, 15:
				col = 36
			case 16, 17, 18:
				col = 40
			case 19, 20, 21:
				col = 59
			case 22, 23, 24:
				col = 100
			case 25, 26, 27:
				col = 101
			case 28, 29, 30:
				col = 94
			case 31, 32, 33:
				col = 166
			case 34, 35, 36:
				col = 52
			default:
				if temp > 0 {
					col = 196
				}
			}
		}
		if config.Imperial {
			temp = (temp*18 + 320) / 10
		}
		if explicitPlus {
			return fmt.Sprintf("\033[38;5;%03dm+%d\033[0m", col, temp)
		}
		return fmt.Sprintf("\033[38;5;%03dm%d\033[0m", col, temp)
	}
	t := c.TempC
	if t == 0 {
		t = c.TempC2
	}

	// hyphen := " - "

	// if (config.Lang == "sl") {
	//     hyphen = "-"
	// }

	// hyphen = ".."

	explicitPlus1 := false
	explicitPlus2 := false
	if c.FeelsLikeC != t {
		if t > 0 {
			explicitPlus1 = true
		}
		if c.FeelsLikeC > 0 {
			explicitPlus2 = true
		}
		if explicitPlus1 {
			explicitPlus2 = false
		}
		return pad(
			fmt.Sprintf("%s(%s) °%s",
				color(t, explicitPlus1),
				color(c.FeelsLikeC, explicitPlus2),
				unitTemp[config.Imperial]),
			15)
	}
	// if c.FeelsLikeC < t {
	// 	if c.FeelsLikeC < 0 && t > 0 {
	// 		explicitPlus = true
	// 	}
	// 	return pad(fmt.Sprintf("%s%s%s °%s", color(c.FeelsLikeC, false), hyphen, color(t, explicitPlus), unitTemp[config.Imperial]), 15)
	// } else if c.FeelsLikeC > t {
	// 	if t < 0 && c.FeelsLikeC > 0 {
	// 		explicitPlus = true
	// 	}
	// 	return pad(fmt.Sprintf("%s%s%s °%s", color(t, false), hyphen, color(c.FeelsLikeC, explicitPlus), unitTemp[config.Imperial]), 15)
	// }
	return pad(fmt.Sprintf("%s °%s", color(c.FeelsLikeC, false), unitTemp[config.Imperial]), 15)
}

func formatWind(c cond) string {
	windInRightUnits := func(spd int) int {
		if config.WindMS {
			spd = (spd * 1000) / 3600
		} else {
			if config.Imperial {
				spd = (spd * 1000) / 1609
			}
		}
		return spd
	}
	color := func(spd int) string {
		var col = 46
		switch spd {
		case 1, 2, 3:
			col = 82
		case 4, 5, 6:
			col = 118
		case 7, 8, 9:
			col = 154
		case 10, 11, 12:
			col = 190
		case 13, 14, 15:
			col = 226
		case 16, 17, 18, 19:
			col = 220
		case 20, 21, 22, 23:
			col = 214
		case 24, 25, 26, 27:
			col = 208
		case 28, 29, 30, 31:
			col = 202
		default:
			if spd > 0 {
				col = 196
			}
		}
		spd = windInRightUnits(spd)

		return fmt.Sprintf("\033[38;5;%03dm%d\033[0m", col, spd)
	}

	unitWindString := unitWind[0]
	if config.WindMS {
		unitWindString = unitWind[2]
	} else {
		if config.Imperial {
			unitWindString = unitWind[1]
		}
	}

	hyphen := " - "
	// if (config.Lang == "sl") {
	//     hyphen = "-"
	// }
	hyphen = "-"

	cWindGustKmph := color(c.WindGustKmph)
	cWindspeedKmph := color(c.WindspeedKmph)
	if windInRightUnits(c.WindGustKmph) > windInRightUnits(c.WindspeedKmph) {
		return pad(fmt.Sprintf("%s %s%s%s %s", windDir[c.Winddir16Point], cWindspeedKmph, hyphen, cWindGustKmph, unitWindString), 15)
	}
	return pad(fmt.Sprintf("%s %s %s", windDir[c.Winddir16Point], cWindspeedKmph, unitWindString), 15)
}

func formatVisibility(c cond) string {
	if config.Imperial {
		c.VisibleDistKM = (c.VisibleDistKM * 621) / 1000
	}
	return pad(fmt.Sprintf("%d %s", c.VisibleDistKM, unitVis[config.Imperial]), 15)
}

func formatRain(c cond) string {
	rainUnit := float32(c.PrecipMM)
	if config.Imperial {
		rainUnit = float32(c.PrecipMM) * 0.039
	}
	if c.ChanceOfRain != "" {
		return pad(fmt.Sprintf("%.1f %s | %s%%", rainUnit, unitRain[config.Imperial], c.ChanceOfRain), 15)
	}
	return pad(fmt.Sprintf("%.1f %s", rainUnit, unitRain[config.Imperial]), 15)
}

func formatCond(cur []string, c cond, current bool) (ret []string) {
	var icon []string
	if i, ok := codes[c.WeatherCode]; !ok {
		icon = iconUnknown
	} else {
		icon = i
	}
	if config.Inverse {
		// inverting colors
		for i := range icon {
			icon[i] = strings.Replace(icon[i], "38;5;226", "38;5;94", -1)
			icon[i] = strings.Replace(icon[i], "38;5;250", "38;5;243", -1)
			icon[i] = strings.Replace(icon[i], "38;5;21", "38;5;18", -1)
			icon[i] = strings.Replace(icon[i], "38;5;255", "38;5;245", -1)
			icon[i] = strings.Replace(icon[i], "38;5;111", "38;5;63", -1)
			icon[i] = strings.Replace(icon[i], "38;5;251", "38;5;238", -1)
		}
	}
	//desc := fmt.Sprintf("%-15.15v", c.WeatherDesc[0].Value)
	desc := c.WeatherDesc[0].Value
	if config.RightToLeft {
		for runewidth.StringWidth(desc) < 15 {
			desc = " " + desc
		}
		for runewidth.StringWidth(desc) > 15 {
			_, size := utf8.DecodeLastRuneInString(desc)
			desc = desc[size:]
		}
	} else {
		for runewidth.StringWidth(desc) < 15 {
			desc += " "
		}
		for runewidth.StringWidth(desc) > 15 {
			_, size := utf8.DecodeLastRuneInString(desc)
			desc = desc[:len(desc)-size]
		}
	}
	if current {
		if config.RightToLeft {
			desc = c.WeatherDesc[0].Value
			if runewidth.StringWidth(desc) < 15 {
				desc = strings.Repeat(" ", 15-runewidth.StringWidth(desc)) + desc
			}
		} else {
			desc = c.WeatherDesc[0].Value
		}
	} else {
		if config.RightToLeft {
			if frstRune, size := utf8.DecodeRuneInString(desc); frstRune != ' ' {
				desc = "…" + desc[size:]
				for runewidth.StringWidth(desc) < 15 {
					desc = " " + desc
				}
			}
		} else {
			if lastRune, size := utf8.DecodeLastRuneInString(desc); lastRune != ' ' {
				desc = desc[:len(desc)-size] + "…"
				//for numberOfSpaces < runewidth.StringWidth(fmt.Sprintf("%c", lastRune)) - 1 {
				for runewidth.StringWidth(desc) < 15 {
					desc = desc + " "
				}
			}
		}
	}
	if config.RightToLeft {
		ret = append(ret, fmt.Sprintf("%v %v %v", cur[0], desc, icon[0]))
		ret = append(ret, fmt.Sprintf("%v %v %v", cur[1], formatTemp(c), icon[1]))
		ret = append(ret, fmt.Sprintf("%v %v %v", cur[2], formatWind(c), icon[2]))
		ret = append(ret, fmt.Sprintf("%v %v %v", cur[3], formatVisibility(c), icon[3]))
		ret = append(ret, fmt.Sprintf("%v %v %v", cur[4], formatRain(c), icon[4]))
	} else {
		ret = append(ret, fmt.Sprintf("%v %v %v", cur[0], icon[0], desc))
		ret = append(ret, fmt.Sprintf("%v %v %v", cur[1], icon[1], formatTemp(c)))
		ret = append(ret, fmt.Sprintf("%v %v %v", cur[2], icon[2], formatWind(c)))
		ret = append(ret, fmt.Sprintf("%v %v %v", cur[3], icon[3], formatVisibility(c)))
		ret = append(ret, fmt.Sprintf("%v %v %v", cur[4], icon[4], formatRain(c)))
	}
	return
}

func justifyCenter(s string, width int) string {
	appendSide := 0
	for runewidth.StringWidth(s) <= width {
		if appendSide == 1 {
			s = s + " "
			appendSide = 0
		} else {
			s = " " + s
			appendSide = 1
		}
	}
	return s
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func printDay(w weather) (ret []string) {
	hourly := w.Hourly
	ret = make([]string, 5)
	for i := range ret {
		ret[i] = "│"
	}

	// find hourly data which fits the desired times of day best
	var slots [slotcount]cond
	for _, h := range hourly {
		c := int(math.Mod(float64(h.Time), 100)) + 60*(h.Time/100)
		for i, s := range slots {
			if math.Abs(float64(c-slotTimes[i])) < math.Abs(float64(s.Time-slotTimes[i])) {
				h.Time = c
				slots[i] = h
			}
		}
	}

	if config.RightToLeft {
		slots[0], slots[3] = slots[3], slots[0]
		slots[1], slots[2] = slots[2], slots[1]
	}

	for i, s := range slots {
		if config.Narrow {
			if i == 0 || i == 2 {
				continue
			}
		}
		ret = formatCond(ret, s, false)
		for i := range ret {
			ret[i] = ret[i] + "│"
		}
	}

	d, _ := time.Parse("2006-01-02", w.Date)
	// dateFmt := "┤ " + d.Format("Mon 02. Jan") + " ├"

	if val, ok := locale[config.Lang]; ok {
		lctime.SetLocale(val)
	} else {
		lctime.SetLocale("en_US")
	}
	dateName := ""
	if config.RightToLeft {
		dow := lctime.Strftime("%a", d)
		day := lctime.Strftime("%d", d)
		month := lctime.Strftime("%b", d)
		dateName = reverse(month) + " " + day + " " + reverse(dow)
	} else {
		dateName = lctime.Strftime("%a %d %b", d)
		if config.Lang == "ko" {
			dateName = lctime.Strftime("%b %d일 %a", d)
		}
		if config.Lang == "zh" || config.Lang == "zh-tw" || config.Lang == "zh-cn" {
			dateName = lctime.Strftime("%b%d日%A", d)
		}
	}
	// appendSide := 0
	// // for utf8.RuneCountInString(dateName) <= dateWidth {
	// for runewidth.StringWidth(dateName) <= dateWidth {
	//     if appendSide == 1 {
	//         dateName = dateName + " "
	//         appendSide = 0
	//     } else {
	//         dateName = " " + dateName
	//         appendSide = 1
	//     }
	// }

	dateFmt := "┤" + justifyCenter(dateName, 12) + "├"

	trans := daytimeTranslation["en"]
	if t, ok := daytimeTranslation[config.Lang]; ok {
		trans = t
	}
	if config.Narrow {

		names := "│      " + justifyCenter(trans[1], 16) +
			"└──────┬──────┘" + justifyCenter(trans[3], 16) + "      │"

		ret = append([]string{
			"                        ┌─────────────┐                        ",
			"┌───────────────────────" + dateFmt + "───────────────────────┐",
			names,
			"├──────────────────────────────┼──────────────────────────────┤"},
			ret...)

		return append(ret,
			"└──────────────────────────────┴──────────────────────────────┘")

	}

	names := ""
	if config.RightToLeft {
		names = "│" + justifyCenter(trans[3], 29) + "│      " + justifyCenter(trans[2], 16) +
			"└──────┬──────┘" + justifyCenter(trans[1], 16) + "      │" + justifyCenter(trans[0], 29) + "│"
	} else {
		names = "│" + justifyCenter(trans[0], 29) + "│      " + justifyCenter(trans[1], 16) +
			"└──────┬──────┘" + justifyCenter(trans[2], 16) + "      │" + justifyCenter(trans[3], 29) + "│"
	}

	ret = append([]string{
		"                                                       ┌─────────────┐                                                       ",
		"┌──────────────────────────────┬───────────────────────" + dateFmt + "───────────────────────┬──────────────────────────────┐",
		names,
		"├──────────────────────────────┼──────────────────────────────┼──────────────────────────────┼──────────────────────────────┤"},
		ret...)

	return append(ret,
		"└──────────────────────────────┴──────────────────────────────┴──────────────────────────────┴──────────────────────────────┘")
}

func unmarshalLang(body []byte, r *resp) error {
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
				langs, ok := cc["lang_"+config.Lang].([]interface{})
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
						langs, ok := h["lang_"+config.Lang].([]interface{})
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

func getDataFromAPI() (ret resp) {
	var params []string

	if len(config.APIKey) == 0 {
		log.Fatal("No API key specified. Setup instructions are in the README.")
	}
	params = append(params, "key="+config.APIKey)

	// non-flag shortcut arguments will overwrite possible flag arguments
	for _, arg := range flag.Args() {
		if v, err := strconv.Atoi(arg); err == nil && len(arg) == 1 {
			config.Numdays = v
		} else {
			config.City = arg
		}
	}

	if len(config.City) > 0 {
		params = append(params, "q="+url.QueryEscape(config.City))
	}
	params = append(params, "format=json")
	params = append(params, "num_of_days="+strconv.Itoa(config.Numdays))
	params = append(params, "tp=3")
	if config.Lang != "" {
		params = append(params, "lang="+config.Lang)
	}

	if debug {
		fmt.Fprintln(os.Stderr, params)
	}

	res, err := http.Get(wuri + strings.Join(params, "&"))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	if debug {
		var out bytes.Buffer
		json.Indent(&out, body, "", "  ")
		out.WriteTo(os.Stderr)
		fmt.Print("\n\n")
	}

	if config.Lang == "" {
		if err = json.Unmarshal(body, &ret); err != nil {
			log.Println(err)
		}
	} else {
		if err = unmarshalLang(body, &ret); err != nil {
			log.Println(err)
		}
	}
	return
}

func init() {
	flag.IntVar(&config.Numdays, "days", 3, "Number of days of weather forecast to be displayed")
	flag.StringVar(&config.Lang, "lang", "en", "Language of the report")
	flag.StringVar(&config.City, "city", "New York", "City to be queried")
	flag.BoolVar(&debug, "debug", false, "Print out raw json response for debugging purposes")
	flag.BoolVar(&config.Imperial, "imperial", false, "Use imperial units")
	flag.BoolVar(&config.Inverse, "inverse", false, "Use inverted colors")
	flag.BoolVar(&config.Narrow, "narrow", false, "Narrow output (two columns)")
	flag.StringVar(&config.LocationName, "location_name", "", "Location name (used in the caption)")
	flag.BoolVar(&config.WindMS, "wind_in_ms", false, "Show wind speed in m/s")
	flag.BoolVar(&config.RightToLeft, "right_to_left", false, "Right to left script")
	configpath = os.Getenv("WEGORC")
	if configpath == "" {
		usr, err := user.Current()
		if err != nil {
			log.Fatalf("%v\nYou can set the environment variable WEGORC to point to your config file as a workaround.", err)
		}
		configpath = path.Join(usr.HomeDir, ".wegorc")
	}
	config.APIKey = ""
	config.Imperial = false
	config.Lang = "en"
	err := configload()
	if _, ok := err.(*os.PathError); ok {
		log.Printf("No config file found. Creating %s ...", configpath)
		if err2 := configsave(); err2 != nil {
			log.Fatal(err2)
		}
	} else if err != nil {
		log.Fatalf("could not parse %v: %v", configpath, err)
	}

	ansiEsc = regexp.MustCompile("\033.*?m")
}

func main() {
	flag.Parse()

	r := getDataFromAPI()

	if r.Data.Req == nil || len(r.Data.Req) < 1 {
		if r.Data.Err != nil && len(r.Data.Err) >= 1 {
			log.Fatal(r.Data.Err[0].Msg)
		}
		log.Fatal("Malformed response.")
	}
	locationName := r.Data.Req[0].Query
	if config.LocationName != "" {
		locationName = config.LocationName
	}
	if config.Lang == "he" || config.Lang == "ar" || config.Lang == "fa" {
		config.RightToLeft = true
	}
	if caption, ok := localizedCaption[config.Lang]; !ok {
		// r.Data.Req[0].Type,
		fmt.Printf("Weather report: %s\n\n", locationName)
	} else {
		if config.RightToLeft {
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
	out := formatCond(make([]string, 5), r.Data.Cur[0], true)
	for _, val := range out {
		if config.RightToLeft {
			fmt.Fprint(stdout, strings.Repeat(" ", 94))
		} else {
			fmt.Fprint(stdout, " ")
		}
		fmt.Fprintln(stdout, val)
	}

	if config.Numdays == 0 {
		return
	}
	if r.Data.Weather == nil {
		log.Fatal("No detailed weather forecast available.")
	}
	for _, d := range r.Data.Weather {
		for _, val := range printDay(d) {
			fmt.Fprintln(stdout, val)
		}
	}
}
