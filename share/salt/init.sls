wttr:
  service.running:
    - enable: True
    - watch:
      - file: /srv/ephemeral/start.sh
      - git: wttr-repo
    - require:
      - pkg: wttr-dependencies
      - git: wttr-repo
      - cmd: wego
      - archive: geolite-db

# package names are from Ubuntu 18.04, you may need to adjust if on a different distribution
wttr-dependencies:
  pkg.installed:
    - pkgs:
      - golang
      - gawk
      - python-setuptools
      - python-dev
      - python-dnspython
      - python-geoip2
      - python-geopy
      - python-gevent
      - python-flask
      - python-pil
      - authbind

wttr-repo:
  git.latest:
    - name: https://github.com/chubin/wttr.in
    - rev: master
    - target: /srv/ephemeral/wttr.in
    - require:
      - /srv/ephemeral

wttr-start:
  file.managed:
    - name: /srv/ephemeral/start.sh
    - source: salt://wttr/start.sh
    - mode: '0770'
    - user: srv
    - group: srv

wegorc:
  file.managed:
    - name: /srv/ephemeral/.wegorc
    - user: srv
    - group: srv
    - source: salt://wttr/wegorc
    - template: jinja
    - context:
      apikey: {{ pillar['wttr']['apikey'] }}

geolite-db:
   archive.extracted:
    - name: /srv/ephemeral
    - source: http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.tar.gz
    - source_hash: http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.tar.gz.md5
    - keep_source: True
    - options: --strip-components=1 # flatten directory structure
    - enforce_toplevel: False

# Could benefit from improvement, won't get updated automatically at all
wego:
  cmd.run:
    - onlyif: 'test ! -e /srv/ephemeral/bin/wego'
    - env:
      - GOPATH: /srv/ephemeral
    - name: go get -u github.com/schachmat/wego && go install github.com/schachmat/wego
    - cwd: /srv/ephemeral/
    - require:
      - pkg: wttr-dependencies
      - file: wegorc

/srv/ephemeral:
  file.directory:
    - makedirs: True

{% for dir in '/srv/ephemeral/wttr.in/log','/srv/ephemeral/wttr.in/cache' %}
{{ dir }}:
  file.directory:
    - user: srv
    - group: srv
    - makedirs: True
    - recurse:
        - user
        - group
    - require_in:
      - service: wttr
{% endfor %}

/etc/systemd/system/wttr.service:
  file:
    - managed
    - source: salt://wttr/wttr.service
    - require:
      - file: wttr-start
      - file: authbind-80
    - require_in:
      - service: wttr

authbind-80:
  file:
    - managed
    - name: /etc/authbind/byport/80
    - user: srv
    - group: srv
    - mode: 770
    - replace: False
    - require:
      - pkg: wttr-dependencies
