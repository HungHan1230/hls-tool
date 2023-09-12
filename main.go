package main

import (
	"log"
	"os"
	"sort"

	"github.com/HungHan1230/hls-tool/actions"
	"github.com/urfave/cli/v2"
)

var appName string = "hls-cmder"
var version string = "0.0.1"
var usage string = "self-testing hls tool"

func main() {
	app := &cli.App{
		Name:    appName,
		Version: version,
		Usage:   usage,
	}
	app.Authors = append(app.Authors, &cli.Author{Email: "hanksuworking@gmail.com", Name: "hanksu"})
	app.Commands = []*cli.Command{
		{
			Name: actions.DOWNLOAD,
			Usage: "The download command supports downloading the streaming video by providing the url of the .m3u8 file and the name of the video. In this action, you can choose whether download .m3u8 directly (subcommand: --directly true) from a url or copy the content of m3u8 (subcommand: --txt <path to m3u8.txt>)",
			Action: actions.DownloadVideo,
			Flags: actions.GetCmdFlags(actions.DOWNLOAD),
		},
		{
			Name: actions.CONVERT,
			Usage: "The convert command supports to convert the .ts files to .mp4 files by providing a path.",
			Action: actions.ConvertVideos,
			Flags: actions.GetCmdFlags(actions.CONVERT),
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
