package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "rblambda"
	app.Usage = "updates lambda syntax for specified version of ruby"
	app.Commands = []cli.Command{
		cli.Command{
			Name:  "1.9",
			Usage: "update lambda syntax to ruby 1.9",
			Action: func(c *cli.Context) {
				updateToOldSyntax()
			},
		},
		cli.Command{
			Name:  "2",
			Usage: "update lambda syntax to ruby 2",
			Action: func(c *cli.Context) {
				updateToNewSyntax()
			},
		},
	}
	app.Run(os.Args)
}
