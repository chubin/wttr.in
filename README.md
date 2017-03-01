*wttr.in — the right way to check the weather.*

wttr.in is a console oriented weather forecast service, that supports various information representation methods
like terminal oriented ANSI-sequences for console HTTP clients such as curl, httpie or wget;
HTML for web browsers; or PNG for graphical viewers. 
wttr.in uses [wego](http://github.com/schachmat/wego) for visualization
and various data sources for weather forecast information.

You can check it at [wttr.in](http://wttr.in).

## Usage

You can access the service from a shell or from a Web browser:

    $ curl wttr.in
    Weather for City: Paris, France

         \   /     Clear
          .-.      10 – 11 °C     
       ― (   ) ―   ↑ 11 km/h      
          `-’      10 km          
         /   \     0.0 mm         


That is how the actual weather report for your location looks like (it is live!):

![Weather Report](http://wttr.in/MyLocation.png?)

(it's not your location actually, becasue GitHub's CDN hides your real IP address with its own IP address,
but it is still a live weather report in your language).

You can specify the location, for that you want to get the weather information.
If you omit the location name, you will get the information for your current location,
based on your IP address.

    $ curl wttr.in/London
    $ curl wttr.in/Moscow

You can use 3-letters airport codes if you want to get the weather information
about some airports:

    $ curl wttr.in/muc      # Weather for IATA: muc, Munich International Airport, Germany
    $ curl wttr.in/ham      # Weather for IATA: ham, Hamburg Airport, Germany

You can also use IP-addresses (direct) or domain names (prefixed with @)
as a location specificator:

    $ curl wttr.in/@github.com
    $ curl wttr.in/@msu.ru

To get detailed information online, you can access the [/:help](http://wttr.in/:help) page:

    $ curl wttr.in/:help

## Additional options

By default the USCS units are used for the queries from the USA and the 
metric system for the rest for the world.
You can override this behavior with the following options:

    $ curl wttr.in/Amsterdam?u
    $ curl wttr.in/Amsterdam?m

## Special pages

wttr.in can be used not only to check the wheather, but for some other purposes also:

    $ curl wttr.in/Moon

To see the current Moon phase (uses [pyphoon](https://github.com/chubin/pyphoon) as its backend).

    $ curl wttr.in/Moon@2016-12-25

To see the Moon phase for the specified date (2016-12-25).

## Translations

wttr.in is currently translated in more than 30 languages.

The preferred language is detected automatically basing on the query headers (`Accept-Language`)
which is automatically set by the browser or can be set explicitly using command line options
in console clients (for example: `curl -H "Accept-Language: fr"`).

The preferred language can be forced using the `lang` option:

    $ curl wttr.in/Berlin?lang=de

See [/:translation](http://wttr.in/:translation) to learn more about the translation process, 
to see the list of the supported languages, or to know how you can help to translate wttr.in in your language.

## Installation 

To install the program you need:

1. Install external dependencies
2. Install python dependencies used by the service
3. Get WorldWeatherOnline API Key
4. Configure wego
5. Configure wttr.in
6. Configure HTTP-frontend service

### Install external dependencies

External requirements:

* [wego](https://github.com/schachmat/wego), weather client for terminal

To install `wego` you must have [golang](https://golang.org/doc/install) installed. After that:

    go get -u github.com/schachmat/wego
    go install github.com/schachmat/wego

### Install python dependencies

Python requirements:

* Flask
* geoip2
* geopy
* requests
* gevent

You can install them using `pip`. 

If `virtualenv` is used:

    virtualenv ve
    ve/bin/pip install -r requirements.txt
    ve/bin/pip bin/srv.py

Also, you need to install the geoip2 database.
You can use a free database GeoLite2, that can be downloaded from http://dev.maxmind.com/geoip/geoip2/geolite2/

### Get WorldWeatherOnline key

To get the WorldWeatherOnline API key, you must register here:
 
    https://developer.worldweatheronline.com/auth/register

### Configure wego

After you have the key, configure `wego`:

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

Configure the following environment variables specifing the path to the local `wttr.in`
installation, to the GeoLite database and to the `wego` installation. For example:

    export WTTR_MYDIR="/home/igor/wttr.in"
    export WTTR_GEOLITE="/home/igor/wttr.in/GeoLite2-City.mmdb"
    export WTTR_WEGO="/home/igor/go/bin/wego"
    export WTTR_LISTEN_HOST="0.0.0.0"
    export WTTR_LISTEN_PORT="8002"


### Configure HTTP-frontend service

Configure the web server, that will be used
to access the service (if you want to use a web frontend; it's recommended):

    server {
        listen [::]:80;
        server_name  wttr.in *.wttr.in;
        access_log  /var/log/nginx/wttr.in-access.log  main;
        error_log  /var/log/nginx/wttr.in-error.log;

        location /clouds_files { root /var/www/igor/; }
        location /clouds_images { root /var/www/igor/; }

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


