# Task description
Implement recursive web-crawler of the site.
1. Crawler is a command-line tool that accept starting URL and destination directory
2. Crawler downloads starting URL and searches for links inside of it
3. Crawler does not walk outside of initial URL, for example if starting URL is "https://start.url/abc", then it goes to "https://start.url/abc/123" and "https://start.url/abc/456", but skips "https://another.domain/" or "https://start.url/def".
4. Crawler should correctly process Ctrl+C hotkey
5. Crawler should be parallel
6. Crawler should support continue to load if the destination directory already has loaded data (if we cancel the download and then continue).

# Note
In general, it's very similar to `wget --mirror` with few extra options.
Does not necessary to implement all listed above, just do some (more important from your point of view) part, but take into account to do it in the way where every listed above is possible to add by extending your program without re-writing from scratch.

# Timeframe for the task
4 hours.

---
# How to build
1. Install Go Development Enviromnent as described [here](https://go.dev/doc/install)
2. Clone project's source code
3. Run:
```
cd <project directory>
go build 
```

# How to run
There are two ways to run this app do exist:

1. If you have GO development environment installed you can clone this project and run:
```
cd <project directory>
go run ./app.go <Target URL> <Path to output directory>
```
