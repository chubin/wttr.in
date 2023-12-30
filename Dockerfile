# Build stage
FROM golang:1-alpine as builder

WORKDIR /app

COPY ./share/we-lang/ /app

RUN apk add --no-cache git

RUN go get -u github.com/mattn/go-colorable && \
    go get -u github.com/klauspost/lctime && \
    go get -u github.com/mattn/go-runewidth && \
    cd /app && CGO_ENABLED=0 go build .

# Application stage
FROM alpine:3.19

WORKDIR /app

COPY ./requirements.txt /app

ENV LLVM_CONFIG=/usr/bin/llvm11-config

RUN apk add --no-cache --virtual .build \
    autoconf \
    automake \
    g++ \
    gcc \
    jpeg-dev \
    llvm11-dev\
    make \
    zlib-dev \
    && apk add --no-cache \
    python3 \
    py3-pip \
    py3-scipy \
    py3-wheel \
    py3-gevent \
    zlib \
    jpeg \
    llvm11 \
    libtool \
    supervisor \
    py3-numpy-dev \
    python3-dev && \
    mkdir -p /app/cache && \
    mkdir -p /var/log/supervisor && \
    mkdir -p /etc/supervisor/conf.d && \
    chmod -R o+rw /var/log/supervisor && \
    chmod -R o+rw /var/run && \
    pip install -r requirements.txt --no-cache-dir && \
    apk del --no-cache -r .build

COPY --from=builder /app/wttr.in /app/bin/wttr.in
COPY ./bin /app/bin
COPY ./lib /app/lib
COPY ./share /app/share
COPY share/docker/supervisord.conf /etc/supervisor/supervisord.conf

ENV WTTR_MYDIR="/app"
ENV WTTR_GEOLITE="/app/GeoLite2-City.mmdb"
ENV WTTR_WEGO="/app/bin/wttr.in"
ENV WTTR_LISTEN_HOST="0.0.0.0"
ENV WTTR_LISTEN_PORT="8002"

EXPOSE 8002

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/supervisord.conf"]
