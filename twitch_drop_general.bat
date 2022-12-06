@echo off
set arg1=%1
set arg2=%2

start https://www.twitch.tv/drops/inventory
timeout 10
start %arg1%
timeout %arg2%
echo close