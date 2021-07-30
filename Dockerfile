# Build stage
FROM golang:1-buster as builder

WORKDIR /app

COPY ./share/we-lang/we-lang.go /app
COPY ./share/we-lang/go.mod /app

RUN apt update && apt install -y git


RUN go get -u github.com/mattn/go-colorable && \
    go get -u github.com/klauspost/lctime && \
    go get -u github.com/mattn/go-runewidth && \
    go get -u github.com/schachmat/wego && \
    CGO_ENABLED=0 go build /app/we-lang.go && \
    cp $GOPATH/bin/wego /app/wego
# Results in /app/we-lang & /app/wego


FROM python:3.9-slim-buster

WORKDIR /app

RUN mkdir -p cache data log \
      /var/log/supervisor && \
    chmod -R o+rw /var/log/supervisor /var/run cache log

COPY --from=builder /app/we-lang /app/bin/we-lang
COPY --from=builder /app/wego /app/bin/wego
COPY ./bin /app/bin
COPY ./lib /app/lib
COPY ./share /app/share
COPY share/docker/supervisord.conf /etc/supervisor/supervisord.conf

# Get GeoLite2 & airports.dat
ARG geolite_license_key
RUN ( apt update && apt install -y wget && \
  cd data/ && \
  wget "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=${geolite_license_key}&suffix=tar.gz" -O - | tar -xz --strip=1 --wildcards --no-anchored '*.mmdb' && \
  wget "https://raw.githubusercontent.com/jpatokal/openflights/master/data/airports.dat" -O airports.dat )


COPY ./requirements.txt /app
# Python build time dependencies
RUN apt install -y \
    build-essential \
    autoconf \
    libtool \
    git \
    # Runtime deps
    gawk \
    tzdata && \
    pip install -r requirements.txt --no-cache-dir && \
    apt remove -y \
    build-essential \
    autoconf \
    libtool \
    git && \
    apt autoremove -y

ENV WTTR_MYDIR="/app"
ENV WTTR_LISTEN_HOST="0.0.0.0"
ENV WTTR_LISTEN_PORT="8002"

EXPOSE 8002

CMD ["/usr/local/bin/supervisord", "-c", "/etc/supervisor/supervisord.conf"]
