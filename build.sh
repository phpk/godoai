#!/bin/bash

wails build -m -platform windows/amd64 -webview2 embed -nsis
wails build -m -platform linux/amd64 -webview2 embed