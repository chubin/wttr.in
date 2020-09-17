# Build stage
FROM golang:1-alpine as builder

WORKDIR /app

COPY ./share/we-lang/we-lang.go /app

RUN apk add --no-cache git


RUN true && \ 
    go get -u github.com/mattn/go-colorable && \
    go get -u github.com/klauspost/lctime && \
    go get -u github.com/mattn/go-runewidth && \
    CGO_ENABLED=0 go build /app/we-lang.go
# Results in /app/we-lang


# FROM ubuntu:18.04
# RUN apt-get update && \
#     apt-get install -y curl \
#                        git \
#                        python \
#                        python-pip \
#                        python-dev \
#                        autoconf \
#                        libtool \
#                        gawk
FROM alpine:3

WORKDIR /app

COPY ./requirements.txt /app

ENV LLVM_CONFIG=/usr/bin/llvm9-config

RUN apk add --no-cache \
    python3 \
    py3-pip \
    py3-scipy \
    py3-wheel \
    py3-gevent \
    zlib \
    jpeg \
    gcc \
    g++ \
    llvm9 \
    make \
    autoconf \
    automake \
    libtool \
    supervisor \
    zlib-dev \
    jpeg-dev \
    llvm9-dev \
    py3-numpy-dev \
    python3-dev && \
    mkdir -p /app/cache && \
    mkdir -p /var/log/supervisor && \
    mkdir -p /etc/supervisor/conf.d && \
    chmod -R o+rw /var/log/supervisor && \
    chmod -R o+rw /var/run && \
    pip install -r requirements.txt && \
    apk del --no-cache -r gcc g++ make autoconf automake libtool zlib-dev jpeg-dev llvm9-dev

COPY --from=builder /app/we-lang /app/bin/we-lang
COPY ./bin /app/bin
COPY ./lib /app/lib
COPY ./share /app/share

# These files should be mounted by the user at runtime:
# /root/.wegorc
# /root/.ip2location.key (optional)
# /app/airports.dat
# /app/GeoLite2-City.mmdb

COPY share/docker/supervisord.conf /etc/supervisor/supervisord.conf

ENV WTTR_MYDIR="/app"
ENV WTTR_GEOLITE="/app/GeoLite2-City.mmdb"
ENV WTTR_WEGO="/app/bin/we-lang"
ENV WTTR_LISTEN_HOST="0.0.0.0"
ENV WTTR_LISTEN_PORT="8002"

EXPOSE 8002

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/supervisord.conf"]
