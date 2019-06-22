*wttr.in — the right way to check the weather!*

wttr.in is a console-oriented weather forecast service that supports various information
representation methods like terminal-oriented ANSI-sequences for console HTTP clients
(curl, httpie, or wget), HTML for web browsers, or PNG for graphical viewers.

wttr.in uses [wego](http://github.com/schachmat/wego) for visualization
and various data sources for weather forecast information.

You can see it running here: [wttr.in](http://wttr.in).

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

![Weather Report](http://wttr.in/MyLocation.png?)

(It's not your actual location - GitHub's CDN hides your real IP address with its own IP address,
but it's still a live weather report in your language.)

Want to get the weather information for a specific location? You can add the desired location to the URL in your
request like this:

    $ curl wttr.in/London
    $ curl wttr.in/Moscow

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

To get detailed information online, you can access the [/:help](http://wttr.in/:help) page:

    $ curl wttr.in/:help

### Weather Units

By default the USCS units are used for the queries from the USA and the metric system for the rest of the world.
You can override this behavior by adding `?u` or `?m` to a URL like this:

    $ curl wttr.in/Amsterdam?u
    $ curl wttr.in/Amsterdam?m

## Supported output formats

wttr.in currently supports three output formats:

* ANSI for the terminal;
* ANSI for the terminal, one-line mode;
* HTML for the browser;
* PNG for the graphical viewers.

The ANSI and HTML formats are selected basing on the User-Agent string.
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

## One-line output

For one-line output format, specify additional URL parameter `format`:

```
$ curl wttr.in/Nuremberg?format=3
Nuremberg: 🌦 +11⁰C
```

Available preconfigured formats: 1, 2, 3, 4 and the custom format using the percent notation (see below).

You can sepcify multiple locations separated with `:` (for repeating queries):

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
    h    Humidity,
    t    Temperature,
    w    Wind,
    l    Location,
    m    Moonphase 🌑🌒🌓🌔🌕🌖🌗🌘,
    M    Moonday,
    p    precipitation (mm),
    P    pressure (hPa),
```

So, these two calls are the same:

```
    $ curl wttr.in/London?format=3
    London: ⛅️ +7⁰C
    $ curl wttr.in/London?format="%l:+%c+%t"
    London: ⛅️ +7⁰C
```
Keep in mind, that when using in `tmux.conf`, you have to escape `%` with `%`, i.e. write there `%%` instead of `%`.

In progams, that are querying the service automatically (such as tmux),it is better to use some reasonable update interval. In tmux, you can configure it with `status-interval`.

If several, `:` separated locations, are specified in the query, specify update period
as an additional query parameter `period=`:
```
set -g status-interval 60
WEATHER='#(curl -s wttr.in/London:Stockholm:Moscow\?format\="%%l:+%%c%%20%%t%%60%%w&period=60")'
set -g status-right "$WEATHER ..."
```
![wttr.in in tmux status bar](https://wttr.in/files/example-tmux-status-line.png)

To see emojis in terminal, you need:

1. Terminal support for emojis (was added to Cairo 1.15.8);
2. Font with emojis support.

For the Emoji font, we recommend *Noto Color Emoji*, and a good alternative option would be the *Emoji One* font;
both of them support all necessary emoji glyphs.

Font configuration:

```
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

(to apply the configuration, run `fc-cache -f -v`)

## Moon phases

wttr.in can also be used to check the phase of the Moon. This example shows how to see the current Moon phase:

    $ curl wttr.in/Moon

Get the Moon phase for a particular date by adding `@YYYY-MM-DD`:

    $ curl wttr.in/Moon@2016-12-25

The Moon phase information uses [pyphoon](https://github.com/chubin/pyphoon) as its backend.

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

See [/:translation](http://wttr.in/:translation) to learn more about the translation process, 
to see the list of supported languages and contributors, or to know how you can help to translate wttr.in
in your language.

![Queries to wttr.in in various languages](https://pbs.twimg.com/media/C7hShiDXQAES6z1.jpg)

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

If `virtualenv` is used:

    $ virtualenv ve
    $ ve/bin/pip install -r requirements.txt
    $ ve/bin/python bin/srv.py

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

### Get a WorldWeatherOnline key and configure wego

To get a WorldWeatherOnline API key, you must register here:
 
    https://developer.worldweatheronline.com/auth/register

After you have a WorldWeatherOnline key, you can save it into the
WWO key file: `~/.wwo.key`

Also, you have to specify the key in the `wego` configuration:

    $ cat ~/.wegorc 
    {
        "APIKey": "00XXXXXXXXXXXXXXXXXXXXXXXXXXX",
        "City": "London",
        "Numdays": 3,
        "Imperial": false,
        "Lang": "en"
    }

The `City` parameter in `~/.wegorc` is ignored.

### Configure wttr.in

Configure the following environment variables that define the path to the local `wttr.in`
installation, to the GeoLite database, and to the `wego` installation. For example:

    export WTTR_MYDIR="/home/igor/wttr.in"
    export WTTR_GEOLITE="/home/igor/wttr.in/GeoLite2-City.mmdb"
    export WTTR_WEGO="/home/igor/go/bin/wego"
    export WTTR_LISTEN_HOST="0.0.0.0"
    export WTTR_LISTEN_PORT="8002"


### Configure the HTTP-frontend service

It's recommended that you also configure the web server that will be used to access the service:

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
