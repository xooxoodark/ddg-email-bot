# DuckDuckGo Email Bot

![GitHub Repo stars](https://img.shields.io/github/stars/xooxoodark/ddg-email-bot?style=social)

## Features

* No More DuckDuckGo Browser Extension
* Single Binary File and A Bot Under Ur Management

## DemoBOT

[DuckDuckGoEmailBot](https://t.me/duckduckgoemail_bot)

## Deployment

0. Get A UserName From DuckDuckGO Email

```
visited https://duckduckgo.com/email
```

1. Compile

```bash
git clone https://github.com/xooxoodark/ddg-email-bot.git

cd ddg-email-bot
CGO_ENABLED=0 go build -trimpath  -ldflags "-s -w"
```

2. Run

```bash
export TELEGRAM_APITOKEN=YOUR_TOKEN
./ddg-email-bot
```

## Referenced Repo and Special Thanks To Them

[DDG Email Panel](https://github.com/whatk233/ddg-email-panel)

[ddg-get-api-token](https://github.com/timedin-de/ddg-get-api-token)