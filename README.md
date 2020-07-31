
# *wttr.in — the right way to check the weather!*

wttr.in is a console-oriented weather forecast service that supports various information
representation methods like terminal-oriented ANSI-sequences for console HTTP clients
(curl, httpie, or wget), HTML for web browsers, or PNG for graphical viewers.

wttr.in uses [wego](http://github.com/schachmat/wego) for visualization
and various data sources for weather forecast information.

You can see it running here: [wttr.in](http://wttr.in).

## TOC

<details>
    <summary> Toggle Table of Contents</summary>

* [Usage](#usage)
  * [Weather Units](#weather-units)
* [Supported output formats and views](#supported-output-formats-and-views)
* [One-line output](#one-line-output)
* [Data-rich output format](#data-rich-output-format)
  * [URXVT](#urxvt)
* [Different output formats](#different-output-formats)
  * [JSON output](#json-output)
  * [Prometheus Metrics Output](#prometheus-metrics-output)
* [Moon phases](#moon-phases)
* [Internationalization and localization](#internationalization-and-localization)
* [Windows Users](#windows-users)
  * [Garbage characters in the output](#garbage-characters-in-the-output)
  * [Missing or double wide diagonal wind direction characters](#missing-or-double-wide-diagonal-wind-direction-characters)
* [Installation](#installation)
  * [Install OS dependencies](#install-os-dependencies)
  * [Configure Environment vars and $PATH](#configure-env-vars-and-$path)
  * [Install GO dependencies](#install-go-dependencies)
  * [Install Python dependencies](#install-python-dependencies)
  * [Download Geolite2-City DB](#download-geolite2-city-db)
  * [Get a WorldWeatherOnline API and configure wego](#get-a-worldweatherwnline-api-and-configure-wego)
  * [Configure IP2Location (optional)](#configure-ip2location-(optional))
  * [Configure the HTTP-frontend service](#configure-the-http-frontend-service)
  * [Start in dev mode](#start-in-dev-mode)
* [Docker](#docker)

</details>

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

Or in PowerShell:

```PowerShell
Invoke-RestMethod http://wttr.in
```

Want to get the weather information for a specific location? You can add the desired location to the URL in your
request like this:

    curl wttr.in/London
    curl wttr.in/Moscow
    curl wttr.in/Salt+Lake+City

If you omit the location name, you will get the report for your current location based on your IP address.

Use 3-letter airport codes in order to get the weather information at a certain airport:

    curl wttr.in/muc      # Weather for IATA: muc, Munich International Airport, Germany
    curl wttr.in/ham      # Weather for IATA: ham, Hamburg Airport, Germany

Let's say you'd like to get the weather for a geographical location other than a town or city - maybe an attraction
in a city, a mountain name, or some special location. Add the character `~` before the name to look up that special
location name before the weather is then retrieved:

    curl wttr.in/~Vostok+Station
    curl wttr.in/~Eiffel+Tower
    curl wttr.in/~Kilimanjaro

For these examples, you'll see a line below the weather forecast output that shows the geolocation
results of looking up the location:

    Location: Vostok Station, станция Восток, AAT, Antarctica [-78.4642714,106.8364678]
    Location: Tour Eiffel, 5, Avenue Anatole France, Gros-Caillou, 7e, Paris, Île-de-France, 75007, France [48.8582602,2.29449905432]
    Location: Kilimanjaro, Northern, Tanzania [-3.4762789,37.3872648]

You can also use IP-addresses (direct) or domain names (prefixed with `@`) to specify a location:

    curl wttr.in/@github.com
    curl wttr.in/@msu.ru

To get detailed information online, you can access the [/:help](http://wttr.in/:help) page:

    curl wttr.in/:help

### Weather Units

By default the USCS units are used for the queries from the USA and the metric system for the rest of the world.
You can override this behavior by adding `?u` or `?m` to a URL like this:

    curl wttr.in/Amsterdam?u
    curl wttr.in/Amsterdam?m

[table of contents](#toc)

## Supported output formats and views

wttr.in currently supports five output formats:

* ANSI for the terminal;
* Plain-text for the terminal and scripts;
* HTML for the browser;
* PNG for the graphical viewers;
* JSON for scripts and APIs;
* Prometheus metrics for scripts and APIs.

The ANSI and HTML formats are selected basing on the User-Agent string.
The PNG format can be forced by adding `.png` to the end of the query:

    wget wttr.in/Paris.png

You can use all of the options with the PNG-format like in an URL, but you have
to separate them with `_` instead of `?` and `&`:

    wget wttr.in/Paris_0tqp_lang=fr.png

Useful options for the PNG format:

* `t` for transparency (`transparency=150`);
* transparency=0..255 for a custom transparency level.

Transparency is a useful feature when weather PNGs are used to add weather data to pictures:

    convert source.jpg <( curl wttr.in/Oymyakon_tqp0.png ) -geometry +50+50 -composite target.jpg

In this example:

* `source.jpg` - source file;
* `target.jpg` - target file;
* `Oymyakon` - name of the location;
* `tqp0` - options (recommended).

![Picture with weather data](https://pbs.twimg.com/media/C69-wsIW0AAcAD5.jpg)

You can embed a special wttr.in widget, that displays the weather condition for the current or a selected location, into a HTML page using the [wttr-switcher](https://github.com/midzer/wttr-switcher). That is how it looks like: [wttr-switcher-example](https://midzer.github.io/wttr-switcher/) or [real-world-website](https://feuerwehr-eisolzried.de/.)

![Embedded wttr.in example at feuerwehr-eisolzried.de](https://user-images.githubusercontent.com/3875145/65265457-50eac180-db11-11e9-8f9b-2e1711dfc436.png)

[table of contents](#toc)

## One-line output

For one-line output format, specify additional URL parameter `format`:

    $ curl wttr.in/Nuremberg?format=3
    Nuremberg: 🌦 +11⁰C

Available preconfigured formats: 1, 2, 3, 4 and the custom format using the percent notation (see below).

You can specify multiple locations separated with `:` (for repeating queries):

    $ curl wttr.in/Nuremberg:Hamburg:Berlin?format=3
    Nuremberg: 🌦 +11⁰C

Or to process all this queries at once:

    $ curl -s 'wttr.in/{Nuremberg,Hamburg,Berlin}?format=3'
    Nuremberg: 🌦 +11⁰C
    Hamburg: 🌦 +8⁰C
    Berlin: 🌦 +8⁰C

To specify your own custom output format, use the special `%`-notation:

        c    Weather condition,
        C    Weather condition textual name,
        h    Humidity,
        t    Temperature (Actual),
        f    Temperature (Feels Like),
        w    Wind,
        l    Location,
        m    Moonphase 🌑🌒🌓🌔🌕🌖🌗🌘,
        M    Moonday,
        p    precipitation (mm),
        o    Probability of Precipitation,
        P    pressure (hPa),

        D    Dawn*,
        S    Sunrise*,
        z    Zenith*,
        s    Sunset*,
        d    Dusk*.

    (*times are shown in the local timezone)

So, these two calls are the same:

    $ curl wttr.in/London?format=3
    London: ⛅️ +7⁰C
    $ curl wttr.in/London?format="%l:+%c+%t\n"
    London: ⛅️ +7⁰C

Keep in mind, that when using in `tmux.conf`, you have to escape `%` with `%`, i.e. write there `%%` instead of `%`.

In programs, that are querying the service automatically (such as tmux), it is better to use some reasonable update interval. In tmux, you can configure it with `status-interval`.

If several, `:` separated locations, are specified in the query, specify update period
as an additional query parameter `period=`:

    set -g status-interval 60
    WEATHER='#(curl -s wttr.in/London:Stockholm:Moscow\?format\="%%l:+%%c%%20%%t%%60%%w&period=60")'
    set -g status-right "$WEATHER ..."

![wttr.in in tmux status bar](https://wttr.in/files/example-tmux-status-line.png)

To see emojis in terminal, you need:

1. Terminal support for emojis (was added to Cairo 1.15.8);
2. Font with emojis support.

For the Emoji font, we recommend *Noto Color Emoji*, and a good alternative option would be the *Emoji One* font;
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

(to apply the configuration, run `fc-cache -f -v`)
(If using MACO OSX run `brew install fontcon`)

[table of contents](#toc)

## Data-rich output format

In the experimental data-rich output format, that is available under the view code `v2`,
a lot of additional weather and astronomical information is available:

* Temperature, and precepetation changes forecast throughout the days;
* Moonphase for today and the next three days;
* The current weather condition, temperature, humidity, windspeed and direction, pressure;
* Timezone;
* Dawn, sunrise, noon, sunset, dusk time for he selected location;
* Precise geographical coordinates for the selected location.

    curl v2.wttr.in/München

or

    curl wttr.in/München?format=v2

or, if you prefer Nerd Fonts instead of Emoji, `v2d` (day) or `v2n` (night):

    curl v2d.wttr.in/München

![data-reach output format](https://wttr.in/files/example-wttr-v2.png)

(The mode is experimental, and it has several limitations currently:

* It works only in terminal;
* Only English is supported).

Currently, you need some tweaks for some terminals, to get the best possible visualization.

### URXVT

Depending on your configuration you might be taking all steps, or only a few. URXVT currently doesn't support emoji related fonts, but we can get almost the same effect using *Font-Symbola*. So add to your `.Xresources` file the following line:

    xft:symbola:size=10:minspace=False

You can add it _after_ your preferred font and it will only show up when required.
Then, if you see or feel like you're having spacing issues, add this: `URxvt.letterSpace: 0`
For some reason URXVT sometimes stops deciding right the word spacing and we need to force it this way.

The result, should look like:

![URXVT Emoji line](https://user-images.githubusercontent.com/24360204/63842949-1d36d480-c975-11e9-81dd-998d1329bd8a.png)

[table of contents](#toc)

## Different output formats

### JSON output

The JSON format is a feature providing access to *wttr.in* data through an easy-to-parse format, without requiring the user to create a complex script to reinterpret wttr.in's graphical output.

To fetch information in JSON format, use the following syntax:

    curl wttr.in/Detroit?format=j1

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

    curl wttr.in/Detroit?format=p1

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

[table of contents](#toc)

## Moon phases

wttr.in can also be used to check the phase of the Moon. This example shows how to see the current Moon phase
in the full-output mode:

    curl wttr.in/Moon

Get the Moon phase for a particular date by adding `@YYYY-MM-DD`:

    curl wttr.in/Moon@2016-12-25

The Moon phase information uses [pyphoon](https://github.com/chubin/pyphoon) as its backend.

To get the moon phase information in the online mode, use `%m`:

    $ curl wttr.in/London?format=%m
    🌖

Keep in mid that the Unicode representation of moonphases suffers 2 caveats:

* With some fonts, the representation `🌘` is ambiguous, for it either seem
  almost-shadowed or almost-lit, depending on whether your terminal is in
  light mode or dark mode. Relying on colored fonts like `noto-fonts` works
  around this problem.

* The representation `🌘` is also ambiguous, for it means "last quarter" in
  northern hemisphere, but "first quarter" in souther hemisphere. It also means
  nothing in tropical zones. This is a limitation that
  [Unicode](https://www.unicode.org/L2/L2017/17304-moon-var.pdf) is aware about.
  But it has not been worked around at `wttr.in` yet.

See #247, #364 for the corresponding tracking issues,
and [pyphoon#1](https://github.com/chubin/pyphoon/issues/1) for pyphoon. Any help is welcome.

[table of contents](#toc)

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

[table of contents](#toc)

## Windows Users

There are currently two Windows related issues that prevent the examples found on this page from working exactly as expected out of the box. Until Microsoft fixes the issues, there are a few workarounds. To circumvent both issues you may use a shell such as `bash` on the [Windows Subsystem for Linux (WSL)](https://docs.microsoft.com/en-us/windows/wsl/install-win10) or read on for alternative solutions.

### Garbage characters in the output

There is a limitation of the current Win32 version of `curl`. Until the [Win32 curl issue](https://github.com/chubin/wttr.in/issues/18#issuecomment-474145551) is resolved and rolled out in a future Windows release, it is recommended that you use Powershell’s `Invoke-Web-Request` command instead:

* `(Invoke-WebRequest http://wttr.in).Content`

### Missing or double wide diagonal wind direction characters

The second issue is regarding the width of the diagonal arrow glyphs that some Windows Terminal Applications such as the default `conhost.exe` use. At the time of writing this, `ConEmu.exe`, `ConEmu64.exe` and Terminal Applications built on top of ConEmu such as Cmder (`cmder.exe`) use these double-wide glyphs by default. The result is the same with all of these programs, either a missing character for certain wind directions or a broken table in the output or both. Some third-party Terminal Applications have addressed the wind direction glyph issue but that fix depends on the font and the Terminal Application you are using.
One way to display the diagonal wind direction glyphs in your Terminal Application is to use [Windows Terminal](https://www.microsoft.com/en-us/p/windows-terminal-preview/9n0dx20hk701?activetab=pivot:overviewtab) which is currently available in the [Microsoft Store](https://www.microsoft.com/en-us/p/windows-terminal-preview/9n0dx20hk701?activetab=pivot:overviewtab). Windows Terminal is currently a preview release and will be rolled out as the default Terminal Application in an upcoming release. If your output is still skewed after using Windows Terminal then try maximizing the terminal window.
Another way you can display the diagonal wind direction is to swap out the problematic characters with forward and backward slashes as shown [here](https://github.com/chubin/wttr.in/issues/18#issuecomment-405640892).

[table of contents](#toc)

## Installation

**Usage of wttr.in is cross platform, however, the installation has only been fully tested using Ubuntu.**

To install the application:

1. [Install OS dependencies](#install-os-dependencies)
2. [Configure Environment vars and $PATH](#configure-env-vars-and-$path)
3. [Install GO dependencies](#install-go-dependencies)
4. [Install Python dependencies](#install-python-dependencies)
5. [Download Geolite2-City DB](#download-geolite2-city-db)
6. [Get a WorldWeatherOnline API and configure wego](#get-a-worldweatherwnline-api-and-configure-wego)
7. [Configure IP2Location (optional)](#configure-ip2location-(optional))
8. [Configure the HTTP-frontend service](#configure-the-http-frontend-service)
9. [Start in dev mode](#start-in-dev-mode)

### Install OS dependencies

Ensure the OS that the wttr.in app is installed on has these dependencies:

    curl
    git
    python3
    python3-pip
    python3-dev
    autoconf
    libtool
    gettext
    gawk
    sed or gnu-sed depending on your OS

If you do not have any one/all of these then run:

    apt-get update &&           \
    apt-get install -y curl     \
        git                     \
        python3                 \
        python3-pip             \
        python3-dev             \
        autoconf                \
        libtool                 \
        gawk                    \
        gettext

or a similar command pertaining to the package manager you use (Homebrew/Chocolate).

Some packages may not be readily available to your OS, as listed above, so you will have to Google around for the package distribution that is right for your OS.

### Configure Environment vars and $PATH

    Add these environment variables to your $HOME/.bashrc, $HOME/.bash_profile or $HOME/.zshrc in the form of:
        export VARIABLE=VALUE

    Then make sure to restart all open terminals, or you can re-source the variables in all open terminals.

    ====

    $GOBIN:
        The path where installed packages are placed. Typically found /usr/local/go/bin but also dependent on where you installed GO

    $WTTR_MYDIR:
        The path to where you installed/cloned the wttr.in application. ex. $HOME/app

    $WTTR_WEGO:
        This is the path to where WEGO will be installed and if $GOBIN is set correctly can be: "$GOBIN/wego"

    $WTTR_LISTEN_HOST:
        Host for the wttr.in app server. Default: "0.0.0.0"
    
    $WTTR_LISTEN_PORT:
        Ports for the wttr.in app to listen to. Default: "8002"

    $WTTR_GEOLITE:
        The path to your GeoLite-city.mmdb file. ex. $HOME/GeoLite2-City.mmdb

    Add the $GOBIN and wttr.in/lib directories to your $PATH. example PATH=$PATH:$GOBIN:$WTTR_MYDIR/lib

### Install GO dependencies

install GO by running:

    curl "https://dl.google.com/go/goGO_VERSION.linux-amd64.tar.gz" | tar -xz -C /usr/local

you can install GO into a different location other than `/usr/local` (which ends up being `/usr/local/go`), but you will get errors, or unexpected results when searching the weather,
if you do not update the environment variables you set up in [env and path configuration](#configure-environment-vars-and-$path), accordingly.

After GO has been installed, install tools and packages:

    go get golang.org/x/tools/cmd/godoc                             && \
    go get golang.org/x/lint/golint                                 && \
    go get -u github.com/mattn/go-colorable                         && \
    go get -u github.com/klauspost/lctime                           && \
    go get -u github.com/mattn/go-runewidth                         && \
    go get github.com/schachmat/wego

After all the packages have been installed, install WEGO:

    go install github.com/schachmat/wego

### Install Python dependencies

It is best if you create a python virtual environment to run the application. You can do this by simply running:

    python -m venv PATH/TO/VIRTUAL/ENVIRONMENT # ex. $HOME/pyenvs/wttr.in

After creating enter this into newly created terminal sessions:

    source PATH/TO/VIRTUAL/ENVIRONMENT/bin/activate # ex. $HOME/pyenvs/wttr.in/bin/activate

Using the example above, you should see something along the lines of this in your terminal:

    (wttr.in)$

If not using a pyvenv, you will have to reassign your python handler:

    python3 ./some-py-file

to:
    python ./some-py-file

and can do this by (symlink):
    ln -s $(which python3) /usr/bin/python
    ln -s $(which pip3) /usr/bin/pip

This way python 2x does not get used over python 3x.

Once you have the python command updated, start with installing Pyphoon, in order to utilize moon phase requests:

    git clone https://github.com/chubin/pyphoon.git $HOME/pyphoon
    pip install $HOME/pyphoon

Then install the project dependencies from pip:

    pip install -r $WTTR_MYDIR/requirements.txt

### Download Geolite2-City DB

GeoLite2 https://www.maxmind.com/en/geolite2/signup for free database access. After registration login to your [account portal](https://www.maxmind.com/en/account/login) and on the left-hand side select `Download Files` under the `GeoIP2 / GeoLite2` section and download/unpack the GeoLite2 City database. Then add the GeoLite2-city.mmdb file to $WTTR_GEOLITE/GeoLite2-City.mmdb

### Get a WorldWeatherOnline key and configure wego

[Worldweatheronline](http://www.worldweatheronline.com/developer/signup.aspx)  
**Ths API is required to run SERVER:PORT/PATH?format=option requests**

If you want to use the Worldweatheronline API upate the $HOME/.wegorc:

    wwo-api-key=API_KEY
    backend=worldweatheronline

You will also need to add this API into $HOME/.wwo.key ex. echo "API_KEY" >> $HOME/.wwo.key

[Openweathermap]( https://home.openweathermap.org/users/sign_up)

If you want to use the Openweathermap API upate the $HOME/.wegorc:

    owm-api-key=API_KEY
    backend=openweathermap

[forcast.io](https://blog.darksky.net/) has been absorbed by Apple and only available for existing API holders, until the end of 2021.  
No longer can get new API keys but if you have one and it still works, by all means

If you want to use the forcast.io API upate the $HOME/.wegorc:

    forecast-api-key=API_KEY
    backend=forecast.io

### Configure IP2Location (optional)

[IP2Location](https://www.ip2location.com/register) **Optional**

If you have an ip2location API add it to:

    $HOME/.ip2location.key ex. echo "IP2_API" >> $HOME/.ip2location.key

If you don't have this file, the service will be silently skipped (it is not a big problem,
because the MaxMind database is pretty good).

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

### Start in dev mode

For dev environment you can run:

    python $WTTR_MYDIR/bin/srv.py
    python $WTTR_MYDIR/bin/proxy.py
    python $WTTR_MYDIR/bin/geo-proxy.py

Running each command in its own terminal session, or you can utilize [supervisord](http://supervisord.org/introduction.html):

    mkdir -p /var/log/supervisor        && \
    mkdir -p /etc/supervisor/conf.d     && \
    chmod -R o+rw /var/log/supervisor   && \
    chmod -R o+rw /var/run

Then copy the supervisord.conf to your host:

    cp ./share/docker/supervisord.conf /etc/supervisor/supervisord.conf

Run this in the terminal:

    /usr/local/bin/supervisord

[table of contents](#toc)

## Docker

Start by running:

    make init

This command creates .env template if nonexistent and shows this message:

    # The .env file has been created and located in the root of this project, where this Makefile is also located.

    <<-- ADD YOUR API KEYS TO THE .env FILE BEFORE RUNNING THIS CONTAINER -->>

    # WorldWeatherOnline API key is required. You must also set the WTTR_BACKEND ENV var in the Dockerfile, respectively.

    {
            # openworldmap
            OWM_API=API_KEY

            # worldweatheronline
            WWO_API=API_KEY

            # This option has been absorbed by Apple (https://blog.darksky.net/) and only available for existing API holders, until the end of 2021.
            FORCAST_API=API_KEY (forcast.io)
    } 

    # This is an optional API key and may have trials/free subscriptions available.
    IP2_LOCATION_API=API_KEY

    # You will also need to create the GeoLite2-City db.
    Download the GeoLite2-city db https://www.maxmind.com/en/geolite2/signup for free database access.
    After registration login to your account portal https://www.maxmind.com/en/account/login and on the left-hand side 
    select Download Files under the GeoIP2 / GeoLite2 section and download/unpack the GeoLite2 City database
    add the GeoLite2-city.mmdb file to PATH/TO/WTTR.IN/share/docker/GeoLite2-City.mmdb location.

**Follow the steps above, in the message, then continue forward.**

Start docker container normally

    make run

Start docker container in detached mode (no output to terminal)

    make run-detached

Start docker container and forces a docker rebuild (Usually when you update config files and the yaml files)

    make run-build

Start docker container in detached mode and forces a docker rebuild (Usually when you update config files and the yaml files)

    make run-build-d

For Docker specific commands check out [Docker's Documentation](https://docs.docker.com/compose/reference/)
