# my-line-simple-bot

A LINE bot server utilized with [liff-react-boilerplate](https://github.com/BensonLiao/liff-react-boilerplate) and
[imgur-api-go-v3](https://github.com/BensonLiao/imgur-api-go-v3),
User can chat with bot, upload image to imgur or search
subreddit/account on imgur.com.

Also a barebones Go app, which can easily be deployed to Heroku.

This application supports the [Getting Started with Go on Heroku](https://devcenter.heroku.com/articles/getting-started-with-go) article - check it out.

## Running Locally

Make sure you have [Go](http://golang.org/doc/install) version 1.12 or newer and the [Heroku Toolbelt](https://toolbelt.heroku.com/) installed.

```sh
$ git clone https://github.com/heroku/my-line-simple-bot.git
$ cd my-line-simple-bot
$ go build -o bin/my-line-simple-bot -v .
github.com/mattn/go-colorable
gopkg.in/bluesuncorp/validator.v5
golang.org/x/net/context
github.com/heroku/x/hmetrics
github.com/gin-gonic/gin/render
github.com/manucorporat/sse
github.com/heroku/x/hmetrics/onload
github.com/gin-gonic/gin/binding
github.com/gin-gonic/gin
github.com/heroku/my-line-simple-bot
$ heroku local
```

Your app should now be running on [localhost:8888](http://localhost:8888/).

**Note. server read and access variable like PORT from `.env`,
usually put that in the project root.**

For example:

```
LINEBOT_CHANNEL_SECRET=1234
LINEBOT_CHANNEL_TOKEN=5678
IMGUR_CLIENT_ID=1234
IMGUR_CLIENT_SECRET=5678
PORT=8888
```

## Deploying to Heroku

```sh
$ heroku create
$ git push heroku master
$ heroku open
```

or

[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)

## Demo

Add the LINE bot as a friend:

Press the button
<a href="http://nav.cx/3tDhraO" target="_blank" rel="noopener noreferrer">
<img src="https://scdn.line-apps.com/n/line_add_friends/btn/zh-Hant.png" alt="加入好友" height="36" border="0">
</a>
or scan the QRcode:

<img src="static/my-line-bot.png" alt="reactToPost" width="360" height="360">

## Documentation

For more information about using Go on Heroku, see these Dev Center articles:

- [Go on Heroku](https://devcenter.heroku.com/categories/go)
