# Kotaku UK RSS (for fun)

So [Kotaku UK](http://www.kotaku.co.uk) does not have an RSS feed right
now. And as I'm still teaching myself [Go](http://golang.org/), I figured it
would be a fun small project to hack together a small scraper that serves up
an RSS feed for my personal use.

## Building/Running

```
go get github.com/kr/godep
godep go build
PORT=5000 ./kotaku-uk-rss
```

## Deploying to Heroku

```
git clone git@github.com:jimeh/kotaku-uk-rss.git
cd kotaku-uk-rss
heroku create -b https://github.com/kr/heroku-buildpack-go.git
git push heroku master
```

## Technical

When launched, the process does two things:

1. Starts a goroutine that will fetch a list of articles from the first 3
   pages of kotaku.co.uk, and building and caching a RSS XML string from those
   articles. This is repeated every 60 seconds.
2. Starts a web-server on `PORT` number serving up cached RSS feed on `/rss`.

## License

```
        DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE
                    Version 2, December 2004

 Copyright (C) 2014 Jim Myhrberg

 Everyone is permitted to copy and distribute verbatim or modified
 copies of this license document, and changing it is allowed as long
 as the name is changed.

            DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE
   TERMS AND CONDITIONS FOR COPYING, DISTRIBUTION AND MODIFICATION

  0. You just DO WHAT THE FUCK YOU WANT TO.
```
