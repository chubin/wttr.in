FROM ubuntu
ENV DEBIAN_FRONTEND=noninteractive

########################################## INSTALL OS DEPENDENCIES ##########################################

RUN apt-get update &&       \
    apt-get install -y curl \
    git                     \
    python3                 \
    python3-pip             \
    python3-dev             \
    autoconf                \
    libtool                 \
    gawk                    \
    gettext

########################################## ENVIRONMENTAL VARIABLES ##########################################

# Root
ENV HOME="/root"

# GO VARS
ENV GOBIN="/usr/local/go/bin"

# GO VERSION TO BE INSTALLED
ENV GO_VERSION="go1.14.6"

# WTTR/WEGO VARS
ENV WTTR_MYDIR="$HOME/app"
ENV WTTR_WEGO="$GOBIN/wego"
ENV WTTR_LISTEN_HOST="0.0.0.0"
ENV WTTR_LISTEN_PORT="8002"

# GeoLite2 https://www.maxmind.com/en/geolite2/signup for free database access
# After registration login to your account portal https://www.maxmind.com/en/account/login and on the left-hand side 
# select `Download Files` under the `GeoIP2 / GeoLite2` section and download/unpack the GeoLite2 City database
# add the GeoLite2-city.mmdb file to PATH/TO/WTTR.IN/share/docker/GeoLite2-City.mmdb
ENV WTTR_GEOLITE="$HOME/GeoLite2-City.mmdb"

# PYPHOON VAR
ENV PYPHOON_ROOT="$HOME/pyphoon"

# WEGORC CONFIG DATA
ENV WTTR_BACKEND="openweathermap"
ENV WTTR_DEFAULT_LOCATION="Denver"

# PATH VAR UPDATE
ENV PATH=$PATH:$GOBIN:$WTTR_MYDIR/lib

#=========================================================== APIS ===========================================================#

# Openweathermap https://home.openweathermap.org/users/sign_up
ARG OWM_API

# Worldweatheronline account http://www.worldweatheronline.com/developer/signup.aspx
# Required to run http(s)://SERVER/LOCATION?format=formatter requests
ARG WWO_API

# No longer can get new API keys but if you have one and it still works, by all means
ARG FORCAST_API

# IP2Location https://www.ip2location.com/register
ARG IP2_LOCATION_API

#=========================================================== APIS ===========================================================#

# Python 3 reassignments (symlink)
RUN ln -s $(which python3) /usr/bin/python
RUN ln -s $(which pip3) /usr/bin/pip

########################################## PROJECT SET UP ##########################################

# Copy app and GEOLite DB files into container
RUN mkdir -p $WTTR_MYDIR
COPY . $WTTR_MYDIR
# Move Geolite DB to root
RUN mv $WTTR_MYDIR/share/docker/GeoLite2-City.mmdb $HOME/GeoLite2-City.mmdb

# Format IP2 Location API and .wegorc API's and copy to container, along with wwo_key
RUN envsubst < $WTTR_MYDIR/share/docker/.ip2location.key > $HOME/.ip2location.key   && \
    echo $WWO_API >> $HOME/.wwo.key                                                 && \
    envsubst < $WTTR_MYDIR/share/docker/.wegorc > $HOME/.wegorc

# Download GO
RUN curl "https://dl.google.com/go/${GO_VERSION}.linux-amd64.tar.gz" | tar -xz -C /usr/local

# Install GO tools and packages
RUN go get golang.org/x/tools/cmd/godoc                             && \
    go get golang.org/x/lint/golint                                 && \
    go get -u github.com/mattn/go-colorable                         && \
    go get -u github.com/klauspost/lctime                           && \
    go get -u github.com/mattn/go-runewidth                         && \
    go get github.com/schachmat/wego

# Install WEGO
RUN go install github.com/schachmat/wego

# Install PYPHOON
RUN git clone https://github.com/chubin/pyphoon.git $PYPHOON_ROOT   && \
    pip install $PYPHOON_ROOT

# Install project python dependencies
RUN pip install -r $WTTR_MYDIR/requirements.txt

# Create and configure supervisor
RUN mkdir -p /var/log/supervisor                                    && \
    mkdir -p /etc/supervisor/conf.d                                 && \
    chmod -R o+rw /var/log/supervisor                               && \
    chmod -R o+rw /var/run
COPY ./share/docker/supervisord.conf /etc/supervisor/supervisord.conf

# Expose port to host 
EXPOSE $WTTR_LISTEN_PORT

########################################## START ##########################################

# Startup command
CMD ["/usr/local/bin/supervisord"]
