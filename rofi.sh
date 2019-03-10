#!/usr/bin/bash

STREAM=`twitch -list | rofi -dmenu -matching fuzzy -location 0 -p "> "`

twitch $STREAM
