FROM ubuntu:18.04

RUN apt-get update && \
    apt-get install -y curl \
                       git \
                       python3 \
                       python3-pip \
                       python3-dev \
                       autoconf \
                       libtool \
                       gawk
RUN curl -O https://dl.google.com/go/go1.10.3.linux-amd64.tar.gz
RUN tar xvf go1.10.3.linux-amd64.tar.gz

WORKDIR /app

COPY ./bin /app/bin
COPY ./lib /app/lib
COPY ./share /app/share
COPY ./requirements.txt /app
COPY ./share/we-lang/we-lang.go /app

# There are several files that must be fetched/created manually
# before building the image
COPY ./.wegorc /root
# COPY ./.ip2location.key /root
COPY ./airports.dat /app
COPY ./GeoLite2-City.mmdb /app

RUN export PATH=$PATH:/go/bin && \ 
    go get -u github.com/mattn/go-colorable && \
    go get -u github.com/klauspost/lctime && \
    go get -u github.com/mattn/go-runewidth && \
    export GOBIN="/root/go/bin" && \
    go install /app/we-lang.go

RUN pip3 install -r requirements.txt

RUN mkdir /app/cache
RUN mkdir -p /var/log/supervisor && \
    mkdir -p /etc/supervisor/conf.d
RUN chmod -R o+rw /var/log/supervisor && \
    chmod -R o+rw /var/run
COPY share/docker/supervisord.conf /etc/supervisor/supervisord.conf

ENV WTTR_MYDIR="/app"
ENV WTTR_GEOLITE="/app/GeoLite2-City.mmdb"
ENV WTTR_WEGO="/root/go/bin/we-lang"
ENV WTTR_LISTEN_HOST="0.0.0.0"
ENV WTTR_LISTEN_PORT="8002"

EXPOSE 8002

CMD ["/usr/local/bin/supervisord", "-c", "/etc/supervisor/supervisord.conf"]
