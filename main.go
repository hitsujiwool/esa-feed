package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Post struct {
	Name     string   `json:"name,omitempty"`
	Body     string   `json:"body_md,omitempty"`
	Tags     []string `json:"tags,omitempty"`
	Category string   `json:"category,omitempty"`
	WIP      bool     `json:"wip"`
	Message  string   `json:"message,omitempty"`
	User     string   `json:"user,omitempty"`
}

func feed(p Post, team string, token string) error {
	wrap := struct {
		Post `json:"post"`
	}{p}
	bs, err := json.Marshal(wrap)
	if err != nil {
		return err
	}
	_, err = http.Post(
		fmt.Sprintf("https://api.esa.io/v1/teams/%s/posts?access_token=%s", team, token),
		"application/json",
		bytes.NewReader(bs),
	)
	return err
}

func main() {
	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.Name}} [options] name

VERSION:
   {{.Version}}{{if or .Author .Email}}

AUTHOR:{{if .Author}}
  {{.Author}}{{if .Email}} - <{{.Email}}>{{end}}{{else}}
  {{.Email}}{{end}}{{end}}

OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}
`
	app := cli.NewApp()
	app.Name = "esa-feed"
	app.Usage = "feed posts to esa.io"
	app.Version = "0.1.0"
	app.Action = func(c *cli.Context) {
		if c.IsSet("help") {
			cli.ShowAppHelp(c)
			os.Exit(0)
		}
		if len(c.Args()) == 0 {
			cli.ShowAppHelp(c)
			os.Exit(1)
		}
		bs, _ := ioutil.ReadAll(os.Stdin)
		name := c.Args().First()
		p := Post{
			Name:     name,
			Body:     string(bs),
			WIP:      c.Bool("wip"),
			Category: c.String("category"),
			Message:  c.String("message"),
			User:     c.String("user"),
		}		
		team := c.String("team")
		if len(team) == 0 {
			fmt.Fprintf(os.Stderr, "error: %s\n", "team cannot be empty")
			os.Exit(1)
		}
		token := c.String("token")
		if len(token) == 0 {
			fmt.Fprintf(os.Stderr, "error: %s\n", "access token cannot be empty")
			os.Exit(1)
		}
		if tags := c.String("tags"); len(tags) > 0 {
			p.Tags = strings.Split(tags, ",")
		}
		if err := feed(p, team, token); err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		}
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "team",
			Usage:  "name of your team",
			EnvVar: "ESA_TEAM",
		},
		cli.StringFlag{
			Name:   "token",
			Usage:  "your API token",
			EnvVar: "ESA_ACCESS_TOKEN",
		},
		cli.StringFlag{
			Name:  "category, c",
			Usage: "category",
		},
		cli.StringFlag{
			Name:  "tags, t",
			Usage: "comma separated tags",
		},
		cli.BoolFlag{
			Name:  "wip, w",
			Usage: "WIP (default: false)",
		},
		cli.StringFlag{
			Name:  "message, m",
			Usage: "message",
		},
		cli.StringFlag{
			Name:  "user, u",
			Usage: "user",
		},
	}
	app.Run(os.Args)
}
