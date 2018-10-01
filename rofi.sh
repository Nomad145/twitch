#!/usr/bin/bash

STREAM=`twitch -list | rofi -dmenu -matching fuzzy -only-match -location 0 -p "> "`

twitch $STREAM
