# esa-feed

Feed posts to [esa.io](https://esa.io)

## Installation

```
$ go get github.com/hitsujiwool/esa-feed
```

## Example

```
$ export ESA_TOKEN=your_access_token
$ export ESA_TEAM=your_team
$ cat post.md | esa-feed -c foo/bar -m "initial post" -w "hello"
```
