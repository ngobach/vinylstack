package main

import (
	"fmt"
	"os"
	"vinylstack/core"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

func pack(ctx *cli.Context) error {
	srcDir := ctx.Path("src")
	destDir := ctx.Path("dest")
	shouldClean := ctx.Bool("clean")
	pack := core.NewPack(destDir)

	if shouldClean {
		fmt.Println("Cleaning destination directory")
		pack.Clean()
	}

	if len(ctx.String("default-cover")) > 0 {
		pack.SetDefaultCover(ctx.String("default-cover"))
	} else {
		pack.SetDefaultCover("https://zmp3-photo-fbcrawler.zmdcdn.me/avatars/7/3/73688444a73a76169d03b689a7e785cf_1404904575.jpg")
	}

	pack.ImportMp3Directory(srcDir)
	pack.WriteIndex()
	color.Green("Directory imported successfully")

	return nil
}

func main() {
	app := &cli.App{
		Name:  "vinylstack",
		Usage: "working with Vinyl mp3 packages",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:     "src",
				Required: true,
			},
			&cli.PathFlag{
				Name:     "dest",
				Required: true,
			},
			&cli.BoolFlag{
				Name: "clean",
			},
			&cli.StringFlag{
				Name: "default-cover",
			},
		},
		Action: pack,
	}

	if err := app.Run(os.Args); err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}
}
