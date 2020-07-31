#!/usr/bin/python3
import os

SHARE_ROOT = os.path.dirname(
    os.path.dirname(
        os.path.abspath(__file__)
    )
)

ENV_TEMPLATE_FILE = os.path.join(
    SHARE_ROOT,
    'docker/.env.template'
)

ENV_FILE = os.path.join(
     SHARE_ROOT,
    '../.env'
)

# Check if .env exists and if not create using template from WTTR_MYDIR/docker/.env.template
if 	not os.path.exists( ENV_FILE ):
	proc = os.popen(f"cp {ENV_TEMPLATE_FILE} {ENV_FILE}")
else:
	print("Project has already been initialized.")
	exit()

# Text Highlights
GRN = '\033[1;32m'
BLU = '\033[1;34m'
RED = '\033[1;31m'
YEL = '\033[1;33m'
NC = '\033[0m'

# Formatted Message
INIT_MESSAGE = f"""
{BLU}# The .env file has been created and located in the root of this project, where this Makefile is also located.

{RED}<<-- ADD YOUR API KEYS TO THE {GRN}.env {RED}FILE BEFORE RUNNING THIS CONTAINER -->>

{BLU}# {RED}WorldWeatherOnline API{BLU} key is required. You must also set the WTTR_BACKEND ENV var in the Dockerfile, respectively.

{YEL}{{{GRN}
	{BLU}# openworldmap{GRN}
	OWM_API=API_KEY

	{BLU}# worldweatheronline{GRN}
	WWO_API=API_KEY

	{BLU}# This option has been absorbed by Apple (https://blog.darksky.net/) and only available for existing API holders, until the end of 2021.
	{GRN}FORCAST_API=API_KEY (forcast.io)
{YEL}}} 

{BLU}# This is an optional API key and may have trials/free subscriptions available.
{GRN}IP2_LOCATION_API=API_KEY

{BLU}# You will also need to create the GeoLite2-City db.{NC}
Download the GeoLite2-city db {RED}https://www.maxmind.com/en/geolite2/signup{NC} for free database access.
After registration login to your account portal {RED}https://www.maxmind.com/en/account/login{NC} and on the left-hand side 
select {RED}Download Files{NC} under the {RED}GeoIP2 / GeoLite2{NC} section and download/unpack the GeoLite2 City database
add the {RED}GeoLite2-city.mmdb{NC} file to {RED}PATH/TO/WTTR.IN/share/docker/GeoLite2-City.mmdb{NC} location.
"""
# Print to console
print(INIT_MESSAGE)
