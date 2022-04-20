*wttr.in — the right way to ~check~ `curl` the weather!*

wttr.in is a console-oriented weather forecast service that supports various information
representation methods like terminal-oriented ANSI-sequences for console HTTP clients
(curl, httpie, or wget), HTML for web browsers, or PNG for graphical viewers.

Originally started as a small project, a wrapper for [wego](https://github.com/schachmat/wego),
intended to demonstrate the power of the console-oriented services,
*wttr.in* became a popular weather reporting service, handling tens millions of queries daily.

You can see it running here: [wttr.in](https://wttr.in).

[Documentation](https://wttr.in/:help) | [Usage](https://github.com/chubin/wttr.in#usage) | [One-line output](https://github.com/chubin/wttr.in#one-line-output) | [Data-rich output format](https://github.com/chubin/wttr.in#data-rich-output-format-v2) | [Map view](https://github.com/chubin/wttr.in#map-view-v3) | [Output formats](https://github.com/chubin/wttr.in#different-output-formats) | [Moon phases](https://github.com/chubin/wttr.in#moon-phases) | [Internationalization](https://github.com/chubin/wttr.in#internationalization-and-localization) | [Windows issues](https://github.com/chubin/wttr.in#internationalization-and-localization) | [Installation](https://github.com/chubin/wttr.in#installation)

## Usage

You can access the service from a shell or from a Web browser like this:

    $ curl wttr.in
    Weather for City: Paris, France

         \   /     Clear
          .-.      10 – 11 °C  
       ― (   ) ―   ↑ 11 km/h  
          `-’      10 km  
         /   \     0.0 mm  


Here is an actual weather report for your location (it's live!):

![Weather Report](https://wttr.in/MyLocation.png?)

(It's not your actual location - GitHub's CDN hides your real IP address with its own IP address,
but it's still a live weather report in your language.)

Or in PowerShell:

```PowerShell
Invoke-RestMethod https://wttr.in
```

Want to get the weather information for a specific location? You can add the desired location to the URL in your
request like this:

    $ curl wttr.in/London
    $ curl wttr.in/Moscow
    $ curl wttr.in/Salt+Lake+City

If you omit the location name, you will get the report for your current location based on your IP address.

Use 3-letter airport codes in order to get the weather information at a certain airport:

    $ curl wttr.in/muc      # Weather for IATA: muc, Munich International Airport, Germany
    $ curl wttr.in/ham      # Weather for IATA: ham, Hamburg Airport, Germany

Let's say you'd like to get the weather for a geographical location other than a town or city - maybe an attraction
in a city, a mountain name, or some special location. Add the character `~` before the name to look up that special
location name before the weather is then retrieved:

	$ curl wttr.in/~Vostok+Station
	$ curl wttr.in/~Eiffel+Tower
	$ curl wttr.in/~Kilimanjaro

For these examples, you'll see a line below the weather forecast output that shows the geolocation
results of looking up the location:

	Location: Vostok Station, станция Восток, AAT, Antarctica [-78.4642714,106.8364678]
    Location: Tour Eiffel, 5, Avenue Anatole France, Gros-Caillou, 7e, Paris, Île-de-France, 75007, France [48.8582602,2.29449905432]
	Location: Kilimanjaro, Northern, Tanzania [-3.4762789,37.3872648]

You can also use IP-addresses (direct) or domain names (prefixed with `@`) to specify a location:

    $ curl wttr.in/@github.com
    $ curl wttr.in/@msu.ru

To get detailed information online, you can access the [/:help](https://wttr.in/:help) page:

    $ curl wttr.in/:help


### Weather Units

By default the USCS units are used for the queries from the USA and the metric system for the rest of the world.
You can override this behavior by adding `?u`, `?m` or `?M`   to a URL like this:

    $ curl wttr.in/Amsterdam?u  # USCS (used by default in US)
    $ curl wttr.in/Amsterdam?m  # metric (SI) (used by default everywhere except US)
    $ curl wttr.in/Amsterdam?M  # metric (SI), but show wind speed in m/s

If you have several options to pass, write them without delimiters in between for the one-letter options,
and use `&` as a delimiter for the long options with values:

    $ curl 'wttr.in/Amsterdam?m2&lang=nl'

It would be a rough equivalent of `-m2 --lang nl` for the GNU CLI syntax.

## Supported output formats and views

wttr.in currently supports five output formats:

* ANSI for the terminal;
* Plain-text for the terminal and scripts;
* HTML for the browser;
* PNG for the graphical viewers;
* JSON for scripts and APIs;
* Prometheus metrics for scripts and APIs.

The ANSI and HTML formats are selected based on the User-Agent string.

To force plain text, which disables colors:

    $ curl wttr.in/?T

The PNG format can be forced by adding `.png` to the end of the query:

    $ wget wttr.in/Paris.png

You can use all of the options with the PNG-format like in an URL, but you have
to separate them with `_` instead of `?` and `&`:

    $ wget wttr.in/Paris_0tqp_lang=fr.png

Useful options for the PNG format:

* `t` for transparency (`transparency=150`);
* transparency=0..255 for a custom transparency level.

Transparency is a useful feature when weather PNGs are used to add weather data to pictures:

    $ convert source.jpg <( curl wttr.in/Oymyakon_tqp0.png ) -geometry +50+50 -composite target.jpg

In this example:

* `source.jpg` - source file;
* `target.jpg` - target file;
* `Oymyakon` - name of the location;
* `tqp0` - options (recommended).

![Picture with weather data](https://pbs.twimg.com/media/C69-wsIW0AAcAD5.jpg)

You can embed a special wttr.in widget, that displays the weather condition for the current or a selected location, into a HTML page using the [wttr-switcher](https://github.com/midzer/wttr-switcher). That is how it looks like: [wttr-switcher-example](https://midzer.github.io/wttr-switcher/) or on a real world web site: https://feuerwehr-eisolzried.de/.

![Embedded wttr.in example at feuerwehr-eisolzried.de](https://user-images.githubusercontent.com/3875145/65265457-50eac180-db11-11e9-8f9b-2e1711dfc436.png)

## One-line output

One-line output format is convenient to be used to show weather info
in status bar of different programs, such as *tmux*, *weechat*, etc.

For one-line output format, specify additional URL parameter `format`:

```
$ curl wttr.in/Nuremberg?format=3
Nuremberg: 🌦 +11⁰C
```

Available preconfigured formats: 1, 2, 3, 4 and the custom format using the percent notation (see below).

You can specify multiple locations separated with `:` (for repeating queries):

```
$ curl wttr.in/Nuremberg:Hamburg:Berlin?format=3
Nuremberg: 🌦 +11⁰C
```
Or to process all this queries at once:

```
$ curl -s 'wttr.in/{Nuremberg,Hamburg,Berlin}?format=3'
Nuremberg: 🌦 +11⁰C
Hamburg: 🌦 +8⁰C
Berlin: 🌦 +8⁰C
```

To specify your own custom output format, use the special `%`-notation:

```
    c    Weather condition,
    C    Weather condition textual name,
    x    Weather condition, plain-text symbol,
    h    Humidity,
    t    Temperature (Actual),
    f    Temperature (Feels Like),
    w    Wind,
    l    Location,
    m    Moon phase 🌑🌒🌓🌔🌕🌖🌗🌘,
    M    Moon day,
    p    Precipitation (mm/3 hours),
    P    Pressure (hPa),

    D    Dawn*,
    S    Sunrise*,
    z    Zenith*,
    s    Sunset*,
    d    Dusk*,
    T    Current time*,
    Z    Local timezone.

(*times are shown in the local timezone)
```

So, these two calls are the same:

```
    $ curl wttr.in/London?format=3
    London: ⛅️ +7⁰C
    $ curl wttr.in/London?format="%l:+%c+%t\n"
    London: ⛅️ +7⁰C
```

### tmux

When using in `tmux.conf`, you have to escape `%` with `%`, i.e. write there `%%` instead of `%`.

The output does not contain new line by default, when the %-notation is used, but it does contain it when preconfigured format (`1`,`2`,`3` etc.)
are used. To have the new line in the output when the %-notation is used, use '\n' and single quotes when doing a query from the shell.

In programs, that are querying the service automatically (such as tmux), it is better to use some reasonable update interval. In tmux, you can configure it with `status-interval`.

If several, `:` separated locations, are specified in the query, specify update period
as an additional query parameter `period=`:
```
set -g status-interval 60
WEATHER='#(curl -s wttr.in/London:Stockholm:Moscow\?format\="%%l:+%%c%%20%%t%%60%%w&period=60")'
set -g status-right "$WEATHER ..."
```
![wttr.in in tmux status bar](https://wttr.in/files/example-tmux-status-line.png)

### Weechat

To embed in to an IRC ([Weechat](https://github.com/weechat/weechat)) client's existing status bar:

```
/alias add wttr /exec -pipe "/set plugins.var.python.text_item.wttr all" url:wttr.in/Montreal?format=%l:+%c+%f+%h+%p+%P+%m+%w+%S+%s
/trigger add wttr timer 60000;0;0 "" "" "/wttr"
/eval /set weechat.bar.status.items ${weechat.bar.status.items},wttr
/eval /set weechat.startup.command_after_plugins ${weechat.startup.command_after_plugins};/wttr
```
![wttr.in in weechat status bar](https://i.imgur.com/IyvbxjL.png)

### Emojis support

To see emojis in terminal, you need:

1. Terminal support for emojis (was added to Cairo 1.15.8);
2. Font with emojis support.

For the emoji font, we recommend *Noto Color Emoji*, and a good alternative option would be the *Emoji One* font;
both of them support all necessary emoji glyphs.

Font configuration:

```xml
$ cat ~/.config/fontconfig/fonts.conf
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE fontconfig SYSTEM "fonts.dtd">
<fontconfig>
  <alias>
    <family>serif</family>
    <prefer>
      <family>Noto Color Emoji</family>
    </prefer>
  </alias>
  <alias>
    <family>sans-serif</family>
    <prefer>
      <family>Noto Color Emoji</family>
    </prefer>
  </alias>
  <alias>
    <family>monospace</family>
    <prefer>
      <family>Noto Color Emoji</family>
    </prefer>
  </alias>
</fontconfig>
```

(to apply the configuration, run `fc-cache -f -v`).

In some cases, `tmux` and the terminal understanding of some emoji characters may differ, which may
cause strange effects similar to that described in #579.

## Data-rich output format (v2)

In the experimental data-rich output format, that is available under the view code `v2`,
a lot of additional weather and astronomical information is available:

* Temperature, and precipitation changes forecast throughout the days;
* Moonphase for today and the next three days;
* The current weather condition, temperature, humidity, wind speed and direction, pressure;
* Timezone;
* Dawn, sunrise, noon, sunset, dusk time for he selected location;
* Precise geographical coordinates for the selected location.

```
  $ curl v2.wttr.in/München
```

or

```
  $ curl wttr.in/München?format=v2
```

or, if you prefer Nerd Fonts instead of Emoji, `v2d` (day) or `v2n` (night):

```
  $ curl v2d.wttr.in/München
```


![data-reach output format](https://wttr.in/files/example-wttr-v2.png)

(The mode is experimental, and it has several limitations currently:

* It works only in terminal;
* Only English is supported).

Currently, you need some tweaks for some terminals, to get the best possible visualization.

### URXVT

Depending on your configuration you might be taking all steps, or only a few. URXVT currently doesn't support emoji related fonts, but we can get almost the same effect using *Font-Symbola*. So add to your `.Xresources` file the following line:
```
    xft:symbola:size=10:minspace=False
```
You can add it _after_ your preferred font and it will only show up when required.
Then, if you see or feel like you're having spacing issues, add this: `URxvt.letterSpace: 0`
For some reason URXVT sometimes stops deciding right the word spacing and we need to force it this way.

The result, should look like:

![URXVT Emoji line](https://user-images.githubusercontent.com/24360204/63842949-1d36d480-c975-11e9-81dd-998d1329bd8a.png)

## Map view (v3)

In the experimental map view, that is available under the view code `v3`,
weather information about a geographical region is available:

```
    $ curl v3.wttr.in/Bayern.sxl
```

![v3.wttr.in/Bayern](https://v3.wttr.in/Bayern.png)

or directly in browser:

*   https://v3.wttr.in/Bayern

The map view currently supports three formats:

* PNG (for browser and messengers);
* Sixel (terminal inline images support);
* IIP (terminal with iterm2 inline images protocol support).

Terminal with inline images protocols support:

⟶ *Detailed article: [Images in terminal](doc/terminal-images.md)*

| Terminal              | Environment    | Images support | Protocol |
| --------------------- | --------- | ------------- | --------- |
| uxterm                |   X11     |   yes         |   Sixel   |
| mlterm                |   X11     |   yes         |   Sixel   |
| kitty                 |   X11     |   yes         |   Kitty   |
| wezterm               |   X11     |   yes         |   IIP     |
| Darktile              |   X11     |   yes         |   Sixel   |
| Jexer                 |   X11     |   yes         |   Sixel   |
| GNOME Terminal        |   X11     |   [in-progress](https://gitlab.gnome.org/GNOME/vte/-/issues/253) |   Sixel   |
| alacritty             |   X11     |   [in-progress](https://github.com/alacritty/alacritty/issues/910) |  Sixel   |
| foot                  |  Wayland  |   yes         |   Sixel   |
| DomTerm               |   Web     |   yes         |   Sixel   |
| Yaft                  |   FB      |   yes         |   Sixel   |
| iTerm2                |   Mac OS X|   yes         |   IIP     |
| mintty                | Windows   |   yes         |   Sixel   |
| Windows Terminal  |   Windows     |   [in-progress](https://github.com/microsoft/terminal/issues/448) |   Sixel   |
| [RLogin](http://nanno.dip.jp/softlib/man/rlogin/) | Windows | yes         |   Sixel   |   |


## Different output formats

### JSON output

The JSON format is a feature providing access to *wttr.in* data through an easy-to-parse format, without requiring the user to create a complex script to reinterpret wttr.in's graphical output.

To fetch information in JSON format, use the following syntax:

    $ curl wttr.in/Detroit?format=j1

This will fetch information on the Detroit region in JSON format. The j1 format code is used to allow for the use of other layouts for the JSON output.

The result will look something like the following:
```json
{
	"current_condition": [
		{
		    "FeelsLikeC": "25",
		    "FeelsLikeF": "76",
		    "cloudcover": "100",
		    "humidity": "76",
		    "observation_time": "04:08 PM",
		    "precipMM": "0.2",
		    "pressure": "1019",
		    "temp_C": "22",
		    "temp_F": "72",
		    "uvIndex": 5,
		    "visibility": "16",
		    "weatherCode": "122",
		    "weatherDesc": [
			{
			    "value": "Overcast"
			}
		    ],
		    "weatherIconUrl": [
			{
			    "value": ""
			}
		    ],
		    "winddir16Point": "NNE",
		    "winddirDegree": "20",
		    "windspeedKmph": "7",
		    "windspeedMiles": "4"
		}
	],
...
```

Most of these values are self-explanatory, aside from `weatherCode`. The `weatherCode` is an enumeration which you can find at either [the WorldWeatherOnline website](https://www.worldweatheronline.com/developer/api/docs/weather-icons.aspx) or [in the wttr.in source code](https://github.com/chubin/wttr.in/blob/master/lib/constants.py).

### Prometheus Metrics Output

The [Prometheus](https://github.com/prometheus/prometheus) Metrics format is a feature providing access to *wttr.in* data through an easy-to-parse format for monitoring systems, without requiring the user to create a complex script to reinterpret wttr.in's graphical output. 

To fetch information in Prometheus format, use the following syntax:

    $ curl wttr.in/Detroit?format=p1

This will fetch information on the Detroit region in Prometheus Metrics format. The `p1` format code is used to allow for the use of other layouts for the Prometheus Metrics output.

A possible configuration for Prometheus could look like this:

```yaml
    - job_name: 'wttr_in_detroit'
        static_configs:
            - targets: ['wttr.in']
        metrics_path: '/Detroit'
        params:
            format: ['p1']
```

The result will look something like the following:


    # HELP temperature_feels_like_celsius Feels Like Temperature in Celsius
    temperature_feels_like_celsius{forecast="current"} 7
    # HELP temperature_feels_like_fahrenheit Feels Like Temperature in Fahrenheit
    temperature_feels_like_fahrenheit{forecast="current"} 45
    [truncated]
...


## Moon phases

wttr.in can also be used to check the phase of the Moon. This example shows how to see the current Moon phase
in the full-output mode:

    $ curl wttr.in/Moon

Get the moon phase for a particular date by adding `@YYYY-MM-DD`:

    $ curl wttr.in/Moon@2016-12-25

The moon phase information uses [pyphoon](https://github.com/chubin/pyphoon) as its backend.

To get the moon phase information in the online mode, use `%m`:

    $ curl wttr.in/London?format=%m
    🌖

Keep in mind that the Unicode representation of moon phases suffers 2 caveats:

- With some fonts, the representation `🌘` is ambiguous, for it either seem
  almost-shadowed or almost-lit, depending on whether your terminal is in
  light mode or dark mode. Relying on colored fonts like `noto-fonts` works
  around this problem.

- The representation `🌘` is also ambiguous, for it means "last quarter" in
  northern hemisphere, but "first quarter" in souther hemisphere. It also means
  nothing in tropical zones. This is a limitation that
  [Unicode](https://www.unicode.org/L2/L2017/17304-moon-var.pdf) is aware about.
  But it has not been worked around at `wttr.in` yet.

See #247, #364 for the corresponding tracking issues,
and [pyphoon#1](https://github.com/chubin/pyphoon/issues/1) for pyphoon. Any help is welcome.

## Internationalization and localization

wttr.in supports multilingual locations names that can be specified in any language in the world
(it may be surprising, but many locations in the world don't have an English name).

The query string should be specified in Unicode (hex-encoded or not). Spaces in the query string
must be replaced with `+`:

    $ curl wttr.in/станция+Восток
    Weather report: станция Восток

                   Overcast
          .--.     -65 – -47 °C
       .-(    ).   ↑ 23 km/h
      (___.__)__)  15 km
                   0.0 mm

The language used for the output (except the location name) does not depend on the input language
and it is either English (by default) or the preferred language of the browser (if the query
was issued from a browser) that is specified in the query headers (`Accept-Language`).

The language can be set explicitly when using console clients by using command-line options like this:

    curl -H "Accept-Language: fr" wttr.in
    http GET wttr.in Accept-Language:ru

The preferred language can be forced using the `lang` option:

    $ curl wttr.in/Berlin?lang=de

The third option is to choose the language using the DNS name used in the query:

    $ curl de.wttr.in/Berlin

wttr.in is currently translated into 54 languages, and the number of supported languages is constantly growing.

See [/:translation](https://wttr.in/:translation) to learn more about the translation process,
to see the list of supported languages and contributors, or to know how you can help to translate wttr.in
in your language.

![Queries to wttr.in in various languages](https://pbs.twimg.com/media/C7hShiDXQAES6z1.jpg)

## Windows Users

There are currently two Windows related issues that prevent the examples found on this page from working exactly as expected out of the box. Until Microsoft fixes the issues, there are a few workarounds. To circumvent both issues you may use a shell such as `bash` on the [Windows Subsystem for Linux (WSL)](https://docs.microsoft.com/en-us/windows/wsl/install-win10) or read on for alternative solutions.

### Garbage characters in the output
There is a limitation of the current Win32 version of `curl`. Until the [Win32 curl issue](https://github.com/chubin/wttr.in/issues/18#issuecomment-474145551) is resolved and rolled out in a future Windows release, it is recommended that you use Powershell’s `Invoke-Web-Request` command instead:
- `(Invoke-WebRequest https://wttr.in).Content`

### Missing or double wide diagonal wind direction characters
The second issue is regarding the width of the diagonal arrow glyphs that some Windows Terminal Applications such as the default `conhost.exe` use. At the time of writing this, `ConEmu.exe`, `ConEmu64.exe` and Terminal Applications built on top of ConEmu such as Cmder (`cmder.exe`) use these double-wide glyphs by default. The result is the same with all of these programs, either a missing character for certain wind directions or a broken table in the output or both. Some third-party Terminal Applications have addressed the wind direction glyph issue but that fix depends on the font and the Terminal Application you are using.
One way to display the diagonal wind direction glyphs in your Terminal Application is to use [Windows Terminal](https://www.microsoft.com/en-us/p/windows-terminal-preview/9n0dx20hk701?activetab=pivot:overviewtab) which is currently available in the [Microsoft Store](https://www.microsoft.com/en-us/p/windows-terminal-preview/9n0dx20hk701?activetab=pivot:overviewtab). Windows Terminal is currently a preview release and will be rolled out as the default Terminal Application in an upcoming release. If your output is still skewed after using Windows Terminal then try maximizing the terminal window.
Another way you can display the diagonal wind direction is to swap out the problematic characters with forward and backward slashes as shown [here](https://github.com/chubin/wttr.in/issues/18#issuecomment-405640892).

## Installation

To install the application:

1. Install external dependencies
2. Install Python dependencies used by the service
3. Configure IP2Location (optional)
4. Get a WorldWeatherOnline API and configure wego
5. Configure wttr.in
6. Configure the HTTP-frontend service

### Install external dependencies

wttr.in has the following external dependencies:

* [golang](https://golang.org/doc/install), wego dependency
* [wego](https://github.com/schachmat/wego), weather client for terminal

After you install [golang](https://golang.org/doc/install), install `wego`:

    $ go get -u github.com/schachmat/wego
    $ go install github.com/schachmat/wego

### Install Python dependencies

Python requirements:

* Flask
* geoip2
* geopy
* requests
* gevent

If you want to get weather reports as PNG files, you'll also need to install:

* PIL
* pyte (>=0.6)
* necessary fonts

You can install most of them using `pip`.

Some python package use LLVM, so install it first:

    $ apt-get install llvm-7 llvm-7-dev

If `virtualenv` is used:

    $ virtualenv -p python3 ve
    $ ve/bin/pip3 install -r requirements.txt
    $ ve/bin/python3 bin/srv.py

Also, you need to install the geoip2 database.
You can use a free database GeoLite2 that can be downloaded from (http://dev.maxmind.com/geoip/geoip2/geolite2/).

### Configure IP2Location (optional)

If you want to use the IP2location service for IP-addresses that are not covered by GeoLite2,
you have to obtain a API key of that service, and after that save into the `~/.ip2location.key` file:

```
$ echo 'YOUR_IP2LOCATION_KEY' > ~/.ip2location.key
```

If you don't have this file, the service will be silently skipped (it is not a big problem,
because the MaxMind database is pretty good).

### Installation with Docker

* Install Docker
* Build Docker Image
   ```
   docker build . --build-arg geolite_license_key=************
   ```
* These files should be mounted by the user at runtime:

```
/root/.wegorc
/root/.ip2location.key (optional)
```

### Get a WorldWeatherOnline key and configure wego

To get a WorldWeatherOnline API key, you must register here:

    https://developer.worldweatheronline.com/auth/register

After you have a WorldWeatherOnline key, you can save it into the
WWO key file: `~/.wwo.key`

Also, you have to specify the key in the `wego` configuration:

```json
$ cat ~/.wegorc
{
	"APIKey": "00XXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"City": "London",
	"Numdays": 3,
	"Imperial": false,
	"Lang": "en"
}
```

The `City` parameter in `~/.wegorc` is ignored.

### Configure wttr.in

Configure the following environment variables that define the path to the local `wttr.in`
installation, to the GeoLite database, and to the `wego` installation. For example:

```bash
export WTTR_MYDIR="/home/igor/wttr.in"
export WTTR_GEOLITE="/home/igor/wttr.in/GeoLite2-City.mmdb"
export WTTR_WEGO="/home/igor/go/bin/wego"
export WTTR_LISTEN_HOST="0.0.0.0"
export WTTR_LISTEN_PORT="8002"
```


### Configure the HTTP-frontend service

It's recommended that you also configure the web server that will be used to access the service:

```nginx
server {
	listen [::]:80;
	server_name  wttr.in *.wttr.in;
	access_log  /var/log/nginx/wttr.in-access.log  main;
	error_log  /var/log/nginx/wttr.in-error.log;

	location / {
	    proxy_pass         http://127.0.0.1:8002;

	    proxy_set_header   Host             $host;
	    proxy_set_header   X-Real-IP        $remote_addr;
	    proxy_set_header   X-Forwarded-For  $remote_addr;

	    client_max_body_size       10m;
	    client_body_buffer_size    128k;

	    proxy_connect_timeout      90;
	    proxy_send_timeout         90;
	    proxy_read_timeout         90;

	    proxy_buffer_size          4k;
	    proxy_buffers              4 32k;
	    proxy_busy_buffers_size    64k;
	    proxy_temp_file_write_size 64k;

	    expires                    off;
	}
}
```
