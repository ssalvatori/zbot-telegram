#!/bin/env bash
#
# Requirements: curl & jq and a valid openweathermap key
#
get_weather () {
    API="http://api.openweathermap.org/data/2.5/weather"
    API_KEY=""
    CITY=${ORIGIN}
    TOKEN="appid=${API_KEY}"
    URL="${API}?q=${CITY}&${TOKEN}&units=metric&lang=es_es"
    DATA=$(curl -s ${URL})
    CURRENT=$(echo ${DATA} | jq '.main.temp')
    MIN=$(echo ${DATA} | jq '.main.temp_min')
    MAX=$(echo ${DATA} | jq '.main.temp_max')
    HUM=$(echo ${DATA} | jq '.main.humidity')
    DESC=$(echo ${DATA} | jq '.weather[0].description')
    echo "City: ${CITY} Current: ${CURRENT} ºC min: ${MIN} ºC max: ${MAX} ºC hum: ${HUM}% Des: ${DESC}"
}

case $4 in
  tll|TLL)
    ORIGIN="Tallinn,EE"
  ;;
  bru|BRU)
    ORIGIN="Brussels,Belgium"
  ;;
  scl|SCL)
    ORIGIN="Santiago,CL"
  ;;
  mad|MAD)
    ORIGIN="Madrid,ES"
  ;;
  ccp|CCP)
    ORIGIN="Concepcion,CL"
  ;;
  str|STR)
    ORIGIN="Stuttgart,DE"
  ;;
  soc|SOC)
    ORIGIN="Socorro,New%20Mexico"
  ;;
  *)
    echo "Options: mad, bru, tll, scl, ccp, str, soc"
    exit 0
  ;;
esac
get_weather $ORIGIN
exit 0