package actions

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func convertCmdFlags() []cli.Flag {

	return []cli.Flag{
		&cli.StringFlag{
			Name:     "video-path",
			Aliases:  []string{"p"},
			Usage:    "specify the path of the video or the directory that includes videos.",
			Required: true,
		},
	}
}

func ConvertVideos(c *cli.Context) (err error) {
	videoPath := c.String("video-path")
	if videoPath == "" {
		return fmt.Errorf("please specify the path of the video or directory")
	}

	successToDelete := false

	video, err := os.Open(videoPath)
	defer func() {
		video.Close()
		if successToDelete {
			os.Remove(videoPath)
		}
	}()

	if err != nil {
		// handle the error and return
		return fmt.Errorf("failed to open video, err=%s", err)
	}

	// This returns an *os.FileInfo type
	fileInfo, err := video.Stat()
	if err != nil {
		// error handling
		return fmt.Errorf("failed to open video stat, err=%s", err)
	}

	// IsDir is short for fileInfo.Mode().IsDir()
	if fileInfo.IsDir() {
		// file is a directory, read file names
		fmt.Printf("%s is directory\n", videoPath)
		return
	}

	// file is not a directory, execute ffmpeg directly
	err = RunFFmpeg(videoPath)
	if err != nil {
		return
	}

	// successfully convert .ts to .mp4, delete .ts
	successToDelete = true
	return
}
