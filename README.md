# DuckDuckGo Email Bot

![GitHub Repo stars](https://img.shields.io/github/stars/xooxoodark/ddg-email-bot?style=social)

## Features

* No More DuckDuckGo Browser Extension
* Single Binary File and A Bot Under Ur Management

## DemoBOT

[DuckDuckGoEmailBot](https://t.me/duckduckgoemail_bot)

## Deployment

1. Compile
```
git clone https://github.com/xooxoodark/ddg-email-bot.git

cd ddg-email-bot
CGO_ENABLED=0 go build -trimpath  -ldflags "-s -w"
```

2. Run
```
export TELEGRAM_APITOKEN=YOUR_TOKEN
./ddg-email-bot
```

## Referenced Repo and Special Thanks for them

[DDG Email Panel](https://github.com/whatk233/ddg-email-panel)

[ddg-get-api-token](https://github.com/timedin-de/ddg-get-api-token)