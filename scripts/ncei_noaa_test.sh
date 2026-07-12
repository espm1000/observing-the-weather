#!/bin/bash
# shellcheck disable=SC2034

set -e

# Documentation: https://www.ncei.noaa.gov/cdo-web/webservices/v2#gettingStarted

TOKEN="${NCEI_TOKEN:?Error: Required NCEI token.}"
STATION='USW00014922' # KMSP
BASE_URL="https://www.ncei.noaa.gov/cdo-web/api/v2"
DATA_ENDPOINT="data"
DATASET_ENDPOINT="datasets"
DATATYPE_ENDPOINT="datatypes"
LOCATION_ENDPOINT="locations"
STATION_ENDPOINT="stations"
DATASET="GHCND" # daily-summaries
START_DATE="2025-07-06" # max range is 1 year
END_DATE="2026-07-07"
UNITS="standard" # options: "metric" and "standard"
LIMIT=1000 # max limit
OUTPUT_FILE="noaa_${START_DATE}_${END_DATE}.json"


curl -G "$BASE_URL/$DATA_ENDPOINT" \
  -H "token:$TOKEN" \
  --data-urlencode "datasetid=$DATASET" \
  --data-urlencode "stationid=$DATASET:$STATION" \
  --data-urlencode "startdate=$START_DATE" \
  --data-urlencode "enddate=$END_DATE" \
  --data-urlencode "units=$UNITS" \
  --data-urlencode "limit=$LIMIT" \
  -o $OUTPUT_FILE
