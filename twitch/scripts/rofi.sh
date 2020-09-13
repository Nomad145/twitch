#!/usr/bin/bash

export TWITCH_CLIENT_SECRET=
STREAM=`twitch -list | rofi -dmenu -matching fuzzy -location 0 -p "> "`

twitch $STREAM
